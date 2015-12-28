/* 微信SDK包-内置server
 * by woylin 2015/12/24
 */
package wechat

import (
	"fmt"
	"log"
	"net/http"
)

var (
	port = ":80"
	dev  = true
)

//内置微信Server接口，实现这些方法则可调用Run()启动运行内置ApiServer
type WxTexter interface {
	GoText() //文字类消息接口
}
type WxImager interface {
	GoImage() //图片类消息接口
}
type WxEventer interface {
	GoEvent() //事件类消息接口
}

type WxApp struct {
	Name    string //公众号名称
	RespStr string //回复消息
	RespB   []byte //回复消息字，用于转换xml
	Req     *WxReq
	Fn      interface{} //方法重写接口
}

func (w *WxApp) GoText()  {}
func (w *WxApp) GoImage() {}
func (w *WxApp) GoEvent() {
	if w.Req.Event == "subscribe" {
		w.RespB, _ = RespText(w.Req.ToUserName, w.Req.FromUserName, fmt.Sprintf("%v", "欢迎关注"+w.Name))
	}
}
func (w *WxApp) ParseReq(r *http.Request) {
	w.Req = ParseWxReq(r)
}

//内置Server的HTTP默认处理函数
func (w *WxApp) WxHandler(w2 http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, "| URL:", r.URL.String()) //控制台打印请求信息
	r.ParseForm()                                   //解析请求
	if !Valid(w2, r) {                              //验证是否为微信请求
		return
	}
	switch r.Method {

	//微信服务器一般发送POST请求到API，这里只对POST请求进行处理
	case "POST":
		w.ParseReq(r) //解析微信请求
		if w.Req != nil {
			if dev {
				fmt.Printf("ReqXml:%#v\n", w.Req)
			}

			switch w.Req.MsgType {
			case "text":
				if v, ok := w.Fn.(WxTexter); ok {
					v.GoText()
				}
			case "image":
				if v, ok := w.Fn.(WxImager); ok {
					v.GoImage()
				}
			case "event":
				if v, ok := w.Fn.(WxEventer); ok {
					v.GoEvent()
				}
			}
		}
		w2.Header().Set("Content-Type", "text/xml")
		if dev {
			fmt.Println(string(w.RespB))
		}
		fmt.Fprintf(w2, string(w.RespB))
	}
}

//运行内置Server
func (w *WxApp) Run(app interface{}) {
	w.Fn = app
	log.Println("Wechat: Start at", port)
	http.HandleFunc("/", w.WxHandler)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Wechat: ListenAndServe failed, ", err)
	}
	log.Println("Wechat: Stop!")
}

//设置内置Server端口
func SetPort(portNo string) {
	port = portNo
}

//设置内置开发模式
func SetDev(d bool) {
	dev = d
}
