package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"qinim/common"
	"qinim/extra"
	"qinim/models"
	"strconv"
	"strings"
	"time"
)

type UserApiController struct {
	beego.Controller
}

//将用户设置到聊天系统
func (c *UserApiController) SetUser() {
	user_id := c.GetString("user_id")
	item := c.GetString("item")
	if user_id=="" || item=="" {
		logs.Info("参数错误,请检查")
		c.Data["json"] = models.ApiResponse{Code: models.ErrParam, Message: "参数错误,请检查", Data: ""}
		c.ServeJSON()
		return
	}

	d := common.RedisConn.Get()
	defer d.Close()

	//查看哈希列表中是否存在
	ok, err := d.Do(models.REDIS_HEXISTS,models.USER_LIST,item+"_"+user_id)
	if err != nil {
		logs.Info("查询redis错误,Err:%v", err)
		c.Data["json"] = models.ApiResponse{Code: models.ErrSystem, Message: "redis保存账号错误1", Data: ""}
		c.ServeJSON()
		return
	}

	if ok.(int64)==1 {
		logs.Info("账号存在")
		c.Data["json"] = models.ApiResponse{Code: models.ErrSystem, Message: "账号存在", Data: ""}
		c.ServeJSON()
		return
	}

	strNowTime:=strconv.FormatInt(time.Now().Unix(),10)
	_, err1 := d.Do(models.REDIS_HSET,models.USER_LIST, item+"_"+user_id, strNowTime)
	if err1 != nil {
		logs.Info("redis保存账号错误, Err:%v", err1)
		c.Data["json"] = models.ApiResponse{Code: models.ErrSystem, Message: "redis保存账号错误", Data: ""}
		c.ServeJSON()
		return
	}

	logs.Info("redis保存账号成功")
	c.Data["json"] = models.ApiResponse{Code: models.SUCCESS, Message: "redis保存账号成功", Data: ""}
	c.ServeJSON()
	return
}

//获取群消息
func (c *UserApiController) GetGroupMessage() {
	userId := c.GetString("user_id")
	groupId := c.GetString("group_id")
	page, _ := c.GetInt("page")
	if groupId=="" || userId=="" {
		logs.Info("参数错误,请检查")
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrParam,
			Message: "参数错误,请检查",
			Data: "",
		}
		c.ServeJSON()
		return
	}
	d := common.RedisConn.Get()
	defer d.Close()

	//查看哈希列表中是否存在
	ok, err := d.Do(models.REDIS_HEXISTS,models.USER_LIST,userId)
	if err != nil {
		logs.Info("查询redis错误,Err:%v", err)
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrSystem,
			Message: "redis获取失败",
			Data: "",
		}
		c.ServeJSON()
		return
	}
	if ok.(int64)!=1 {
		logs.Info("账号不存在")
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrSystem,
			Message: "账号不存在",
			Data: "",
		}
		c.ServeJSON()
		return
	}

	//查找该用户用有的群
	groupDdat:=GetUidToGroupDataRedis(userId)
	if !strings.Contains(groupDdat,groupId) {
		logs.Info("你没有该群权限")
		c.Data["json"] = models.ApiResponse{
			Code: models.ErrSystem,
			Message: "你没有该群权限",
			Data: "",
		}
		c.ServeJSON()
		return
	}

	//获取群聊天历史记录
	limit,_ :=beego.AppConfig.Int("max_history_msg_record")
	data:= models.GetGroupChatRecord(groupId,page,limit)

	c.Data["json"] = models.ApiResponse{
		Code: models.SUCCESS,
		Message: "ok",
		Data: data,
		Page: strconv.Itoa(page),
	}
	c.ServeJSON()
	return
}

