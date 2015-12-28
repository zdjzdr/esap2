/* 微信SDK包-基础接口
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
	Content      CDATA //type:text
	Image        Image //type:image
	ArticleCount int   //type:news
	Articles     item  //type:news
}

//图片
type Image struct {
	MeidaId CDATA
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
func CreArt(title, desc, picUrl, url string) Article {
	art := Article{cCDATA(title),
		cCDATA(desc),
		cCDATA(picUrl),
		cCDATA(url)}
	return art
}

var token = "esap" //默认token

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
