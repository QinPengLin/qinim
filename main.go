package main

import (
	"qinim/common"
	_ "qinim/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"qinim/models"
	"qinim/controllers"
)

func main() {
	common.RedisConn = common.RegisterRedis()

	// 初始化日志
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"log/log.log","separate":["error","info","notice"]}`)
	beego.SetLogFuncCall(true)

	// 初始化数据库
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 注册默认数据库
	myws := beego.AppConfig.String("myws_mysql")
	orm.RegisterDataBase("default", "mysql", myws)

	//初始化
	models.Initializes()

	//错误页面自定义
	beego.ErrorController(&controllers.ErrorController{})

	beego.Run()
}

