/**
 * 企业号API示例, @woylin, 2016-1-12
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
	corpId         = "wx1d2f333555566666"
	encodingAesKey = "dCpkRHPR246zreBWKzNKCwSKBAQTmLHW6IpjptE13eB"

	port = ":80"
)

func init() {
	wechat.SetBiz(token, corpId, encodingAesKey)
}
func main() {
	log.Println("Wechat: Started at", port)
	http.HandleFunc("/", wxhander)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Wechat: Running failed, ", err)
	}
	log.Println("Wechat: Stop!")
}

func wxhander(w http.ResponseWriter, r *http.Request) {
	wb := wechat.WxBiz{}

	switch r.Method {
	case "GET":
		wb.Vurl(r) //首次认证
		fmt.Fprintf(w, wb.Echostr)
	case "POST":
		wr, err := wb.Vurl(r) //认证是否来自微信
		if err != nil {
			return
		}

		var resp1 []byte
		switch wr.MsgType {
		case "text":
			//回复文本
			resp1, _ = wechat.RespText(wr.ToUserName, wr.FromUserName, "U say:"+wr.Content)
		case "image":
			//回复图片
			resp1, _ = wechat.RespImg(wr.ToUserName, wr.FromUserName, wr.MediaId)
		case "voice":
			//回复语音
			resp1, _ = wechat.RespVoice(wr.ToUserName, wr.FromUserName, wr.MediaId)
		case "shortvideo":
			//回复视频
			resp1, _ = wechat.RespVideo(wr.ToUserName, wr.FromUserName, wr.MediaId, "看一看", "瞧一瞧")
		case "video":
			//回复视频
			resp1, _ = wechat.RespVideo(wr.ToUserName, wr.FromUserName, wr.MediaId, "看一看", "瞧一瞧")
		case "location":
			//回复文本
			resp1, _ = wechat.RespText(wr.ToUserName, wr.FromUserName, "本次签到地点："+wr.Label)
		case "event":
			//回复图文
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
			resp1, _ = wechat.RespArt(wr.ToUserName, wr.FromUserName, art, art2, art3)
		}
		fmt.Println("--respEnc\n", string(resp1))

		resp, _ := wb.EncryptMsg(resp1)
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprintf(w, string(resp))
	}
}
