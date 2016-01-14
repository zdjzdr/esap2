package main

import (
	"fmt"
	//	"time"

	"github.com/woylyn/esap2/db/sqlsrv"
	"github.com/woylyn/esap2/wechat"
)

//改进项目
type Agent1 struct {
	WxAgent
}

func (w *Agent1) Gtext() {
	//回复文本
	bd, _ := wechat.TextMsg(w.req.FromUserName, w.req.Content, w.req.AgentID)
	for i := 0; i < 3; i++ {
		go wechat.SendMsg(bd)
	}
}
func (w *Agent1) Gimage() {
	w.resp, _ = wechat.RespImg(w.req.ToUserName, w.req.FromUserName, w.req.MediaId)
}
func (w *Agent1) Gvoice() {
	w.resp, _ = wechat.RespVoice(w.req.ToUserName, w.req.FromUserName, w.req.MediaId)
}
func (w *Agent1) Gshortvideo() {
	w.resp, _ = wechat.RespVideo(w.req.ToUserName, w.req.FromUserName, w.req.MediaId, "看一看", "瞧一瞧")
}
func (w *Agent1) Gvideo() {
	w.resp, _ = wechat.RespVideo(w.req.ToUserName, w.req.FromUserName, w.req.MediaId, "看一看", "瞧一瞧")
}
func (w *Agent1) Glocation() {
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "本次签到地点："+w.req.Label)
}
func (w *Agent1) Gevent() {
	switch w.req.Event {
	case "view":
	case "click":
		switch w.req.EventKey {
		case "jxzxm":
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询进行中项目...")
			bd, _ := wechat.TextMsg(w.req.FromUserName, "未找到项目...", w.req.AgentID)
			go wechat.SendMsg(bd)
		case "ywcxm":
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询已完成项目...")
			ywcxm(w.req.FromUserName, w.req.AgentID)
		}
	}
}
func ywcxm(user string, id int) {
	//	time.Sleep(time.Second * 3)
	//	arr := sqlsrv.Fetch(fmt.Sprintf("select resp from wxr1 where charindex('%s',keyword)>0", "默认"))
	arr := sqlsrv.FetchAll(fmt.Sprintf("select username,userid from Es_user where userlogin='%s'", "w"))
	//	arr := sqlsrv.FetchOne("select * from Es_user")
	fmt.Println("--arr", *arr)

	//	arr := sqlsrv.FetchAll("select 年,季,授权任务 from [改进项目记录_主表] where 验收='未通过'")
	bd, _ := wechat.TextMsg(user, "未找到项目...", id)
	//	if len(*arr) != 0 {
	if (*arr) != nil {
		bd, _ = wechat.TextMsg(user, fmt.Sprintf("%v", *arr), id)
	}
	wechat.SendMsg(bd)
}

//ESAP
type Agent3 struct {
	WxAgent
}

func (w *Agent3) Gtext() {
	//回复文本
	bd, _ := wechat.TextMsg(w.req.FromUserName, w.req.Content, w.req.AgentID)
	for i := 0; i < 3; i++ {
		go wechat.SendMsg(bd)
	}
}
func (w *Agent3) Gimage() {
	w.resp, _ = wechat.RespImg(w.req.ToUserName, w.req.FromUserName, w.req.MediaId)
}
func (w *Agent3) Gvoice() {
	w.resp, _ = wechat.RespVoice(w.req.ToUserName, w.req.FromUserName, w.req.MediaId)
}
func (w *Agent3) Gshortvideo() {
	w.resp, _ = wechat.RespVideo(w.req.ToUserName, w.req.FromUserName, w.req.MediaId, "看一看", "瞧一瞧")
}
func (w *Agent3) Gvideo() {
	w.resp, _ = wechat.RespVideo(w.req.ToUserName, w.req.FromUserName, w.req.MediaId, "看一看", "瞧一瞧")
}
func (w *Agent3) Glocation() {
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "本次签到地点："+w.req.Label)
}
func (w *Agent3) Gevent() {
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
	w.resp, _ = wechat.RespArt(w.req.ToUserName, w.req.FromUserName, art, art2, art3)
}
