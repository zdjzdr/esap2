/**
 * 企业号API实例-工具函数
 * @woylin, since 2016-1-6
 */
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/woylyn/esap2/db/sqlsrv"
	"github.com/woylyn/esap2/wechat"
)

func queryMaxDate(user string, id int, table string) {
	date := sqlsrv.Fetch(fmt.Sprintf("select max(cdate) from %s", table))
	if v, ok := (*date).(time.Time); ok {
		bd, _ := wechat.TextMsg(user, fmt.Sprintf("最新记录日期: %s", v.Format("2006-01-02")), id)
		wechat.SendMsg(bd)
	}
}

//通用方法，逐条回复sql查询到的内容
func queryAndSend(user string, id int, sql string, struc interface{}) {
	defer rc("from queryAndSend")
	time.Sleep(time.Second)
	arr := sqlsrv.FetchAllRowsPtr(sql, struc)
	bd, _ := wechat.TextMsg(user, "未找到...", id)
	if (*arr) != nil {
		fmt.Println("--arr:", *arr)
		for k, v := range *arr {
			s := fmt.Sprintf("%v：%v", k+1, v)
			if len(*arr) == 1 {
				s = fmt.Sprintf("%v", v)
			}
			bd, _ = wechat.TextMsg(user, s, id)
			wechat.SendMsg(bd)
		}
	}
}

//通用方法，合并回复sql查询到的内容（更常用）
func queryAndSendArr(user string, id int, sql string, struc interface{}) {
	defer rc("from queryAndSendArr")
	time.Sleep(time.Second)
	arr := sqlsrv.FetchAllRowsPtr(sql, struc)
	bd, _ := wechat.TextMsg(user, "未找到...", id)
	if (*arr) != nil {
		fmt.Println("--arr:", *arr)
		s := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprintf("%v", *arr), "["), "]")
		bd, _ = wechat.TextMsg(user, s, id)
		wechat.SendMsg(bd)
	}
}

//异常恢复
func rc(s ...string) {
	err := recover()
	if err != nil {
		fmt.Println("err:", err, s)
	}
}