//验证文件上传上来的信息是否合法
func (c *UserApiController) VerifyUpFileInfo() {
	//先查看之前是否有关系
	dataBody:=c.Ctx.Input.RequestBody
	data:=string(dataBody)
	data=data[1:len(data)-1]
	logs.Info("VerifyUpFileInfo文件上传数据:",data)
	messageStr,_ := extra.Base64DecodeString(data)
	logs.Info("VerifyUpFileInfo文件上传信息:",messageStr)
	fromData,_:=extra.JsonToFromMessageStruct(messageStr)
	_,ok:=router.Load(fromData.Method)
	if ok && (models.OneToOne==fromData.Method || models.ToGroup==fromData.Method || models.Radio==fromData.Method) {
		//查看接收者和发送者是否相同
		if fromData.ToUid==fromData.Uid {
			c.Data["json"] = models.ApiResponse{
				Code: -6,
				Message: "无法给自己发送",
			}
			c.ServeJSON()
			return
		}
		//查看发送者是否设置uid
		_,uok:=uidToKey.Load(fromData.Uid)
		if !uok {//发送者未设置uid或者说未登录
			c.Data["json"] = models.ApiResponse{
				Code: -2,
				Message: "发送者未登录",
			}
			c.ServeJSON()
			return
		}
		if fromData.Method==models.OneToOne {
			//查看接收者是否在线
			_,tuok:=uidToKey.Load(fromData.ToUid)
			if !tuok {
				c.Data["json"] = models.ApiResponse{
					Code: -3,
					Message: "接受者不在线或者不存在",
				}
				c.ServeJSON()
				return
			}
		}
		if fromData.Method==models.ToGroup {
			//查看群是否存在
			groupData:=GetGroupDataToRedis(fromData.ToGroupId)
			if groupData=="" {//群id不能为空或者群不存在
				c.Data["json"] = models.ApiResponse{
					Code: -4,
					Message: "群id不能为空或者群不存在",
				}
				c.ServeJSON()
				return
			}

			/**检测发送者是否在该群中*/
			//将groupData解析成数组
			var groupDataArr []models.GroupMembers
			json.Unmarshal([]byte(groupData), &groupDataArr)
			inGroup:=false
			for _,gV := range groupDataArr {
				if fromData.Uid==gV.Uid {
					inGroup = true
					break
				}
			}
			if !inGroup {
				c.Data["json"] = models.ApiResponse{
					Code: -5,
					Message: "无权限",
				}
				c.ServeJSON()
				return
			}
			/**检测发送者是否在该群中*/
		}
		//验证通过
		c.Data["json"] = models.ApiResponse{
			Code: 1,
			Message: "成功",
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.ApiResponse{
		Code: -1,
		Message: "参数不合法",
	}
	c.ServeJSON()
	return

}

//推送消息
func (c *UserApiController) PushMsg() {
	dataBody:=c.Ctx.Input.RequestBody
	data:=string(dataBody)
	data=data[1:len(data)-1]
	logs.Info("PushMsg文件上传数据:",data)
	messageStr,_ := extra.Base64DecodeString(data)
	logs.Info("PushMsg文件上传信息:",messageStr)
	fromData,_:=extra.JsonToFromMessageStruct(messageStr)
	_,ok:=router.Load(fromData.Method)
	if ok && (models.OneToOne==fromData.Method || models.ToGroup==fromData.Method || models.Radio==fromData.Method) {
		//查看接收者和发送者是否相同
		if fromData.ToUid==fromData.Uid {
			c.Data["json"] = models.ApiResponse{
				Code: -6,
				Message: "无法给自己发送",
			}
			c.ServeJSON()
			return
		}
		//查看发送者是否设置uid
		key,uok:=uidToKey.Load(fromData.Uid)
		if !uok {//发送者未设置uid或者说未登录
			c.Data["json"] = models.ApiResponse{
				Code: -2,
				Message: "发送者未登录",
			}
			c.ServeJSON()
			return
		}
		if fromData.Method==models.OneToOne {
			//查看接收者是否在线
			_,tuok:=uidToKey.Load(fromData.ToUid)
			if !tuok {
				c.Data["json"] = models.ApiResponse{
					Code: -3,
					Message: "接受者不在线或者不存在",
				}
				c.ServeJSON()
				return
			}
		}
		if fromData.Method==models.ToGroup {
			//查看群是否存在
			groupData:=GetGroupDataToRedis(fromData.ToGroupId)
			if groupData=="" {//群id不能为空或者群不存在
				c.Data["json"] = models.ApiResponse{
					Code: -4,
					Message: "群id不能为空或者群不存在",
				}
				c.ServeJSON()
				return
			}

			/**检测发送者是否在该群中*/
			//将groupData解析成数组
			var groupDataArr []models.GroupMembers
			json.Unmarshal([]byte(groupData), &groupDataArr)
			inGroup:=false
			for _,gV := range groupDataArr {
				if fromData.Uid==gV.Uid {
					inGroup = true
					break
				}
			}
			if !inGroup {
				c.Data["json"] = models.ApiResponse{
					Code: -5,
					Message: "无权限",
				}
				c.ServeJSON()
				return
			}
			/**检测发送者是否在该群中*/
		}
		//获取verifyUser
		verifyUser:=beego.AppConfig.String("verify_user")

		methodStruct,ok:=router.Load(fromData.Method)
		if ok {
			methodOb:=methodStruct.(MethodStruct)
			if methodOb.BeforeMethod!=nil {//存在之前函数
				BeforeExeResult:=methodOb.BeforeMethod(key.(string),verifyUser,fromData)
				if BeforeExeResult {
					c.Data["json"] = models.ApiResponse{
						Code: -6,
						Message: "无权限",
					}
					c.ServeJSON()
					return
				}
			}
			if methodOb.Method!=nil {//真实执行方法
				exeResult:=methodOb.Method(key.(string),verifyUser,fromData)
				if exeResult {
					c.Data["json"] = models.ApiResponse{
						Code: 1,
						Message: "发送成功",
					}
					c.ServeJSON()
					return
				}
			}
		}

		c.Data["json"] = models.ApiResponse{
			Code: -7,
			Message: "未找到方法",
		}
		c.ServeJSON()
		return
	}
	c.Data["json"] = models.ApiResponse{
		Code: -1,
		Message: "参数不合法",
	}
	c.ServeJSON()
	return
}