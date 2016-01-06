/**
 * 对公众平台发送给公众账号的消息加解密示例代码.
 *
 * @woylin, 2016-1-6
 */
package main

import (
	"encoding/base64"
	//	"crypto/aes"
	"crypto/sha1"
	"encoding/xml"
	"io/ioutil"
	//	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

var (
	encodingAesKey = "ZzU9SbVE5UMdzupn2La1DthWp3VEqxyB8xsRp6o8qpM"
	token          = "esap"
	corpId         = "wx1d2f333568746602"
	port           = ":80"
	aesKey         []byte
)

type EncryptRequestBody struct {
	XMLName    xml.Name `xml:"xml"`
	ToUserName string
	Encrypt    string
}

func parseEncryptRequestBody(r *http.Request) *EncryptRequestBody {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	requestBody := &EncryptRequestBody{}
	xml.Unmarshal(body, requestBody)
	return requestBody
}

func parseEncryptRequestBody2(s string) *EncryptRequestBody {
	requestBody := &EncryptRequestBody{}
	xml.Unmarshal([]byte(s), requestBody)
	return requestBody
}

func encodingAESKey2AESKey(encodingKey string) []byte {
	data, _ := base64.StdEncoding.DecodeString(encodingKey + "=")
	return data
}

func wxhander(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "|", r.URL.String())
	r.ParseForm()
	timestamp := strings.Join(r.Form["timestamp"], "")
	nonce := strings.Join(r.Form["nonce"], "")
	//	signature := strings.Join(r.Form["signature"], "")
	//	encryptType := strings.Join(r.Form["encrypt_type"], "")
	msgSignature := strings.Join(r.Form["msg_signature"], "")
	echostr := strings.Join(r.Form["echostr"], "")
	wxbiz := WXBizMsgCrypt{token, encodingAesKey, corpId}
	rep := ""
	switch r.Method {
	case "GET":
		//首次验证

		err := wxbiz.VerifyURL(msgSignature, timestamp, nonce, echostr, &rep)
		if err != nil {
			fmt.Println(err)
			return
		}
		log.Println("Wechat Service: msg_signature validation is ok!")

		fmt.Fprintf(w, rep)
	case "POST":

		if msgSignature != "" {
			log.Println("Wechat Service: in safe mode")
			//			encryptRequestBody := parseEncryptRequestBody(r)

			//Validate msg signature
			err := wxbiz.VerifyURL(msgSignature, timestamp, nonce, echostr, &rep)
			if err != nil {
				fmt.Println(err)
				return
			}
			log.Println("Wechat Service: msg_signature validation is ok!")
		}
	}
}

func main() {
	aesKey = encodingAESKey2AESKey(encodingAesKey)
	fmt.Println("aesK:", aesKey, string(aesKey))
	log.Println("Wechat: Started")
	http.HandleFunc("/", wxhander)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Wechat: ListenAndServe failed, ", err)
	}
	log.Println("Wechat: Stop!")

}

/**
 * 1.第三方回复加密消息给公众平台；
 * 2.第三方收到公众平台发送的消息，验证消息的安全性，并对消息进行解密。
 */
type WXBizMsgCrypt struct {
	Token          string
	EncodingAesKey string
	Corpid         string
}

/*
	*验证URL
    *@param MsgSignature: 签名串，对应URL参数的msg_signature
    *@param TimeStamp: 时间戳，对应URL参数的timestamp
    *@param Nonce: 随机串，对应URL参数的nonce
    *@param EchoStr: 随机串，对应URL参数的echostr
    *@param ReplyEchoStr: 解密之后的echostr，当return返回0时有效
    *@return：失败返回对应的错误码
*/
func (w WXBizMsgCrypt) VerifyURL(MsgSignature, TimeStamp, Nonce, EchoStr string, ReplyEchoStr *string) error {
	//	public function VerifyURL($sMsgSignature, $sTimeStamp, $sNonce, $sEchoStr, &$sReplyEchoStr)
	//检验key有效性
	//	str1, _ := aes.(w.EncodingAesKey)
	//	fmt.Println("AesKey", string(str1))
	if len(w.EncodingAesKey) != 43 {
		return errors.New("ErrorCode: IllegalAesKey")
	}
	//	pc, _ := aes.NewCipher([]byte(w.EncodingAesKey))
	signature := getSha1(w.Token, TimeStamp, Nonce, EchoStr)
	if signature != MsgSignature {
		return errors.New("ErrorCode: ValidateSignatureError")
	}
	encryptXml := parseEncryptRequestBody2(EchoStr)
	fmt.Println("encryptXml:", encryptXml)
	//	pc.Decrypt([]byte(EchoStr), []byte(w.EncodingAesKey))

	return nil
}

