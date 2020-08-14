package controllers

import (
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (c *ErrorController) Error404() {
	c.TplName = "error/error.html"
}

func (c *ErrorController) Error501() {
	c.TplName = "error/501.html"
}

func (c *ErrorController) ErrorDb() {
	c.TplName = "error/dberror.html"
}