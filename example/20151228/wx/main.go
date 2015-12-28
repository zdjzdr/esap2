/**
 *微信API示例 by woylin 2015/12/28
 */
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/woylyn/esap2/db/sqlsrv" //MSSQL,编译前先get
	"github.com/woylyn/esap2/wechat"    //微信SDK包
)

const (
	FTPPATH   = `M:\wx\hr` //网盘路径，ES需启用网盘功能，这映射共享的网盘目录为M盘，再建立了wx\hr目录
	PICPREFIX = "/P00"     //照片前缀，默认是"P00"
	PICSUFFIX = ".jpg"     //照片后缀，默认".jpg"
)

//定义员工
type employee struct {
	Name     string
	Eid      string
	Photo    string
	Rcid     interface{}
	errCount int
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
var empMap = make(map[string]*employee)

//继承内置服务器
type myApp struct {
	wechat.WxApp
}

//重写图片类消息处理
func (w *myApp) GoImage() {
	empMap[w.Req.FromUserName] = &employee{"", "", w.Req.PicUrl, nil, 0}
	w.RespB, _ = wechat.RespText(w.Req.ToUserName, w.Req.FromUserName, "请输入员工信息(格式 姓名,工号)：")
}

//重写文本类消息处理
func (w *myApp) GoText() {
	//1.检验是否已上传图片，上传过则提示输入“姓名，工号”
	if v, ok := empMap[w.Req.FromUserName]; ok {
		//解析姓名工号
		empInfo := strings.Split(w.Req.Content, "，")
		if len(empInfo) == 2 {
			v.Name, v.Eid = empInfo[0], empInfo[1]
		} else {
			//重复出错4次后重置，取消流程
			v.errCount++
			if v.errCount > 3 {
				w.RespStr = "出错已超限，上传流程已取消"
				delete(empMap, w.Req.FromUserName)
			} else {
				w.RespStr = fmt.Sprintf("你填入的信息格式不正确，请重新输入。(%d)", v.errCount)
			}
			w.RespB, _ = wechat.RespText(w.Req.ToUserName, w.Req.FromUserName, w.RespStr)
			return
		}
		//通过员工信息表查找rcid,未找到则提示用户。
		v.Rcid = sqlsrv.Fetch("select Excelserverrcid from employee where name=? and eid=?", v.Name, v.Eid)
		if (v.Rcid) == nil {
			w.RespStr = "未找到该用户，请重新输入或重新上传图片。"
		} else {
			//设置照片名为：前缀 + 工号
			picName := PICPREFIX + v.Eid
			//删除原有数据库照片路径记录
			sqlsrv.Exec("delete from es_casepic where rcid=?", v.Rcid)
			//向数据库插入新照片路径记录,这里的ed\esdisk,wx\hr等信息要改成自己的ES网盘配置目录
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
			//创建四个文章
			art := wechat.CreArt("ESAP第十四弹 手把手教你玩转ES微信开发",
				"来自村长的ESAP系统最新技术分享。",
				"http://iesap.net/wp-content/uploads/2015/12/esap3-1.jpg",
				"http://iesap.net/index.php/2015/12/28/esap14/")
			art2 := wechat.CreArt("打通信息化的“任督二脉”(二)",
				"来自村长的ESAP2.0系统技术分享。",
				"http://iesap.net/wp-content/uploads/2015/12/taiji.jpg",
				"http://iesap.net/index.php/2015/12/16/esap2-1/")
			art3 := wechat.CreArt("打通信息化的“任督二脉”(一)",
				"来自村长的ESAP2.0系统技术分享。",
				"http://iesap.net/wp-content/uploads/2015/12/rdem.jpg",
				"http://iesap.net/index.php/2015/12/11/esap2-0/")
			art4 := wechat.CreArt("我与ES不吐不快的槽",
				"来自村长的工作日志。",
				"http://iesap.net/wp-content/uploads/2015/09/s13.jpg",
				"http://iesap.net/index.php/2015/09/04/es13/")
			//创建文章回复
			w.RespB, _ = wechat.RespArt(w.Req.ToUserName, w.Req.FromUserName, art, art2, art3, art4)
		default:
			//其他查询转入到关键字查询，关键字查询表在ES的数据库中，首先将用户关键字存入查询历史
			sqlsrv.Exec("insert cxls(cDate,uid,keyword) values(?,?,?)", w.Req.CreateTime, w.Req.FromUserName, w.Req.Content)
			//尝试匹配用户关键字
			resp := sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", w.Req.Content))
			//未匹配到则使用“默认”所匹配的信息
			if *resp == nil {
				resp = sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", "默认"))
				//连“默认”关键字也找不到则打印“欢迎来到ESAP部落^_^”
				if *resp == nil {
					*resp = "hi，guys,welcome to ESAPbuluo/::D"
				}
			}
			//创建文本回复
			w.RespB, _ = wechat.RespText(w.Req.ToUserName, w.Req.FromUserName, fmt.Sprintf("%v", *resp))
		}
	}
}

func main() {
	//	wechat.SetToken("esap") //设置token
	//	wechat.SetDev(false)    //关闭开发模式
	//	wechat.SetPort(":8080") //更改监听端口
	app := &myApp{} //实例化微信API副本
	app.Run(app)    //运行SERVER
}