func getSha1(token, timestamp, nonce, encrypt_msg string) string {
	//	fmt.Println("getSha1:", sha1(sort(token, timestamp, nonce, encrypt_msg)))

	sl := []string{token, timestamp, nonce, encrypt_msg}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

/**
 * 将公众平台回复用户的消息加密打包.
 * <ol>
 *    <li>对要发送的消息进行AES-CBC加密</li>
 *    <li>生成安全签名</li>
 *    <li>将消息密文和安全签名打包成xml格式</li>
 * </ol>
 *
 * @param $replyMsg string 公众平台待回复用户的消息，xml格式的字符串
 * @param $timeStamp string 时间戳，可以自己生成，也可以用URL参数的timestamp
 * @param $nonce string 随机串，可以自己生成，也可以用URL参数的nonce
 * @param &$encryptMsg string 加密后的可以直接回复用户的密文，包括msg_signature, timestamp, nonce, encrypt的xml格式的字符串,
 *                      当return返回0时有效
 *
 * @return int 成功0，失败返回对应的错误码
 */

//	public function EncryptMsg($sReplyMsg, $sTimeStamp, $sNonce, &$sEncryptMsg)
//	{
//		$pc = new Prpcrypt($this->m_sEncodingAesKey);

//		//加密
//		$array = $pc->encrypt($sReplyMsg, $this->m_sCorpid);
//		$ret = $array[0];
//		if ($ret != 0) {
//			return $ret;
//		}

//		if ($sTimeStamp == null) {
//			$sTimeStamp = time();
//		}
//		$encrypt = $array[1];

//		//生成安全签名
//		$sha1 = new SHA1;
//		$array = $sha1->getSHA1($this->m_sToken, $sTimeStamp, $sNonce, $encrypt);
//		$ret = $array[0];
//		if ($ret != 0) {
//			return $ret;
//		}
//		$signature = $array[1];

//		//生成发送的xml
//		$xmlparse = new XMLParse;
//		$sEncryptMsg = $xmlparse->generate($encrypt, $signature, $sTimeStamp, $sNonce);
//		return ErrorCode::$OK;
//	}

//	/**
//	 * 检验消息的真实性，并且获取解密后的明文.
//	 * <ol>
//	 *    <li>利用收到的密文生成安全签名，进行签名验证</li>
//	 *    <li>若验证通过，则提取xml中的加密消息</li>
//	 *    <li>对消息进行解密</li>
//	 * </ol>
//	 *
//	 * @param $msgSignature string 签名串，对应URL参数的msg_signature
//	 * @param $timestamp string 时间戳 对应URL参数的timestamp
//	 * @param $nonce string 随机串，对应URL参数的nonce
//	 * @param $postData string 密文，对应POST请求的数据
//	 * @param &$msg string 解密后的原文，当return返回0时有效
//	 *
//	 * @return int 成功0，失败返回对应的错误码
//	 */
func (w WXBizMsgCrypt) DecryptMsg(sMsgSignature, sTimeStamp, sNonce, sPostData string, sMsg *string) error {
	//	public function DecryptMsg($sMsgSignature, $sTimeStamp = null, $sNonce, $sPostData, &$sMsg)
	//	{
	if len(w.EncodingAesKey) != 43 {
		return errors.New("ErrorCode: IllegalAesKey")
	}

	//		$pc = new Prpcrypt($this->m_sEncodingAesKey);

	//		//提取密文
	//		$xmlparse = new XMLParse;
	//		$array = $xmlparse->extract($sPostData);
	//		$ret = $array[0];

	//		if ($ret != 0) {
	//			return $ret;
	//		}

	//		if ($sTimeStamp == null) {
	//			$sTimeStamp = time();
	//		}

	//		$encrypt = $array[1];
	//		$touser_name = $array[2];

	//		//验证安全签名
	//		$sha1 = new SHA1;
	//		$array = $sha1->getSHA1($this->m_sToken, $sTimeStamp, $sNonce, $encrypt);
	//		$ret = $array[0];

	//		if ($ret != 0) {
	//			return $ret;
	//		}

	//		$signature = $array[1];
	//		if ($signature != $sMsgSignature) {
	//			return ErrorCode::$ValidateSignatureError;
	//		}

	//		$result = $pc->decrypt($encrypt, $this->m_sCorpid);
	//		if ($result[0] != 0) {
	//			return $result[0];
	//		}
	//		$sMsg = $result[1];

	//		return ErrorCode::$OK;
	//	}
	return nil
}
