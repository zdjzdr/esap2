package controllers

import (
	"db/sqlsrv"
	"fmt"
)

//esMainData 明细数据API，返回json
type EsvController struct {
	BaseAdminRouter
}

func (c *EsvController) Get() {
	page, _ := c.GetInt("page")
	limit, _ := c.GetInt("limit")
	front := limit * (page - 1)
	sqlstmt0 := fmt.Sprintf("select 1 from %s", c.m)
	sqlstmt1 := fmt.Sprintf("select top %d * from %s where id not in (select top %d id from %s order by id) order by id", limit, c.m, front, c.m)
	c.ServeRsts(true, "", sqlsrv.FetchAll(sqlstmt1), sqlsrv.NumRows(sqlstmt0))
}
