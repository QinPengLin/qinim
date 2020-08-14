package controllers

import (
	"qinim/common"
	"github.com/gorilla/websocket"
	"qinim/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"qinim/extra"
	"sync"
)

var (
	clientsAmount = 0                             //保存链接的数量
	uidToKey sync.Map                             //保存uid和key的对应关系uid=>key
	publicKeyToUid sync.Map                       //保存设置了uid的key关系是key=>uid
	clients sync.Map                              //保存链接
	broadcast = make(chan *models.SendMessageShelf)    //世界消息通道

	router sync.Map                               //保存Read循环中的处理方法
)

type dealFunc func(key,verifyUser string,data *models.FromMessage) bool

type MethodStruct struct {
	BeforeMethod          dealFunc      //之前需要之前的方法
	Method                dealFunc      //实际方法
}

func init() {
	/**注册逻辑处理方法*/
	router.Store(models.AddUid,MethodStruct{//新增uid方法,不需要验证uid是否存在
		Method:AddUid,
		})
	router.Store(models.HEART,MethodStruct{//心跳方法,不需要验证uid是否存在
		Method:Heart,
	})
	router.Store(models.ModifyUid,MethodStruct{//修改uid方法
		BeforeMethod:VerifyUser,//验证是否存在账号
		Method:ModifyUid,
	})
	router.Store(models.OneToOne,MethodStruct{//私聊
		BeforeMethod:VerifyUser,//验证是否存在账号
		Method:PrivateChat,
	})
	router.Store(models.CreateGroup,MethodStruct{//创建群
		BeforeMethod:VerifyUser,//验证是否存在账号
		Method:CreateGroup,
	})
	router.Store(models.ToGroup,MethodStruct{//群发消息
		BeforeMethod:VerifyUser,//验证是否存在账号
		Method:SendToGroup,
	})
	router.Store(models.Radio,MethodStruct{//发送世界消息广播
		BeforeMethod:VerifyUser,//验证是否存在账号
		Method:SendToWorld,
	})
	/**注册逻辑处理方法*/
	//初始化状态码
	initStatusCode()
	//发送世界消息监听协程
	go SendWorld()
}

func  Read(wsLink *websocket.Conn,key string)  {
	verifyUser:=beego.AppConfig.String("verify_user")
	logs.Info("key为【"+key+"】连接成功")
	for {
		logs.Info("【"+key+"】进入获取信息模式")
		_, message, err := wsLink.ReadMessage()
		if err != nil {
			logs.Info("key【"+key+"】离开了 error:",err)
			//关闭连接，删除clients
			wsLink.Close()
			DeleteClients(key)
			break
		}

		messageStr,_ := extra.Base64DecodeString(string(message))
		logs.Info("读取了链接【"+key+"】的发来的信息:",messageStr)
		fromData,_:=extra.JsonToFromMessageStruct(messageStr)

		methodStruct,ok:=router.Load(fromData.Method)
		if ok {
			methodOb:=methodStruct.(MethodStruct)
			if methodOb.BeforeMethod!=nil {//存在之前函数
				BeforeExeResult:=methodOb.BeforeMethod(key,verifyUser,fromData)
				if BeforeExeResult {
					continue
				}
			}
			if methodOb.Method!=nil {//真实执行方法
				exeResult:=methodOb.Method(key,verifyUser,fromData)
				if exeResult {
					continue
				}
			}
		}
	}
}

//获取某个key是否设置了uid
func GetKeyToUid(key string) string {
	data,ok:=publicKeyToUid.Load(key)
	if ok {
		return data.(string)
	}
	return ""
}

//删除链接map
func DeleteClients(key string)  {
	//删除连接map
	clients.Delete(key)
	//删除key=>uid的map
	uid:=GetKeyToUid(key)
	if uid!="" {
		//删除key=>uid的map
		publicKeyToUid.Delete(key)
		//删除uid=>key的map
		uidToKey.Delete(uid)
	}
	clientsAmount = clientsAmount-1
}

//获取指定用户是否在redis中
func GetUserToRedis(user_id string) bool {
	d := common.RedisConn.Get()
	defer d.Close()
	//查看哈希列表中是否存在
	ok, err := d.Do(models.REDIS_HEXISTS,models.USER_LIST,user_id)
	if err != nil {
		logs.Info("GetUserToRedis查询redis错误,Err:%v", err)
		return false
	}

	if ok.(int64)==1 {
		logs.Info("GetUserToRedis查询redis账号存在")
		return true
	}

	return false
}
