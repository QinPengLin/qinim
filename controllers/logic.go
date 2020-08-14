package controllers

import (
	"qinim/extra"
	"qinim/models"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

//验证是否设置uid
func VerifyUser(key,verifyUser string,data *models.FromMessage) bool  {
	//验证玩家是否绑定uid
	if verifyUser==models.VerifyUser {
		//如果当前的nowUid为空说明还未设置uid必须要先设置uid才能发送消息
		if GetKeyToUid(key)=="" {
			ToMe(key,getStatusCode("noaccess"))
			return true
		}else {
			return false
		}
	}
	return false
}

//修改uid
func  ModifyUid(key,verifyUser string,data *models.FromMessage) bool {
	if data.Uid=="" {
		//uid不能为空
		logs.Info("修改的uid【" + data.Uid + "】不能为空")
		ToMe(key,getStatusCode("munoempty"))
		return true
	}
	//先查看之前是否有关系
	value,ok := uidToKey.Load(data.Uid)
	if ok {//存在
		if value.(string)==key {
			logs.Info("uid【"+data.Uid+"】无需修改")
			ToMe(key,getStatusCode("munom"))
			return true
		}else {
			logs.Info("将uid【" + data.Uid + "】已经存在")
			ToMe(key,getStatusCode("muuse"))
			return true
		}
	}

	/**判断用户是否合法*/
	if verifyUser==models.VerifyUser {
		if !GetUserToRedis(data.Uid) {
			logs.Info("修改的uid【" + data.Uid + "】不合法")
			ToMe(key,getStatusCode("uidillegal"))
			return true
		}
	}
	/**判断用户是否合法*/

	//如果不存在，就先查找key是否存在于uidToKey如果存在就删除不存在就新增（这就是修改uid功能）
	uidToKey.Range(func(uidToKeyK, uidToKeyV interface{}) bool {
		if uidToKeyV==key {//存在就删除
			uidToKey.Delete(uidToKeyK)
			logs.Info("删除了之前的uid【"+uidToKeyK.(string)+"】key为【"+key+"】的用户列表")
			return true
		}
		return true
	})
	//新增或者删除
	uidToKey.Store(data.Uid,key)

	//记录当前的uid
	publicKeyToUid.Store(key,data.Uid)

	logs.Info("将uid【"+data.Uid+"】绑定到了key为【"+key+"】")
	ToMe(key,getStatusCode("musuccess"))
	return true
}

//新增uid
func  AddUid(key,verifyUser string,data *models.FromMessage) bool {
	if data.Uid=="" {
		//uid不能为空
		logs.Info("新增的uid【" + data.Uid + "】不能为空")
		ToMe(key,getStatusCode("aunoempty"))
		return true
	}
	//查看当前是否已经设置过uid
	if GetKeyToUid(key)!="" {
		logs.Info("将uid【" + data.Uid + "】已经新增过")
		ToMe(key,getStatusCode("auexist"))
		return true
	}
	//先查看之前是否有关系
	_,ok := uidToKey.Load(data.Uid)
	if ok {//存在
		logs.Info("将uid【" + data.Uid + "】已经存在")
		ToMe(key,getStatusCode("auuse"))
		return true
	}

	/**判断用户是否合法*/
	if verifyUser==models.VerifyUser {
		if !GetUserToRedis(data.Uid) {
			logs.Info("新增的uid【" + data.Uid + "】不合法")
			ToMe(key,getStatusCode("uidillegal"))
			return true
		}
	}
	/**判断用户是否合法*/

	//新增
	uidToKey.Store(data.Uid,key)
	//记录当前的uid
	publicKeyToUid.Store(key,data.Uid)

	logs.Info("将uid【"+data.Uid+"】绑定到了key为【"+key+"】")
	ToMe(key,getStatusCode("ausuccess"))
	return true
}

//心跳
func  Heart(key,verifyUser string,data *models.FromMessage) bool {
	ToMe(key,getStatusCode("heart"))
	return true
}

//私聊
func PrivateChat(key,verifyUser string,data *models.FromMessage) bool {
	//单对单
	//不能给自己发送消息（给自己发消息会产生websocket并发写入）
	publicUid:=GetKeyToUid(key)
	if data.ToUid==publicUid {
		logs.Info("ToUid【" + data.ToUid + "】不能给自己发消息")
		ToMe(key,getStatusCode("pcnotome"))
		return true
	}
	//获取toUid对应的key
	_,ok := uidToKey.Load(data.ToUid)
	if !ok {//不存在就保存记录等到玩家上线发送并且跳过该次循环（后面做）
		logs.Info("内存中无该玩家uid【"+data.ToUid+"】OneToOne")
		ToMe(key,getStatusCode("pcuserno"))
		return true
	}

	/**组装发送的消息*/
	senCode:=getStatusCode("success")
	senMessag:=models.SendMessageShelf{
		Code: senCode.Code,
		Msg:  senCode.Msg,
		Data: &models.SendMessage{
			FromUid:   publicUid,
			ToUid:     data.ToUid,
			Time:      strconv.FormatInt(time.Now().Unix(),10),
			Msg:       data.Content,
			File:      data.File,
		},
	}
	/**组装发送的消息*/

	//发送给对应的玩家
	go OneToOne(data.ToUid,senMessag)

	//给自己发消息
	go OneToOne(publicUid,senMessag)
	return true
}

//创建群
func  CreateGroup(key,verifyUser string,data *models.FromMessage) bool {
	if data.GroupName=="" {//群名不能为空
		ToMe(key,getStatusCode("cgnoe"))
		return true
	}
	successUid := ""
	failUid := ""
	var createUids []models.GroupMembers
	var groupMember models.GroupMembers

	publicUid:=GetKeyToUid(key)

	ds:=beego.AppConfig.String("group_uid_partition_code")
	if data.GroupInitMembers!="" {//需要对传过来的uid做验证
		uidStrArr,l:=extra.StrPartition(data.GroupInitMembers,ds)
		//检测uid个数是否合法
		maxGroupMembers,_:=beego.AppConfig.Int("max_group_members")
		if l>maxGroupMembers {//超过群最大容纳人数
			ToMe(key,getStatusCode("cgmax"))
			return true
		}

		for _,uidV := range uidStrArr {
			if GetUserToRedis(uidV) {
				if uidV!=publicUid {
					successUid = successUid + ds + uidV
					groupMember.Uid = uidV
					groupMember.Level = "2"
					createUids = append(createUids,groupMember)
				}
			}else {
				failUid = failUid+ds+uidV
			}
		}
	}

	groupMember.Uid = publicUid
	groupMember.Level = "1"
	createUids = append(createUids,groupMember)

	//创建群id
	groupId:=extra.CreateOnlKey()

	_,addErr:=AddGroupToRdis(groupId,createUids,data.GroupName)
	if addErr!=nil {
		ToMe(key,getStatusCode("cgerr"))
		return  true
	}

	//成功，组装返回信息
	successUid = successUid + ds + publicUid
	cgsuccess:=getStatusCode("cgsuccess")
	senMessagGadd:=models.SendToMeMessageShelf{
		Code: cgsuccess.Code,
		Msg:  cgsuccess.Msg,
		Data: &models.SendToMeMessage{
			SuccessUid: successUid,
			FailUid:    failUid,
			GroupId:    groupId,
		},
	}
	ToMeData(key,senMessagGadd)
	return true

}

//向某个群发送消息
func SendToGroup(key,verifyUser string,data *models.FromMessage) bool {
	//查找群id是否存在
	groupData:=GetGroupDataToRedis(data.ToGroupId)
	if groupData=="" {//群id不能为空或者群不存在
		ToMe(key,getStatusCode("stgnogroup"))
		return true
	}

	publicUid:=GetKeyToUid(key)

	//将groupData解析成数组
	var groupDataArr []models.GroupMembers

	json.Unmarshal([]byte(groupData), &groupDataArr)
	inGroup:=false
	for _,gV := range groupDataArr {
		if publicUid==gV.Uid {
			inGroup = true
			break
		}
	}
	//需要先判断发消息人是否在该群中
	if inGroup {//存在才能群发
		//组装消息结构体
		strNowTime:=strconv.FormatInt(time.Now().Unix(),10)
		code:=getStatusCode("success")
		senMessagToMe:=models.SendMessageShelf{
			Code: code.Code,
			Msg:  code.Msg,
			Data: &models.SendMessage{
				GroupName: "",
				ToGroupId: data.ToGroupId,
				FromUid:   publicUid,
				Time:      strNowTime,
				Msg:       data.Content,
				File:      data.File,
			},
		}

		//将群发消息写入群发记录表
		insertData := &models.GroupChatRecord{
			GroupId:       data.ToGroupId,
			FormUid:       publicUid,
			Content:       data.Content,
			CreationTime:  strNowTime,
		}
		models.InsertGroupChatRecord(insertData)
		for _,groupDataArrV := range groupDataArr {
			go OneToOne(groupDataArrV.Uid,senMessagToMe)
		}
		return true
	}
	//不存在就提示无权限
	ToMe(key,getStatusCode("noaccess"))
	return true
}

//世界消息，将广播给所有在线用户
func SendToWorld(key,verifyUser string,data *models.FromMessage) bool {
	//将信息放入通道,世界消息(只有发送给所有人才会把消息放入通道)
	/**将fromData信息放到SendMessage*/
	publicUid:=GetKeyToUid(key)
	worldmsg:=getStatusCode("worldmsg")
	senMessag:=models.SendMessageShelf{
		Code: worldmsg.Code,
		Msg:  worldmsg.Msg,
		Data: &models.SendMessage{
			FromUid:   publicUid,
			Time:      strconv.FormatInt(time.Now().Unix(),10),
			Msg:       data.Content,
			File:      data.File,
		},
	}
	/**将fromData信息放到SendMessage*/

	broadcast <- &senMessag
	return true
}