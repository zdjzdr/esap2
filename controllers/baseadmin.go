package controllers

import (
	"db/sqlsrv"
	"encoding/json"
	"esap2/models"
	"fmt"
	"io/ioutil"
)

type Sqlstmt struct {
	Sql  string
	Cond []string
}

type BaseAdminRouter struct {
	baseRouter
	Islogin bool
	ob      models.Models
	curd    map[string]Sqlstmt
	reqMap  map[string]interface{}
	m       string
}

func (c *BaseAdminRouter) readFile(filename string) error {
	bytes, err := ioutil.ReadFile("esap/curd/" + filename + ".json")
	if err != nil {
		fmt.Println("ReadFile:", err)
		return err
	}

	if err := json.Unmarshal(bytes, &(c.curd)); err != nil {
		fmt.Println("Unmarshal:", err)
		return err
	}

	return nil
}

func (c *BaseAdminRouter) NestPrepare() {
	c.sesId = c.Ctx.GetCookie("esapSID")
	c.m = c.Ctx.Input.Param(":m")
	c.curd = make(map[string]Sqlstmt)
	c.reqMap = make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.CopyBody(), &(c.reqMap))
	_ = c.readFile(c.m)
	if c.sesId != "" && models.CheckLogined(c.sesId) {
		c.Islogin = true
	} else {
		c.Islogin = false
	}
	//非ajax请求时跳转到登陆界面,否则返回json{"success": false}
	if !c.Islogin {
		if c.IsAjax() {
			c.ServeRsts(false, "notlogin")
		}
	}
}

func (c *BaseAdminRouter) Get() {
	id := c.GetString("id")
	switch c.m {
	default:
		page, _ := c.GetInt("page")
		limit, _ := c.GetInt("limit")
		front := limit * (page - 1)
		sqlstmt0 := fmt.Sprintf("select 1 from %s", c.m)
		sqlstmt1 := fmt.Sprintf("select top %d * from %s where id not in (select top %d id from %s order by excelserverrcid desc ) order by id desc ", limit, c.m, front, c.m)
		if id != "" {
			sqlstmt0 = fmt.Sprintf(sqlstmt0+" where pid='%s'", id)
			sqlstmt1 = fmt.Sprintf("select top %d * from %s where pid='%s' and id not in (select top %d id from %s where pid='%s' order by id) order by id", limit, c.m, id, front, c.m, id)
		}
		c.ServeRsts(true, "", sqlsrv.FetchAll(sqlstmt1), sqlsrv.NumRows(sqlstmt0))
	}
	c.ServeRst()
}

//parse sql from json to CURD
func (c *BaseAdminRouter) exeCurd(m string) {
	sqlstr := c.curd[m].Sql
	sqlcond := c.curd[m].Cond
	sqlcond2 := make([]interface{}, 0)
	for _, v := range sqlcond {
		if val, ok := c.reqMap[v]; ok {
			sqlcond2 = append(sqlcond2, val)
		} else {
			val2, err := c.GetInt(v)
			if err != nil {
				panic(err)
			}
			sqlcond2 = append(sqlcond2, val2)
		}

	}
	err := sqlsrv.Exec(sqlstr, sqlcond2...)
	if err != nil {
		c.ServeRsts(false, err.Error())
	} else {
		c.ServeRsts(true, "1条记录被更新")
	}
}
