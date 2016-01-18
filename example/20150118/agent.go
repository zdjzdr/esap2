package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/woylyn/esap2/db/sqlsrv"
	"github.com/woylyn/esap2/wechat"
)

//微信提醒
type wxtx struct {
	ToUser  string
	ToAgent int
	Context string
	Id      int
}

//循环扫描微信提醒，在main中go一下即可^_^
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

func (w *Agent1) Gevent() {
	switch w.req.Event {
	case "view":
	case "click":
		switch w.req.EventKey {
		case "jxzxm":
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			bd, _ := wechat.TextMsg(w.req.FromUserName, "未找到项目...", w.req.AgentID)
			go wechat.SendMsg(bd)
		case "ywcxm":
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go ywcxm(w.req.FromUserName, w.req.AgentID)
		case "wtgxm":
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go wtgxm(w.req.FromUserName, w.req.AgentID)
		case "wqtgxm":
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go wqtgxm(w.req.FromUserName, w.req.AgentID)
		}
	}
}

type Xm struct {
	Y, J, R string
}

func ywcxm(user string, id int) {
	queryAndSend(user, id, "select 年,季,授权任务 from [改进项目记录_主表] where 验收='通过' and 年=2015", &Xm{})
}
func wtgxm(user string, id int) {
	queryAndSend(user, id, "select 年,季,授权任务 from [改进项目记录_主表] where 验收='未通过' and 年=2015", &Xm{})
}
func wqtgxm(user string, id int) {
	queryAndSend(user, id, "select 年,季,授权任务 from [改进项目记录_主表] where 验收='通过' and 年<>2015 order by 年 desc,季 desc", &Xm{})
}

//备件管理
type Agent2 struct {
	WxAgent
}

func (w *Agent2) Gtext() {
	//库存查询
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询库存情况...")
	go bjkc(w.req.FromUserName, w.req.AgentID, w.req.Content)
}

func (w *Agent2) Gevent() {
	switch w.req.Event {
	case "enter_agent":
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "近期到料情况：")
		go jqdl(w.req.FromUserName, w.req.AgentID)
	}
}

type bj struct {
	GrDate time.Time
	Mdesc  string
	Rem    string
	Qty    float32
	Unit   string
}

func (c bj) String() string {
	return fmt.Sprintf("%v-%v%v\n到料 = %v %v\n", c.GrDate.Format("1/2"), c.Mdesc, c.Rem, c.Qty, c.Unit)
}

//备件-近期到料
func jqdl(user string, id int) {
	queryAndSendArr(user, id, "select grDate,mDesc,rem,sum(qty),mun from wx_wmgr where lcid2='m09' group by grDate,mDesc,rem,mun ", &bj{})
}

type bjQty struct {
	Mdesc string
	Qty   float32
}

func (c bjQty) String() string {
	return fmt.Sprintf("%v = %v\n", c.Mdesc, c.Qty)
}

type ssq struct {
	ErrorCode int       `json:"error_code"`
	Result    []lottery `json:"result"`
}
type lottery struct {
	LotteryDate   string
	LotteryQh     int
	LotteryNumber string
}

