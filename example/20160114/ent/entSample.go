/**
 * 企业号SDK示例.
 * @woylin, 2016-1-6
 */
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/woylyn/esap2/wechat"
)

var (
	token          = "esap"
	corpId         = "yourCorpId"
	encodingAesKey = "yourEncodingAesKey"
	agentMap       = make(map[int]WxAgenter)
	port           = ":80"
)

type WxAgenter interface {
	Gtext()
	Gimage()
	Gvoice()
	Gshortvideo()
	Gvideo()
	Glocation()
	Gevent()
	GetResp() []byte
	CleanResp()
	SetReq(*wechat.WxReq)
}

type WxAgent struct {
	req  *wechat.WxReq
	resp []byte
}

func (w *WxAgent) Gtext()       {}
func (w *WxAgent) Gimage()      {}
func (w *WxAgent) Gvoice()      {}
func (w *WxAgent) Gshortvideo() {}
func (w *WxAgent) Gvideo()      {}
func (w *WxAgent) Glocation()   {}
func (w *WxAgent) Gevent()      {}
func (w *WxAgent) CleanResp() {
	w.resp = nil
}
func (w *WxAgent) GetResp() []byte {
	return w.resp
}
func (w *WxAgent) SetReq(req *wechat.WxReq) {
	w.req = req
}

func init() {
	agentMap[1] = &Agent1{}
	agentMap[3] = &Agent3{}
	wechat.SetSecret("yourSecret")
	wechat.SetBiz(token, corpId, encodingAesKey)
	go wechat.FetchCorpAccessToken2()
}

func main() {
	http.HandleFunc("/", wxHander)
	http.HandleFunc("/notice", notceHander)
	http.HandleFunc("/gj", gjHander)

	log.Println("Wechat: Started at", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Wechat: Running failed, ", err)
	}
	log.Println("Wechat: Stop!")
}

func wxHander(w http.ResponseWriter, r *http.Request) {
	wb := wechat.WxBiz{}

	switch r.Method {
	case "GET":
		wb.Vurl(r) //主要用于首次认证
		fmt.Fprintf(w, wb.Echostr)
	case "POST":
		wr, err := wb.Vurl(r) //认证是否来自微信
		if err != nil {
			return
		}
		agent, ok := agentMap[wr.AgentID]
		if !ok {
			fmt.Printf("--This Agent[%d] has not WxAgent!\n", wr.AgentID)
			//			fmt.Fprintf(w, "")
			return
		}

		agent.SetReq(wr)

		switch wr.MsgType {
		case "text":
			agent.Gtext()
		case "image":
			agent.Gimage()
		case "voice":
			agent.Gvoice()
		case "shortvideo":
			agent.Gshortvideo()
		case "video":
			agent.Gvideo()
		case "location":
			agent.Glocation()
		case "event":
			agent.Gevent()
		}
		fmt.Println("--respEnc\n", string(agent.GetResp()))

		resp, _ := wb.EncryptMsg(agent.GetResp())
		agent.CleanResp()
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, string(resp))

	}
}
func notceHander(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "暂无公告")

}
func gjHander(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "此功能尚未开放")

}
