/**
 * 企业号API实例-应用分支实现.
 * @woylin, 2016-1-6
 */
package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

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
			bd, _ = wechat.TextMsg(user, s, id)
			//			fmt.Printf("%v\n", string(bd))
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
		//		fmt.Printf("%v\n", string(bd))
		wechat.SendMsg(bd)
	}
}

//7.订单进度
type Agent7 struct {
	WxAgent
}

func (w *Agent7) Gtext() {
	//订单进度
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在查询进度...")
	go ddjd(w.req.FromUserName, w.req.AgentID, w.req.Content)
}

//订单进度
type ddJd struct {
	SoNo  string
	SoRN  string
	Ddate time.Time
	Qty   float32
	Pqty  float32 //计划
	Cqty  float32 //下达
	Fqty  float32 //前
	Bqty  float32 //后
	Rqty  float32 //入库
	Dqty  float32 //发货
	Vqty  float32 //开票
}

func (c ddJd) String() string {
	s := fmt.Sprintf("订单号:%v-%v\n交期:%v\n订单总量 = %v\n", c.SoNo, c.SoRN, c.Ddate.Format("2006-1-2"), c.Qty)
	if c.Pqty > 0 {
		s += fmt.Sprintf("计划生产:%v\n", c.Pqty)
	}
	if c.Cqty > 0 {
		s += fmt.Sprintf("已下达量:%v\n", c.Pqty)
	}
	if c.Fqty > 0 {
		s += fmt.Sprintf("包标完成:%v\n", c.Pqty)
	}
	if c.Bqty > 0 {
		s += fmt.Sprintf("包装完成:%v\n", c.Pqty)
	}
	if c.Rqty > 0 {
		s += fmt.Sprintf("入库完成:%v\n", c.Pqty)
	}
	if c.Dqty > 0 {
		s += fmt.Sprintf("发货完成:%v\n", c.Pqty)
	}
	if c.Vqty > 0 {
		s += fmt.Sprintf("已开票数:%v\n", c.Pqty)
	}
	//	fmt.Println("--s:", s)
	return s
}

type ddJdWl struct {
	Ddate time.Time
	WlNo  string //物流
	KdNo  string //快递
}

func (c ddJdWl) String() string {
	s := fmt.Sprintf("相关发货日期:%v\n", c.Ddate.Format("2006-1-2"))
	if c.WlNo != "" {
		//		s += fmt.Sprintf("物流单号:%v http://m.ickd.cn/result.html?no=%v&com=auto \n", c.WlNo, c.WlNo)
		s += fmt.Sprintf("物流单号:%v\n点击查看物流信息：http://m.ickd.cn/result.html?no=%v", c.WlNo, c.WlNo)
		//		http://m.ickd.cn/result.html#no=04843844&com=auto
	}
	if c.KdNo != "" {
		s += fmt.Sprintf("快递单号:%v\n", c.KdNo)
	}
	return s
}

//订单进度
func ddjd(user string, id int, mDesc string) {
	dd := &ddJd{}
	sql := fmt.Sprintf("SELECT 单号,项,交期,数量,计划=ISNULL(计划,0),下达=ISNULL(下达,0),前道=ISNULL(前道,0),完工=ISNULL(完工,0),入库=ISNULL(入库,0),发货=ISNULL(发货,0),开票=ISNULL(开票,0) FROM [esap].[dbo].[vSDSOplus] where charindex('%s',soNoRn)>0", mDesc)
	//	queryAndSendArr(user, id, sql, dd)
	time.Sleep(time.Second)
	arr := sqlsrv.FetchAllRowsPtr(sql, dd)
	bd, _ := wechat.TextMsg(user, "未找到项目...", id)
	if (*arr) != nil {
		fmt.Println("--arr:", *arr)
		s := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%v", *arr), "["), "]")
		bd, _ = wechat.TextMsg(user, s, id)
		//		fmt.Printf("%v\n", string(bd))
		wechat.SendMsg(bd)
	}
	fmt.Println("len-arr:", len(*arr))
	if len(*arr) == 1 {
		if v, ok := (*arr)[0].(ddJd); ok {
			sql := fmt.Sprintf("SELECT ddate,物流号,快递号 from sdd_d where soNo='%s' and soRN='%s'", v.SoNo, v.SoRN)
			fmt.Println("sql:", sql)
			queryAndSendArr(user, id, sql, &ddJdWl{})
		}
	}
}

