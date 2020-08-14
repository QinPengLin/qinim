package controllers

import (
	"qinim/models"
	"sync"
)

var (
	statusCode         sync.Map
)

func initStatusCode(){
	success:=models.SendMessageShelf{
		Code: "c_1",
		Msg:  "ok",
	}
	statusCode.Store("success",success)

	heart:=models.SendMessageShelf{
		Code: "c_1001",
		Msg:  "heart",
	}
	statusCode.Store("heart",heart)

	aunoempty:=models.SendMessageShelf{
		Code: "c_1002",
		Msg:  "新增uid不能为空",
	}
	statusCode.Store("aunoempty",aunoempty)

	auuse:=models.SendMessageShelf{
		Code: "c_1003",
		Msg:  "新增uid已经被使用",
	}
	statusCode.Store("auuse",auuse)

	auexist:=models.SendMessageShelf{
		Code: "c_1004",
		Msg:  "已经新增过了",
	}
	statusCode.Store("auexist",auexist)

	uidillegal:=models.SendMessageShelf{
		Code: "c_1005",
		Msg:  "uid不合法",
	}
	statusCode.Store("uidillegal",uidillegal)

	ausuccess:=models.SendMessageShelf{
		Code: "c_1006",
		Msg:  "uid新增成功",
	}
	statusCode.Store("ausuccess",ausuccess)

	noaccess:=models.SendMessageShelf{
		Code: "c_1007",
		Msg:  "无权限",
	}
	statusCode.Store("noaccess",noaccess)

	munoempty:=models.SendMessageShelf{
		Code: "c_1008",
		Msg:  "修改uid不能为空",
	}
	statusCode.Store("munoempty",munoempty)

	munom:=models.SendMessageShelf{
		Code: "c_1009",
		Msg:  "不需要修改uid",
	}
	statusCode.Store("munom",munom)

	muuse:=models.SendMessageShelf{
		Code: "c_1010",
		Msg:  "修改的uid被占用",
	}
	statusCode.Store("muuse",muuse)

	musuccess:=models.SendMessageShelf{
		Code: "c_1011",
		Msg:  "修改uid成功",
	}
	statusCode.Store("musuccess",musuccess)

	pcnotome:=models.SendMessageShelf{
		Code: "c_1012",
		Msg:  "不能给自己发送消息",
	}
	statusCode.Store("pcnotome",pcnotome)

	pcuserno:=models.SendMessageShelf{
		Code: "c_1013",
		Msg:  "玩家不存在或者不在线",
	}
	statusCode.Store("pcuserno",pcuserno)

	cgmax:=models.SendMessageShelf{
		Code: "c_1014",
		Msg:  "初始化群成员过多",
	}
	statusCode.Store("cgmax",cgmax)

	cgerr:=models.SendMessageShelf{
		Code: "c_1015",
		Msg:  "新增群失败",
	}
	statusCode.Store("cgerr",cgerr)

	cgsuccess:=models.SendMessageShelf{
		Code: "c_1016",
		Msg:  "新增群成功",
	}
	statusCode.Store("cgsuccess",cgsuccess)

	cgnoe:=models.SendMessageShelf{
		Code: "c_1017",
		Msg:  "群名称不能为空",
	}
	statusCode.Store("cgnoe",cgnoe)

	stgnogroup:=models.SendMessageShelf{
		Code: "c_1018",
		Msg:  "群不存在",
	}
	statusCode.Store("stgnogroup",stgnogroup)

	worldmsg:=models.SendMessageShelf{
		Code: "c_1019",
		Msg:  "世界消息",
	}
	statusCode.Store("worldmsg",worldmsg)
}

func getStatusCode(key string) models.SendMessageShelf {
	data,ok:=statusCode.Load(key)
	if ok {
		return data.(models.SendMessageShelf)
	}
	return models.SendMessageShelf{
		Code: "c_1999",
		Msg:  "未知状态",
	}
}