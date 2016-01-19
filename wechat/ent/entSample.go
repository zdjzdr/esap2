/**
 * 企业号API实例.
 * @woylin, 2016-1-6
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/woylyn/esap2/db/sqlsrv"
	"github.com/woylyn/esap2/wechat"
)

var (
	token          = "esap"
	corpId         = ""
	encodingAesKey = ""
	agentMap       = make(map[int]WxAgenter)
	port           = ":8080"
)

//微信“应用接口”，实现这些接口函数可被API主进程引导调用
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

//实现微信“应用接口”的父应用，定义应用时应继承
type WxAgent struct {
	req  *wechat.WxReq
	resp []byte
}

//实现接口的函数， 默认为空方法
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
	//注册应用分支
	agentMap[1] = &Agent1{}   //改进
	agentMap[2] = &Agent2{}   //备件
	agentMap[3] = &Agent3{}   //ESAP
	agentMap[6] = &Agent6{}   //采集
	agentMap[7] = &Agent7{}   //订单
	agentMap[10] = &Agent10{} //考勤
	agentMap[11] = &Agent11{} //报表
	agentMap[13] = &Agent13{} //报表
	agentMap[15] = &Agent15{} //报表
	//设置管理员密钥
	wechat.SetSecret("qcVynIzDfsTEbmRINrNXQZ4KGAntPWZ_ZxYvlq")
	//设置token,corpId,encodingAesKey
	wechat.SetBiz(token, corpId, encodingAesKey)
	//并发线程定期获取AccessToken
	go wechat.FetchCorpAccessToken2()
	//并发线程定期检查微信提醒通知
	go checkWxtx()
}

func main() {
	http.HandleFunc("/", wxHander)
	http.HandleFunc("/notice", notceHander)
	http.HandleFunc("/gj", noHander)
	http.HandleFunc("/dbsy", noHander)

	log.Println("Wechat: Started at", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Wechat: Running failed, ", err)
	}
	log.Println("Wechat: Stop!")
}

//API主进程
func wxHander(w http.ResponseWriter, r *http.Request) {
	//实例化企业号应用
	wb := wechat.WxBiz{}
	//处理GET，POST请求
	switch r.Method {
	case "GET":
		_, err := wb.Vurl(r) //主要用于首次认证
		if err != nil {
			fmt.Fprintf(w, `<a href="http://m.ickd.cn" target="_blank">快递查询</a>`) //用于快递查询友链，非必需
		}
		fmt.Fprintf(w, wb.Echostr)
	case "POST":
		wr, err := wb.Vurl(r) //认证是否来自微信
		if err != nil {
			return
		}
		//查找已注册的应用，未找到则提示该应用未实现
		agent, ok := agentMap[wr.AgentID]
		if !ok {
			fmt.Printf("--This Agent[%d] has not WxAgent!\n", wr.AgentID)
			//fmt.Printf("--This WxReq is: %+v\n", wr)
			return
		}
		//传递微信请求到应用
		agent.SetReq(wr)
		//根据微信请求类型（MsgType），调用应用接口进行处理
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
		//打印直接回复的内容
		fmt.Println("--respEnc\n", string(agent.GetResp()))
		//将直接回复内容加密发送
		resp, _ := wb.EncryptMsg(agent.GetResp())
		agent.CleanResp()
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, string(resp))
	}
}
func notceHander(w http.ResponseWriter, r *http.Request) {
	//公告子页面，待实现
	fmt.Fprintf(w, "暂无公告")
}

func noHander(w http.ResponseWriter, r *http.Request) {
	//未开发功能的统一跳转页面
	fmt.Fprintf(w, "此功能尚未开放")
}

//循环扫描微信提醒，在main中go一下即可^_^

//微信提醒
type wxtx struct {
	ToUser  string
	ToAgent int
	Context string
	Id      int
}

func checkWxtx() {
	for {
		log.Println("Scanning msg to send")
		arr := sqlsrv.FetchAllRowsPtr("select touser,toagent,context,id from wxtx where isnull(flag,0)=0", &wxtx{})
		if len(*arr) != 0 {
			for _, v := range *arr {
				if v1, ok := v.(wxtx); ok {
					s := fmt.Sprintf("【新待办通知】\n描述：%v\n", v1.Context)
					fmt.Printf("--msg to send:%v", s)
					bd, _ := wechat.TextMsg(v1.ToUser, s, v1.ToAgent)
					wechat.SendMsg(bd)
					sqlsrv.Exec("update wxtx set flag=1 where id=?", v1.Id)
				}
			}
		}
		time.Sleep(time.Minute * 5)
	}
}
