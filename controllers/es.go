package controllers

import (
	"db/sqlsrv"
	"esap2/models"
	"fmt"
)

//es 数据库查询接口，返回json GET/POST
type EsController struct {
	BaseAdminRouter
}

func (c *EsController) Get() {
	uId := models.UserList[c.sesId].Id
	switch c.m {
	case "menu":
		sqlstmt1 := "select * from es_v_UserMenu_0 where userId = ? order by ordPath,rtNo"
		rstArr := sqlsrv.FetchMenuTree(sqlstmt1, uId)
		c.Data["json"] = *rstArr
		fmt.Println("menu", *rstArr)
		c.ServeJson()
		return
	}
	c.ServeRst()
}

func (c *EsController) Post() {
	switch c.m {
	case "chgPwd":
		uId := models.UserList[c.sesId].Id
		if !sqlsrv.CheckBool("SELECT * FROM es_user where userid=? and UserPwd=?", uId, c.reqMap["p1"]) {
			c.ServeRsts(false, "密码不正确。")
			break
		}
		if c.reqMap["p2"] != c.reqMap["p3"] {
			c.ServeRsts(false, "密码不一致。")
			break
		} else {
			err := sqlsrv.Exec("update es_user set UserPwd=? where userid=?", c.reqMap["p2"], uId)
			if err != nil {
				c.ServeRsts(false, "密码更新失败。")
				break
			}
		}
		c.rst["success"] = true
	}
	c.ServeRst()
}
