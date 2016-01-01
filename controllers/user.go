package controllers

import (
	"esap2/models"
	"time"
)

//login 登陆验证 GET/POST
type LoginController struct {
	baseRouter
}

//func (c *LoginController) Get() {
//	c.ParseMap("usr", "err")
//	if c.GetString("n") == "1" {
//		c.Data["err"] = "验证失败，用户或密码错误！"
//	}
//	c.TplNames = "login/login.html"
//}
func (c *LoginController) Get() {
	if _, ok := models.UserList[c.sesId]; ok {
		c.ServeRsts(true)
	} else {
		c.ServeRsts(false)
	}
}

func (c *LoginController) Post() {
	usr := models.User{c.GetString("username"), c.GetString("p"), 0, "", time.Now()}
	ok := usr.CheckPwd()
	if ok {
		models.UserList[c.sesId] = &usr
		c.Ctx.SetCookie("esapSID", c.sesId)
		c.Ctx.SetCookie("esapUsrDisp", usr.DispName)
		if c.IsAjax() {
			c.ServeRsts(true)
		} else {
			c.Redirect("/", 302)
		}
	} else {
		if c.IsAjax() {
			c.ServeRsts(false)
		} else {
			c.Redirect("/login?usr="+usr.Username+"&n=1", 302)
		}
	}
}
