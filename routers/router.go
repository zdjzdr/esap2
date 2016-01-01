package routers

import (
	"esap2/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/bdmp", &controllers.BdmpController{})
	beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/es/:m:string", &controllers.EsController{})
	beego.Router("/esm/:m:string", &controllers.EsmController{})
	beego.Router("/esd/:m:string", &controllers.EsdController{})
	beego.Router("/esv/:m:string", &controllers.EsvController{})
	beego.Router("/wx", &controllers.WxController{})
}
