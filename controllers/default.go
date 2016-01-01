package controllers

import (
	"time"
)

// 主界面
type IndexController struct {
	BaseAdminRouter
}

func (c *IndexController) Get() {
	c.TplNames = "index/index.html"
}

//bdmp 百度地图专用接口 GET
type BdmpController struct {
	baseRouter
}

func (c *BdmpController) Get() {
	c.ParseMap("lng", "lat")
	c.TplNames = "bdmp/map.html"
}

//upload 文件上传接口 GET/POST
type UploadController struct {
	BaseAdminRouter
}

func (c *UploadController) Get() {
	if c.GetString("n") == "1" {
		c.Data["n"] = "上传完毕，您可以继续上传"
	}
	c.TplNames = "upload/index.html"
}

func (c *UploadController) Post() {
	c.SaveToFile("uploadname", "uploads/"+time.Now().Format("20060102150405")+"")
	c.Redirect("/upload?n=1", 302)
}
