/**
 * 微信SDK-基础接口(产生明文), @woylin, 2015/12/24
 */
package wechat

import (
	"crypto/sha1"
	"encoding/json"
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

var (
	token               = "esap" //默认token
	appId               string   //企业号为corpId
	secret              string
	aesKey              []byte //解密的AesKey
	accessTokenFetchUrl = "https://api.weixin.qq.com/cgi-bin/token"
	AccessToken         = ""
	FetchDelay          = time.Minute * 5 //默认5分钟获取一次
)

func SetToken(t string) {
	token = t
}
func SetAppId(a string) {
	appId = a
}
func SetSecret(s string) {
	secret = s
}

//AccessToken回复体
type AccessTokenResp struct {
	AccessToken string  `json:"access_token"`
	ExpiresIn   float64 `json:"expires_in"`
	Errcode     float64
	Errmsg      string
}

//获取AccessToken
func FetchAccessToken() (string, float64, error) {
	requestLine := strings.Join([]string{accessTokenFetchUrl,
		"?grant_type=client_credential&appid=",
		appId,
		"&secret=",
		secret}, "")

	resp, err := http.Get(requestLine)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", 0.0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0.0, err
	}

	accessTokenResp := &AccessTokenResp{}
	json.Unmarshal(body, accessTokenResp)

	fmt.Println(accessTokenResp)
	if accessTokenResp.AccessToken != "" {
		AccessToken = accessTokenResp.AccessToken
	}
	return "", 0.0, err
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
	Location_X   float32 `xml:"Latitude"`  //location
	Location_Y   float32 `xml:"Longitude"` //location
	Precision    float32 //LOCATION
	Scale        int     //location
	Label        string  //location
	MsgId        int
	Event        string //event
	EventKey     string //event
	ScanCodeInfo ScanInfo
	AgentID      int //corp
}

//微信回复体
type WxResp struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA
	FromUserName CDATA
	CreateTime   time.Duration
	MsgType      CDATA
	Content      CDATA //type:text
	Image        Media //type:image
	Voice        Media //type:voice
	Format       CDATA //type:voice
	Video        Video //type:video
	ArticleCount int   //type:news
	Articles     item  //type:news
}

//图片声音
type Media struct {
	MediaId CDATA
}

//视频
type Video struct {
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
	Text string `xml:",innerxml"`
}

//文本转CDATA
func cCDATA(v string) CDATA {
	return CDATA{"<![CDATA[" + v + "]]>"}
}

//验证微信请求
func Valid(w http.ResponseWriter, r *http.Request) bool {
	timestamp := r.FormValue("timestamp")
	nonce := r.Form.Get("nonce")
	signature := r.Form.Get("signature")

	tmpStr := getSHA1(timestamp, nonce, token)
	if tmpStr == signature {
		echostr := r.Form.Get("echostr")
		fmt.Fprintf(w, echostr)
		return true
	}

	log.Println("Wechat: Request is not from Wechat platform!")
	return false
}

//排序并sha1，用于计算signature
func getSHA1(sl ...string) string {
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
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
	wxResp.Video = Video{cCDATA(mediaId), cCDATA(title), cCDATA(desc)}
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
