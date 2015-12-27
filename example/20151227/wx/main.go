package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/woylyn/esap2/db/sqlsrv" //MSSQL,编译前先get
	"github.com/woylyn/esap2/wechat"    //微信SDK
)

const (
	FTPPATH   = `M:\wx\hr` //网盘路径
	PICPREFIX = "/P00"     //照片前缀，默认是"P00"
	PICSUFFIX = ".jpg"     //照片后缀，默认".jpg"
)

//定义员工
type employee struct {
	Name  string
	Eid   string
	Photo string
	Rcid  interface{}
}

//员工照片上传，可为网盘
func (e *employee) download(url string) {
	imgResp, _ := http.Get(url) //从微信服务器获取照片
	defer imgResp.Body.Close()
	m, _, _ := image.Decode(imgResp.Body)         //解析照片信息
	fn := FTPPATH + PICPREFIX + e.Eid + PICSUFFIX //设置照片存放路径及文件名
	fh, _ := os.Create(fn)                        //重建照片
	defer fh.Close()
	jpeg.Encode(fh, m, nil) //以JPG格式保存照片
}

//定义员工表、存放照片链接信息
var empMap = make(map[string]employee)

type myApp struct {
	wechat.WxApp //继承内置服务器
}

//重写文本类消息处理接口
func (w *myApp) GoText() {
	//1.检验是否已上传图片，上传过则提示输入“姓名，工号”
	if v, ok := empMap[w.Req.FromUserName]; ok {
		//解析姓名工号
		empInfo := strings.Split(w.Req.Content, "，")
		v.Name, v.Eid = empInfo[0], empInfo[1]
		//通过员工信息表查找rcid,未找到则提示用户。
		v.Rcid = sqlsrv.Fetch("select Excelserverrcid from employee where name=? and eid=?", v.Name, v.Eid)
		if (v.Rcid) == nil {
			w.RespStr = "未找到该用户，请重新输入或重新上传图片。"
		} else {
			//设置照片名为：前缀 + 工号
			picName := PICPREFIX + v.Eid
			//fmt.Println("picName:", picName)
			//删除原有数据库照片路径记录
			sqlsrv.Exec("delete from es_casepic where rcid=?", v.Rcid)
			//向数据库插入新照片路径记录
			err := sqlsrv.Exec("insert es_casepic(rcid,picNo,fileType,rtfid,sh,r,c,saveinto,nfsfolderid,nfsfolder,relafolder,phyfileName) values(?,?,?,?,?,?,?,?,?,?,?,?)",
				v.Rcid, picName, ".jpg", 60, 1, 3, 4, 1, 1, `ed\esdisk`, `wx\hr`, picName+".jpg")
			if err != nil {
				w.RespStr = "图片上传失败"
			} else {
				//上传图片到网盘，销毁员工数组中的信息，回复处理成功信息
				v.download(v.Photo)
				delete(empMap, w.Req.FromUserName)
				w.RespStr = "员工照片已成功处理"
			}
		}
		//生成微信text消息回复
		w.RespB, _ = wechat.RespText(w.Req.ToUserName, w.Req.FromUserName, w.RespStr)
	} else {
		//1.未上传图片其他消息的处理，仅示例
		switch w.Req.Content {
		case "蛋蛋":
			w.RespB, _ = wechat.RespArt(w.Req.ToUserName, w.Req.FromUserName)
		default:
			sqlsrv.Exec("insert cxls(cDate,uid,keyword) values(?,?,?)", w.Req.CreateTime, w.Req.FromUserName, w.Req.Content)
			resp := sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", w.Req.Content))
			if *resp == nil {
				resp = sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", "默认"))
				if *resp == nil {
					*resp = "hi，guys,welcome to ESAPbuluo/::D"
				}
			}
			w.RespB, _ = wechat.RespText(w.Req.ToUserName, w.Req.FromUserName, fmt.Sprintf("%v", *resp))
		}
	}
}
func (w *myApp) GoImage() {
	fmt.Println("DoImg2")
	empMap[w.Req.FromUserName] = employee{"", "", w.Req.PicUrl, nil}
	w.RespB, _ = wechat.RespText(w.Req.ToUserName, w.Req.FromUserName, fmt.Sprintf("%v", "请输入员工信息(格式 姓名,工号)："))
	//	return w.RespB
}

func main() {
	//	wechat.SetToken("esap") //设置token
	//	wechat.SetDev(false)    //设置开发模式
	app := &myApp{} //实例化微信API副本
	app.Run(app)    //运行SERVER
}
