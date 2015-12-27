/* 微信SDK基础库
 * by woylin 2015/12/24
 */
package wechat

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

//微信请求体
type WxReq struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string  //text
	PicUrl       string  //image
	MediaId      string  //image
	Location_X   float32 //location
	Location_Y   float32 //location
	Scale        int     //location
	Label        string  //location
	MsgId        int
	Event        string //event
}

//微信回复体
type WxResp struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   time.Duration
	MsgType      CDATA
	Content      CDATA    //type:text
	Image        Image    //type:image
	ArticleCount int      //type:news
	Articles     articles //type:news
}

//图片
type Image struct {
	MeidaId CDATA
}

//文章组
type articles struct {
	Item article `xml:"item"`
}

//文章
type article struct {
	Title       CDATA
	Description CDATA
	PicUrl      CDATA
	Url         CDATA
}

//标准CDATA
type CDATA struct {
	//	Text []byte `xml:",innerxml"`
	Text string `xml:",innerxml"`
}

//文本转CDATA
func cCDATA(v string) CDATA {
	//return CDATA{[]byte("<![CDATA[" + v + "]]>")}
	return CDATA{"<![CDATA[" + v + "]]>"}
}

var token = "esap"

func SetToken(t string) {
	token = t
}

//验证微信请求
func Valid(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signature := strings.Join(r.Form["signature"], "")
	if checkSignature(timestamp, nonce, signature, token) {
		echostr := strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
		return true
	}
	log.Println("Wechat: Request is not from Wechat platform!")
	return false
}

func checkSignature(timestamp, nonce, signature, token string) bool {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	tmpStr := fmt.Sprintf("%x", s.Sum(nil))
	if tmpStr != signature {
		return false
	}
	return true
}

//解析微信请求，返回请求体
func ParseWxReq(r *http.Request) *WxReq {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	wxReq := &WxReq{}
	xml.Unmarshal(body, wxReq)
	return wxReq
}

//回复文本类消息
func RespText(fromUserName, toUserName, content string) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("text")
	wxResp.Content = cCDATA(content)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(wxResp, " ", "  ")
}

//回复文章类消息
func RespArt(fromUserName, toUserName string) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("news")
	wxResp.ArticleCount = 1
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	//	WxResp.Articles = article{}
	art1 := article{cCDATA("打通信息化的“任督二脉”"),
		cCDATA("来自村长的ESAP2.0系统最新技术分享。"),
		cCDATA("http://iesap.net/wp-content/uploads/2015/12/rdem.jpg"),
		cCDATA("http://iesap.net/index.php/2015/12/11/esap2-0/")}
	//	arts :=
	wxResp.Articles = articles{art1}
	return xml.MarshalIndent(wxResp, " ", "  ")
}
