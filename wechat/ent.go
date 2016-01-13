/**
 * 微信SDK-企业号接口(产生密文), @woylin, 2016-1-6
 * 企业号加解密库，主要提供URL验证，消息加解密三个接口函数
 * 目前官方未提供golang版，本实现参考了php版官方库
 */
package wechat

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	corpAccessTokenUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	corpPostMsgUrl     = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
)

//企业号AccessToken获取
func FetchCorpAccessToken() error {
	rUrl := fmt.Sprintf(corpAccessTokenUrl, appId, secret)
	resp, err := http.Get(rUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	accessTokenResp := &AccessTokenResp{}
	json.Unmarshal(body, accessTokenResp)
	if accessTokenResp.AccessToken != "" {
		AccessToken = accessTokenResp.AccessToken
	}
	log.Println("GET-AccessToken:", *accessTokenResp)
	return nil
}

//企业号循环获取，在main中go一下即可^_^
func FetchCorpAccessToken2() {
	for {
		err := FetchCorpAccessToken()
		if err != nil {
			log.Println("FetchAccessTokenError:", err)
			continue
		}
		time.Sleep(FetchDelay)
	}
}

//初始化调用，设置token,corpid和aesKey
func SetBiz(t, c, k string) error {
	if len(k) != 43 {
		return errors.New("ErrorCode: IllegalAesKey")
	}
	token = t //注意token在base包中定义
	appId = c
	aesKey = AesKeyDecode(k)
	return nil
}

//微信加密请求体
type WxEncReq struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	AgentID    string
	Encrypt    string
}

//微信加密回复体
type WxEncResp struct {
	XMLName      xml.Name `xml:"xml"`
	Encrypt      CDATA
	MsgSignature CDATA
	TimeStamp    string
	Nonce        CDATA
}

//企业客服消息体
type BizMsg struct {
	ToUser  string     `json:"touser"`
	ToParty string     `json:"toparty"`
	ToTag   string     `json:"totag"`
	MsgType string     `json:"msgtype"`
	AgentId int        `json:"agentid"`
	Text    content    `json:"text"`
	Image   media      `json:"image"`
	Voice   media      `json:"voice"`
	Video   video      `json:"video"`
	File    media      `json:"file"`
	News    articles   `json:"news"`
	MpNews  mparticles `json:"mpnews"`
	Safe    string     `json:"safe"`
}

//创建文本客服消息
func TextMsg(toUser, msg string, agentId int) ([]byte, error) {
	csMsg := &BizMsg{
		ToUser:  toUser,
		MsgType: "text",
		AgentId: agentId,
		Text:    content{Content: msg},
		Safe:    "0",
	}
	return json.MarshalIndent(csMsg, " ", "  ")
}

func SendMsg(body []byte) error {
	rUrl := strings.Join([]string{corpPostMsgUrl, AccessToken}, "")
	postReq, err := http.NewRequest("POST", rUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}
	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")
	client := &http.Client{}
	_, err = client.Do(postReq)
	if err != nil {
		return err
	}
	return nil
}

//企业消息-文本
type content struct {
	Content string `json:"content"`
}

//企业消息-媒体
type media struct {
	MediaId string `json:"media_id"`
}

//企业消息-视频
type video struct {
	MediaId     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

//企业消息-图文组
type articles struct {
	Articles []article `json:"articles"`
}

//企业消息-密图文组
type mparticles struct {
	Articles []mparticle `json:"articles"`
}

//企业消息-单图文
type article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	PicUrl      string `json:"picurl"`
}

//企业消息-密图文
type mparticle struct {
	Title        string `json:"title"`
	ThumbMediaId string `json:"thumb_media_id`
	Author       string `json:"author"`
	Url          string `json:"content_source_url"`
	Content      string `json:"content"`
	Digest       string `json:"digest"`
	ShowCoverPic string `json:"show_cover_pic"`
}

//微信消息处理载体
type WxBiz struct {
	timestamp    string
	nonce        string
	msgSignature string
	Echostr      string
}

//验证URL,验证成功则返回标准请求体（已解密）
func (w *WxBiz) Vurl(r *http.Request) (*WxReq, error) {
	log.Println(r.Method, "|", r.URL.String())
	//解析请求
	w.timestamp = r.FormValue("timestamp")
	w.nonce = r.Form.Get("nonce")
	w.msgSignature = r.Form.Get("msg_signature")
	w.Echostr = r.Form.Get("echostr")
	if r.Method == "POST" {
		w.Echostr = parseEncReq(r).Encrypt //POST请求需解析消息体中的Encrpt
	}
	//验证signature
	signature := getSHA1(token, w.timestamp, w.nonce, w.Echostr)
	if signature != w.msgSignature {
		fmt.Println("--w\n", w, "\n", signature)
		return nil, errors.New("ErrorCode: ValidateSignatureError")
	}
	w.Echostr, _ = w.DecryptMsg(w.Echostr)
	fmt.Println("--Req:\n", w.Echostr)
	wxreq := &WxReq{}
	xml.Unmarshal([]byte(w.Echostr), wxreq)
	return wxreq, nil
}

//将普通进行AES-CBC加密,打包成xml格式回复
func (w *WxBiz) EncryptMsg(resp []byte) ([]byte, error) {
	encXmlData := respEnc(resp)
	encResp := &WxEncResp{}
	encResp.Encrypt = cCDATA(encXmlData)
	encResp.MsgSignature = cCDATA(getSHA1(token, w.timestamp, w.nonce, encXmlData))
	encResp.TimeStamp = w.timestamp
	encResp.Nonce = cCDATA(w.nonce)
	return xml.MarshalIndent(encResp, " ", "  ")
}

// 解密消息,密文->base64Dec->aesDec->去除头部随机字串
func (w *WxBiz) DecryptMsg(s string) (string, error) {
	aesMsg, _ := base64.StdEncoding.DecodeString(s)
	deMsg, _ := AesDecrypt(aesMsg, aesKey)
	buf := bytes.NewBuffer(deMsg[16:20])
	var length int32
	binary.Read(buf, binary.BigEndian, &length)
	idstart := 20 + length
	id := deMsg[idstart : int(idstart)+len(appId)]
	if string(id) != appId {
		return "", errors.New("Appid is invalid")
	}
	rs := string(deMsg[20 : 20+length])
	return rs, nil
}

//转化普通回复为加密回复体，[]byte->string
func respEnc(body []byte) string {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int32(len(body)))
	if err != nil {
		fmt.Println("Binary write err:", err)
	}
	bodyLength := buf.Bytes()
	randomBytes := []byte("abcdefghijklmnop")

	plainData := bytes.Join([][]byte{randomBytes, bodyLength, body, []byte(appId)}, nil)
	encBody, _ := AesEncrypt(plainData, aesKey)
	return base64.StdEncoding.EncodeToString(encBody)
}

//解析微信加密请求
func parseEncReq(r *http.Request) *WxEncReq {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	reqBody := &WxEncReq{}
	xml.Unmarshal(body, reqBody)
	return reqBody
}
