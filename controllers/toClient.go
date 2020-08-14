package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"qinim/extra"
	"qinim/models"
)

//服务器返回给自己的消息
func ToMe(key string ,msg models.SendMessageShelf)  {
	msgJson,_:=extra.StructToFromMessageJson(&msg)
	sendMsg:=extra.Base64EncodeString(msgJson)
	cl,_:=clients.Load(key)
	err := cl.(*websocket.Conn).WriteJSON(sendMsg)
	if err != nil {
		logs.Info("key为【"+key+"】断掉了ToMe error:",err)
		cl.(*websocket.Conn).Close()
		DeleteClients(key)
	}
}

//将消息发送给某一个人
func OneToOne(toUid string,msg models.SendMessageShelf) {
	//获取toUid对应的key
	key,ok := uidToKey.Load(toUid)
	if !ok {//不存在就保存记录等到玩家上线发送并且跳过该次循环（后面做）
		logs.Info("内存中无该玩家uid【"+toUid+"】OneToOne")
		return
	}

	msgJson,_:=extra.StructToFromMessageJson(&msg)
	sendMsg:=extra.Base64EncodeString(msgJson)
	cl,_:=clients.Load(key)
	err := cl.(*websocket.Conn).WriteJSON(sendMsg)
	if err != nil {
		logs.Info("key为【"+key.(string)+"】断掉了OneToOneerror:",err)
		cl.(*websocket.Conn).Close()
		DeleteClients(key.(string))
	}
	logs.Info("信息需要发送给【"+key.(string)+"】OneToOne")
}

//服务器返回带有数据的消息给自己
func ToMeData(key string,msg models.SendToMeMessageShelf)  {
	msgJson,_:=extra.StructToFromMeMessageJson(&msg)
	sendMsg:=extra.Base64EncodeString(msgJson)
	cl,_:=clients.Load(key)
	err := cl.(*websocket.Conn).WriteJSON(sendMsg)
	if err != nil {
		logs.Info("key为【"+key+"】断掉了ToMe error:",err)
		cl.(*websocket.Conn).Close()
		DeleteClients(key)
	}
}

//广播给所有玩家
func SendWorld() {
	for {
		//如果通道中没有消息将会阻塞
		msg := <-broadcast

		msgJson,_:=extra.StructToFromMessageJson(msg)
		sendMsg:=extra.Base64EncodeString(msgJson)

		clients.Range(func(k, client interface{}) bool {
			if GetKeyToUid(k.(string))!="" {//只广播给设置了uid的连接
				err := client.(*websocket.Conn).WriteJSON(sendMsg)
				if err != nil {
					logs.Info("key为【"+k.(string)+"】断掉了 error:",err)
					client.(*websocket.Conn).Close()
					DeleteClients(k.(string))
				}else {
					logs.Info("信息需要发送给【" + k.(string) + "】发送者uid【" + msg.Data.FromUid + "】")
				}
			}
			return true
		})
	}
}
