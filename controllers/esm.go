package controllers

//esMainData 明细数据API，返回json
type EsmController struct {
	BaseAdminRouter
}

func (c *EsmController) Post() {
	c.exeCurd("c")
}

func (c *EsmController) Put() {
	c.exeCurd("u")
}

func (c *EsmController) Delete() {
	c.exeCurd("d")
}
