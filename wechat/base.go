/**
 * 微信SDK-基础接口(产生明文), @woylin, 2015/12/24
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

var token = "esap" //默认token

func SetToken(t string) {
	token = t
}

//微信请求体
type WxReq struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string //text
	PicUrl       string //image
	MediaId      string //image/voice/video
	ThumbMediaId string
	Location_X   float32 //location
	Location_Y   float32 //location
	Scale        int     //location
	Label        string  //location
	MsgId        int
	Event        string //event
	EventKey     string //event
	ScanCodeInfo ScanInfo
	AgentID      string //corp
}

//微信回复体
type WxResp struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   time.Duration
	MsgType      CDATA
	Content      CDATA //type:text
	Image        media //type:image
	Voice        media //type:voice
	Format       CDATA //type:voice
	Video        video //type:video
	ArticleCount int   //type:news
	Articles     item  //type:news
}

//图片声音
type media struct {
	MediaId CDATA
}

//视频
type video struct {
	MediaId     CDATA
	Title       CDATA
	Description CDATA
}

//文章组
type item struct {
	Item []Article `xml:"item"`
}

//文章
type Article struct {
	Title       CDATA
	Description CDATA
	PicUrl      CDATA
	Url         CDATA
}

//扫描（二维码，条码）
type ScanInfo struct {
	ScanType   string
	ScanResult string
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

//验证微信请求
func Valid(w http.ResponseWriter, r *http.Request) bool {
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	signature := strings.Join(r.Form["signature"], "")
	msgSignature := strings.Join(r.Form["msg_signature"], "")
	//检验是否来自企业号
	if msgSignature != "" {

	}
	if CheckSignature(timestamp, nonce, signature, token) {
		echostr := strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
		return true
	}
	log.Println("Wechat: Request is not from Wechat platform!")
	return false
}

func CheckSignature(timestamp, nonce, signature, token string) bool {
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

//解析http请求，返回请求体
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

//解析字符串请求，返回请求体
func ParseWxReqS(s string) *WxReq {
	wxReq := &WxReq{}
	xml.Unmarshal([]byte(s), wxReq)
	return wxReq
}

//回复单文本[]bytes
func RespText(fromUserName, toUserName, content string) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("text")
	wxResp.Content = cCDATA(content)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(wxResp, " ", "  ")
}

//回复多图文[]bytes, art ...Article
func RespArt(fromUserName, toUserName string, art ...Article) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("news")
	wxResp.ArticleCount = len(art)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	items := item{art}
	wxResp.Articles = items
	return xml.MarshalIndent(wxResp, " ", "  ")
}

//回复图片[]bytes
func RespImg(fromUserName, toUserName, mediaId string) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("image")
	wxResp.Image.MediaId = cCDATA(mediaId)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(wxResp, " ", "  ")
}

//回复音频[]bytes
func RespVoice(fromUserName, toUserName, mediaId string) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("voice")
	wxResp.Voice.MediaId = cCDATA(mediaId)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(wxResp, " ", "  ")
}

//回复视频[]bytes
func RespVideo(fromUserName, toUserName, mediaId, title, desc string) ([]byte, error) {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("video")
	wxResp.Video = video{cCDATA(mediaId), cCDATA(title), cCDATA(desc)}
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(wxResp, " ", "  ")
}

//创建单文本,返回struct
func CreText(fromUserName, toUserName, content string) *WxResp {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("text")
	wxResp.Content = cCDATA(content)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	return wxResp
}

//创建图文消息，一般与RespArt配合使用，先创建文章，再传入RespArt
func CreArt(title, desc, picUrl, url string) Article {
	art := Article{cCDATA(title),
		cCDATA(desc),
		cCDATA(picUrl),
		cCDATA(url)}
	return art
}

//创建多图文，art ...Article,返回struct
func CreArts(fromUserName, toUserName string, art ...Article) *WxResp {
	wxResp := &WxResp{}
	wxResp.FromUserName = cCDATA(fromUserName)
	wxResp.ToUserName = cCDATA(toUserName)
	wxResp.MsgType = cCDATA("news")
	wxResp.ArticleCount = len(art)
	wxResp.CreateTime = time.Duration(time.Now().Unix())
	items := item{art}
	wxResp.Articles = items
	return wxResp
}
