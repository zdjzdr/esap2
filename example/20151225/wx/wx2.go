package main

import (
	_ "bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/woylyn/esap2/db/sqlsrv" //MSSQL,编译前先get
	"github.com/woylyn/esap2/wechat"    //微信SDK
)

var (
	empMap map[string]employee //员工数组，临时存放员工信息，当完成上传图片和姓名工号采集流程后销毁
)

const (
	TOKEN     = "esap1"    //微信
	PORT      = ":80"      //web监听端口，默认80
	DEV       = true       //开发模式标记
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
	//尝试打开已有照片，文件不存在则创建
	fh, err := os.OpenFile(fn, 2, 0777)
	if err != nil {
		fh, _ = os.Create(fn)
	}
	defer fh.Close()
	jpeg.Encode(fh, m, nil) //以JPG格式保存照片
}

//HTTP处理函数
func weixinReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, "| URL:", r.URL.String()) //控制台打印请求信息
	r.ParseForm()                                   //解析请求
	if !wechat.Valid(w, r) {                        //验证是否为微信请求
		return
	}
	switch r.Method {
	//微信服务器一般发送POST请求，这里只对POST请求进行处理
	case "POST":
		var (
			respStr  string //回复消息
			respByte []byte //回复消息字，用于转换xml
			err      error  //错误消息
		)
		wxReq := wechat.ParseWxReq(r) //解析微信请求
		if wxReq != nil {
			//fmt.Println("ReqXml:", *wxReq)
			switch wxReq.MsgType {
			case "text":
				//1.检验是否已上传图片，上传过则提示输入“姓名，工号”
				if v, ok := empMap[wxReq.FromUserName]; ok {
					//解析姓名工号
					empInfo := strings.Split(wxReq.Content, "，")
					v.Name, v.Eid = empInfo[0], empInfo[1]
					//通过员工信息表查找rcid,未找到则提示用户。
					v.Rcid = sqlsrv.Fetch("select Excelserverrcid from employee where name=? and eid=?", v.Name, v.Eid)
					if (v.Rcid) == nil {
						respStr = "未找到该用户，请重新输入或重新上传图片。"
					} else {
						//设置照片名为：前缀 + 工号
						picName := PICPREFIX + v.Eid
						//fmt.Println("picName:", picName)
						//删除原有数据库照片路径记录
						sqlsrv.Exec("delete from es_casepic where rcid=?", v.Rcid)
						//向数据库插入新照片路径记录
						err = sqlsrv.Exec("insert es_casepic(rcid,picNo,fileType,rtfid,sh,r,c,saveinto,nfsfolderid,nfsfolder,relafolder,phyfileName) values(?,?,?,?,?,?,?,?,?,?,?,?)",
							v.Rcid, picName, ".jpg", 60, 1, 3, 4, 1, 1, `ed\esdisk`, `wx\hr`, picName+".jpg")
						if err != nil {
							respStr = fmt.Sprintf("%v", "图片上传失败")
						} else {
							//上传图片到网盘，销毁员工数组中的信息，回复处理成功信息
							v.download(v.Photo)
							delete(empMap, wxReq.FromUserName)
							respStr = fmt.Sprintf("%v", "员工照片已成功处理")
						}
					}
					//生成微信text消息回复
					respByte, err = wechat.MakeTextResp(wxReq.ToUserName, wxReq.FromUserName, respStr)
					if err != nil {
						log.Println("Wechat: makeTextResp error: ", err)
						return
					}
				} else {
					//1.未上传图片其他消息的处理，仅示例
					switch wxReq.Content {
					case "蛋蛋":
						respByte, err = wechat.MakeArticleResp(wxReq.ToUserName, wxReq.FromUserName)
						if err != nil {
							log.Println("Wechat: makeTextResp error: ", err)
							return
						}
					default:
						sqlsrv.Exec("insert cxls(cDate,uid,keyword) values(?,?,?)", wxReq.CreateTime, wxReq.FromUserName, wxReq.Content)
						resp := sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", wxReq.Content))
						if *resp == nil {
							resp = sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", "默认"))
							if *resp == nil {
								*resp = "hi，guys,welcome to ESAPbuluo/::D"
							}
						}
						respByte, err = wechat.MakeTextResp(wxReq.ToUserName, wxReq.FromUserName, fmt.Sprintf("%v", *resp))
						if err != nil {
							log.Println("Wechat: makeTextResp error: ", err)
							return
						}
					}
				}
			case "image":
				empMap[wxReq.FromUserName] = employee{"", "", wxReq.PicUrl, nil}
				respByte, err = wechat.MakeTextResp(wxReq.ToUserName, wxReq.FromUserName, fmt.Sprintf("%v", "请输入员工信息(格式 姓名,工号)："))
				if err != nil {
					log.Println("Wechat: makeTextResp error: ", err)
					return
				}
			case "event":
				respByte, err = wechat.MakeTextResp(wxReq.ToUserName, wxReq.FromUserName, fmt.Sprintf("%v", "欢迎关注ESAP部落"))
				if err != nil {
					log.Println("Wechat: makeTextResp error: ", err)
					return
				}
			}
		}
		w.Header().Set("Content-Type", "text/xml")
		if DEV {
			fmt.Println(string(respByte))
		}
		fmt.Fprintf(w, string(respByte))
	}
}

//主函数，监听常量PORT所定义的端口
func main() {
	empMap = make(map[string]employee)
	log.Println("Wechat: Start at:", PORT)
	http.HandleFunc("/", weixinReq)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatal("Wechat: ListenAndServe failed, ", err)
	}
	log.Println("Wechat: Stop!")
}
