package controllers

import (
	"github.com/gorilla/websocket"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"qinim/extra"
	"net/http"
)

type WsController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *WsController) Get() {

	ResponseWriter:=c.Ctx.ResponseWriter
	Request:=c.Ctx.Request
	ws, err := upgrader.Upgrade(ResponseWriter, Request, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(ResponseWriter, "不是正确的websocket握手", 400)
		logs.Error("不是正确的websocket握手")

		c.StopRun()
	} else if err != nil {
		http.Error(ResponseWriter, "不是正确的websocket握手", 400)
		logs.Error("不是正确的websocket握手:", err)

		c.StopRun()
	}

	keyStr:=extra.CreateOnlKey()
	clients.Store(keyStr,ws)
	clientsAmount = clientsAmount+1

	logs.Info("创建了链接:",keyStr)

	go Read(ws,keyStr)

	c.StopRun()
}