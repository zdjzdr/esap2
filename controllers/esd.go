package controllers

//esdata 明细数据API，返回json
type EsdController struct {
	BaseAdminRouter
}

func (c *EsdController) Post() {
	c.exeCurd("c")
}

func (c *EsdController) Put() {
	c.exeCurd("u")
}

func (c *EsdController) Delete() {
	c.exeCurd("d")
}