//10.考勤签到
type Agent10 struct {
	WxAgent
}

func (w *Agent10) Gevent() {
	switch w.req.Event {
	case "LOCATION":
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, fmt.Sprintf("您的地址信息已采集：经度：%v，纬度：%v", w.req.Location_X, w.req.Location_Y))
	}
}

//6.照片采集
type Agent6 struct {
	WxAgent
}

func (w *Agent6) Gtext() {
	//匹配姓名工号并存入ESAP
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在处理...")
	go zpcl(w.req)
}

func (w *Agent6) Gimage() {
	//接收图片，提示录入
	empMap[w.req.FromUserName] = &employee{"", "", w.req.PicUrl, nil, 0}
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "请输入姓名，工号\n例如：“张三，120”")
}

func (w *Agent6) Gevent() {
	switch w.req.Event {
	case "enter_agent":
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "请拍摄或选择相册照片发送后，填写姓名，工号。")
	}
}

const (
	FTPPATH   = `E:\esDisk\hr\emp\` //照片存放路径，ES需启用网盘功能
	PICPREFIX = "/P00"              //照片前缀，默认是"P00"
	PICSUFFIX = ".jpg"              //照片后缀，默认".jpg"
)

//定义员工
type employee struct {
	Name     string
	Eid      string
	Photo    string
	Rcid     interface{}
	errCount int
}

//员工照片上传，可为网盘
func (e *employee) download(url string) {
	imgResp, _ := http.Get(url) //从微信服务器获取照片
	defer imgResp.Body.Close()
	m, _, _ := image.Decode(imgResp.Body)         //解析照片信息
	fn := FTPPATH + PICPREFIX + e.Eid + PICSUFFIX //设置照片存放路径及文件名
	fh, _ := os.Create(fn)                        //重建照片
	defer fh.Close()
	jpeg.Encode(fh, m, nil) //以JPG格式保存照片
}

//定义员工表、存放照片链接信息
var empMap = make(map[string]*employee)

//订单进度
func zpcl(w *wechat.WxReq) {
	//检验是否已上传图片，上传过则提示输入“姓名，工号”
	var bd string
	if v, ok := empMap[w.FromUserName]; ok {
		//解析姓名工号
		empInfo := strings.Split(w.Content, "，")
		if len(empInfo) == 2 {
			v.Name, v.Eid = empInfo[0], empInfo[1]
		} else {
			//重复出错4次后重置，取消流程
			v.errCount++
			if v.errCount > 3 {
				bd = "出错已超限，上传流程已取消"
				delete(empMap, w.FromUserName)
			} else {
				bd = fmt.Sprintf("你填入的信息格式不正确，请重新输入。(%d)", v.errCount)
			}
			bd1, _ := wechat.TextMsg(w.FromUserName, bd, w.AgentID)
			wechat.SendMsg(bd1)
			return
		}
		//通过员工信息表查找rcid,未找到则提示用户。
		v.Rcid = sqlsrv.Fetch("select Excelserverrcid from 员工信息表 where 姓名=? and 工号=?", v.Name, v.Eid)
		if (v.Rcid) == nil {
			bd = "未找到该用户，请重新输入或重新上传图片。"
		} else {
			//设置照片名为：前缀 + 工号
			picName := PICPREFIX + v.Eid
			//删除原有数据库照片路径记录，r,c是照片的Excel行列号，示例是R5C7,也就是[G5]单元格
			sqlsrv.Exec("delete from es_casepic where rcid=? and r=5 and c=7", v.Rcid)
			//向数据库插入新照片路径记录,6496，1，5，7，ed\esys，hr\emp要改成自己的
			err := sqlsrv.Exec("insert es_casepic(rcid,picNo,fileType,rtfid,sh,r,c,saveinto,nfsfolderid,nfsfolder,relafolder,phyfileName) values(?,?,?,?,?,?,?,?,?,?,?,?)",
				v.Rcid, picName, ".jpg", 6496, 1, 5, 7, 1, 1, `ed\esys`, `hr\emp`, picName+".jpg")
			if err != nil {
				bd = "图片上传失败"
			} else {
				//			上传图片到网盘，销毁员工数组中的信息，回复处理成功信息
				v.download(v.Photo)
				delete(empMap, w.FromUserName)
				bd = "员工照片已成功处理"
			}
		}
		bd1, _ := wechat.TextMsg(w.FromUserName, bd, w.AgentID)
		wechat.SendMsg(bd1)
	}
}

//13.工作记录
type Agent13 struct {
	WxAgent
}

func (w *Agent13) Gtext() {
	err := sqlsrv.Exec("insert into wx_gzjl(cdate,usr,ctx) values(getdate(),?,?)", w.req.FromUserName, w.req.Content)
	if err != nil {
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "保存失败")
	}
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "记录已保存")
}

//func (w *Agent13) Gimage() {
//	//接收图片，提示录入
//	empMap[w.req.FromUserName] = &employee{"", "", w.req.PicUrl, nil, 0}
//	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "请输入姓名，工号\n例如：“张三，120”")
//}
type gzJl struct {
	Cdate time.Time
	Usr   string
	Ctx   string
}

func (c gzJl) String() string {
	return fmt.Sprintf("%v %v\n", c.Cdate.Format("2006-1-2"), c.Ctx)
}

func (w *Agent13) Gevent() {
	switch w.req.Event {
	case "enter_agent":
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "切换到输入模式填写工作记录，发送后将上传到ESAP工作记录中(不占用手机存储)。")
	case "click":
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在下载本周记录...")
		sql := fmt.Sprintf("SELECT cDate,usr,ctx from wx_gzjl where usr ='%s' order by cdate ", w.req.FromUserName)
		fmt.Println("sql:", sql)
		go queryAndSendArr(w.req.FromUserName, w.req.AgentID, sql, &gzJl{})
	}
}

//15.资产台账
type Agent15 struct {
	WxAgent
}

func (w *Agent15) Gtext() {
	w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "正在搜索...")
	sql := fmt.Sprintf("SELECT 资产编码,类别,资产名称,型号,变动方式,使用日期,数量,单位,制造商,原值原币 FROM 固定资产台账_主表 where charindex('%s',资产编码)>0", w.req.Content)
	go queryAndSendArr(w.req.FromUserName, w.req.AgentID, sql, &zcTz{})
}

type zcTz struct {
	No     string
	Type   string
	Name   string
	Spec   string
	Method string
	Cdate  time.Time
	Qty    float32
	Unit   string
	Vendor string
	Price  float32
}

func (c zcTz) String() string {
	return fmt.Sprintf("资产编号：%v\n资产类别：%v\n资产名称：%v\n型号：%v\n变动方式：%v\n使用日期：%v\n数量(单位)：%v%v\n原值：￥%v\n制造商：%v\n",
		c.No, c.Type, c.Name, c.Spec, c.Method, c.Cdate.Format("2006-1-2"), c.Qty, c.Unit, c.Price, c.Vendor)
}

func (w *Agent15) Gevent() {
	switch w.req.Event {
	case "enter_agent":
		w.resp, _ = wechat.RespText(w.req.ToUserName, w.req.FromUserName, "温馨提示：填入设备编码即可查询资产信息。")
	}
}
