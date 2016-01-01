package models

import (
	"db/sqlsrv"
	"fmt"
	"time"
)

const SES_LIFE_TIME = -20 //session连接保持20分钟
var UserList map[string]*User

type User struct {
	Username    string `json:"username"`
	Pwd         string `json:"pwd"`
	Id          int64  `json:"id"`
	DispName    string `json:"dispname"`
	LastActTime time.Time
}

func init() {
	UserList = make(map[string]*User)
}

func CheckLogined(sid string) bool {
	userlistGC()
	for k, v := range UserList {
		if k == sid {
			v.LastActTime = time.Now()
			return true
		}
	}
	return false
}

func userlistGC() {
	for k, v := range UserList {
		if v.LastActTime.Before(time.Now().Add(time.Duration(SES_LIFE_TIME) * time.Minute)) {
			delete(UserList, k)
		}
	}
}

func (m *User) CheckPwd() bool {
	ok := sqlsrv.CheckBool("SELECT * FROM es_user where UserLogin=? and UserPwd=?", m.Username, m.Pwd)
	if !ok && m.Pwd == "" {
		ok = sqlsrv.CheckBool("SELECT * FROM es_user where UserLogin=? and UserPwd=''", m.Username)
	}
	if ok {
		a1 := sqlsrv.Fetch("select UserId from Es_User where UserLogin=?", m.Username)
		a2 := sqlsrv.Fetch("select DispName from Es_User where UserLogin=?", m.Username)
		value, ok1 := (*a1).(int64)
		if ok1 {
			m.Id = value
			m.DispName = fmt.Sprintf("%s", *a2)
			return true
		}
	}
	return false
}
