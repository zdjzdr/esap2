package controllers

import (
	"github.com/astaxie/beego"
)

type NestPreparer interface {
	NestPrepare()
}

// baseRouter implemented global settings for all other routers.
type baseRouter struct {
	beego.Controller
	rst   map[string]interface{}
	sesId string
}

// Prepare implemented Prepare method for baseRouter.
func (c *baseRouter) Prepare() {
	c.sesId = c.Ctx.GetCookie("esap2SID")
	c.rst = make(map[string]interface{})
	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}

// Parse []string to this.data by woylin 2015-11
func (c *baseRouter) ParseMap(s ...string) {
	for _, v := range s {
		c.Data[v] = c.GetString(v)
	}
}

// ServeJson() via this.rst by woylin 2015-11
func (c *baseRouter) ServeRst() {
	c.Data["json"] = c.rst
	c.ServeJson()
}

// ServeJson() via sets by woylin 2015-12-19
func (c *baseRouter) ServeRsts(v ...interface{}) {
	c.rst["success"] = v[0]
	if len(v) > 1 {
		c.rst["msg"] = v[1]
	}
	if len(v) > 2 {
		c.rst["data"] = v[2]
	}
	if len(v) > 3 {
		c.rst["total"] = v[3]
	}
	c.ServeRst()
	c.StopRun()
}