//备件-库存查询
func bjkc(user string, id int, mDesc string) {
	if mDesc == "双色球" {
		url := "http://apis.haoservice.com/lifeservice/lottery/query?id=1&date=2016-01-17&key=b6558f20e78e45be976231577ed8dbcb"
		imgResp, _ := http.Get(url) //从API服务器获取开奖信息
		defer imgResp.Body.Close()
		body, _ := ioutil.ReadAll(imgResp.Body)
		fmt.Println("body:", string(body))
		ssq1 := &ssq{}
		json.Unmarshal(body, ssq1)
		fmt.Println("ssq:", ssq1)
		bd, _ := wechat.TextMsg(user, "最新开奖日期："+ssq1.Result[0].LotteryDate+"\n本期号码："+ssq1.Result[0].LotteryNumber, id)
		go wechat.SendMsg(bd)
	}
	sql := fmt.Sprintf("select mDesc,iqty from vlbq2 where lcid='m09' and isnull(iqty,0)>0 and charindex('%s',mdesc)>0", mDesc)
	queryAndSendArr(user, id, sql, &bjQty{})
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
	switch w.req.Event {
	case "click":
		switch w.req.EventKey {
		case "xtgg":
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
	}
}

//11.核心报表
type Agent11 struct {
	WxAgent
}

func (w *Agent11) Gtext() {
	//回复文本
	bd, _ := wechat.TextMsg(w.req.FromUserName, w.req.Content, w.req.AgentID)
	for i := 0; i < 3; i++ {
		go wechat.SendMsg(bd)
	}
}

func (w *Agent11) Gevent() {
	switch w.req.Event {
	case "view":
	case "click":
		switch w.req.EventKey {
		case "zxrb": //主线日报
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go zxrb(w.req.FromUserName, w.req.AgentID)
		case "bzyb": //包装月报
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go bzyb(w.req.FromUserName, w.req.AgentID)
		case "pjyb": //配件月报
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go pjyb(w.req.FromUserName, w.req.AgentID)
		case "ddyb": //订单月报
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go ddyb(w.req.FromUserName, w.req.AgentID)
		case "sdkc": //素电库存
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go sdkc(w.req.FromUserName, w.req.AgentID)
		case "cpkc": //产品库存
			w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询...")
			go cpkc(w.req.FromUserName, w.req.AgentID)
		}
	}
}

type Cl struct {
	X产线 string
	X产量 float32
	X损耗 string
}

func (c Cl) String() string {
	return fmt.Sprintf("%v\n产量(万):%v\n损耗:%v", c.X产线, c.X产量/10000, c.X损耗)
}

func zxrb(user string, id int) {
	xm := &Cl{}
	queryMaxDate(user, id, "wx_zxclzx")

	queryAndSend(user, id, "select 车号,sum(pqty),sum(sqty) from wx_zxclzx group by 车号", xm)
	queryAndSend(user, id, "select '--总产量--',sum(pqty),sum(sqty) from wx_zxclzx", xm)
}

func bzyb(user string, id int) {
	xm := &Cl{}
	queryMaxDate(user, id, "wx_zxclbz")

	queryAndSend(user, id, "select mDesc,sum(pqty)as pqty,sum(isnull(sqty,0)) from wx_zxclbz group by mDesc order by pqty desc", xm)
	queryAndSend(user, id, "select '--总产量--',sum(pqty)as pqty,sum(isnull(sqty,0)) from wx_zxclbz", xm)
}

type sd struct {
	Spec  string
	Type  string
	Qty   float32
	Zqty  float32
	Zrate float32
}

func (c sd) String() string {
	return fmt.Sprintf("%v-%v: \n%v / %v\n暂放率: %v%%", c.Spec, c.Type, c.Qty, c.Zqty, c.Zrate)
}

func sdkc(user string, id int) {
	xm := &sd{}

	queryAndSend(user, id, "select * from wx_sdkc", xm)
	queryAndSend(user, id, "select spec,'汇总',sum(iqty)as iqty,sum(isnull(zqty,0)),round(sum(isnull(zqty,0))/sum(iqty),4)*100 from wx_sdkc group by spec", xm)
}

type cp struct {
	Mdesc string
	Qty   float32
}

func (c cp) String() string {
	return fmt.Sprintf("[%v]\n库存(万)：%v", c.Mdesc, c.Qty/10000)
}

func cpkc(user string, id int) {
	xm := &cp{}

	queryAndSend(user, id, "select mDesc,sum(iqty) as iqty from vlbq2 where lcid='m03' group by mDesc having sum(iqty) >0 order by mDesc", xm)
	queryAndSend(user, id, "select '--总库存--',sum(iqty)as iqty from vlbq2 where lcid='m03'", xm)
}

type dd struct {
	Mon  string
	Sgid string
	Qty  float32
}

func (c dd) String() string {
	return fmt.Sprintf("%v月--%v\n订单量(万):%v", c.Mon, c.Sgid, c.Qty/10000)
}

func ddyb(user string, id int) {
	xm := &dd{}

	queryAndSend(user, id, "select 月份,订单类型,sum(isnull(数量,0)) from wx_sxdd group by 月份,订单类型 order by 月份, 订单类型", xm)
	queryAndSend(user, id, "select '全年全','总订单量--',sum(数量)from wx_sxdd", xm)
}

type pjCl struct {
	Date time.Time
	X产线  string
	X产量  string
	X损耗  string
}

func (c pjCl) String() string {
	return fmt.Sprintf("%v\n产量(万):%v\n损耗:%v\n记录日期：%v", c.X产线, c.X产量, c.X损耗, c.Date.Format("2006-01-02"))
}

func pjyb(user string, id int) {
	queryAndSend(user, id, "select cdate,产线,sum(pqty)as pqty,sum(isnull(sqty,0)) as sqty from wx_zxclpj group by cdate,产线 order by 产线", &pjCl{})
}

func queryMaxDate(user string, id int, table string) {
	date := sqlsrv.Fetch(fmt.Sprintf("select max(cdate) from %s", table))
	if v, ok := (*date).(time.Time); ok {
		bd, _ := wechat.TextMsg(user, fmt.Sprintf("最新记录日期: %s", v.Format("2006-01-02")), id)
		wechat.SendMsg(bd)
	}
}

func queryAndSend(user string, id int, sql string, struc interface{}) {
	time.Sleep(time.Second)
	arr := sqlsrv.FetchAllRowsPtr(sql, struc)
	bd, _ := wechat.TextMsg(user, "未找到项目...", id)
	if (*arr) != nil {
		fmt.Println("--arr:", *arr)
		for k, v := range *arr {
			s := fmt.Sprintf("%v：%v", k+1, v)
			if len(*arr) == 1 {
				s = fmt.Sprintf("%v", v)
			}
			fmt.Printf(" %v\n", v)
			bd, _ = wechat.TextMsg(user, s, id)
			wechat.SendMsg(bd)
		}
	}
}
func queryAndSendArr(user string, id int, sql string, struc interface{}) {
	time.Sleep(time.Second)
	arr := sqlsrv.FetchAllRowsPtr(sql, struc)
	bd, _ := wechat.TextMsg(user, "未找到项目...", id)
	if (*arr) != nil {
		fmt.Println("--arr:", *arr)
		s := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%v", *arr), "["), "]")
		bd, _ = wechat.TextMsg(user, s, id)
		wechat.SendMsg(bd)
	}
}
