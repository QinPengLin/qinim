package models

//获取客户端结构
type FromMessage struct {
	Method                string         `json:"method,omitempty"`              //方法
	Uid                   string         `json:"uid,omitempty"`                 //发送者id(当方法为modify_uid是修改uid，add_uid是新增uid，发送文件时需要传入发送者uid,普通消息不需要发送uid)
	ToUid                 string         `json:"to_uid,omitempty"`              //发送给谁(为all是广播给所有)
	ToGroupId             string         `json:"to_group_id,omitempty"`         //发送给那个群ID
	GroupInitMembers      string         `json:"group_init_members,omitempty"`  //创建群初始成员用户uid,如果为空就将只有群主，初始成员个数上限为群人数上限
	GroupName             string         `json:"group_name,omitempty"`          //群名称，创建或者修改群的时候需要
	Content               string         `json:"content,omitempty"`             //内容
	File                  []string       `json:"file_url,omitempty"`            //上传的文件地址
}

//发送给客户端状态码的结构体
type SendMessageShelf struct {
	Code          string          `json:"code,omitempty"`                //状态码（1无异常，）
	Msg           string          `json:"msg,omitempty"`                 //信息
	Data          *SendMessage    `json:"data,omitempty"`                //数据
}

//发送给客户端的信息结构
type SendMessage struct {
	GroupName       string                  `json:"group_name,omitempty"`         //群名
	GroupId         string                  `json:"group_id,omitempty"`           //群id
	ToGroupId       string                  `json:"to_group_id,omitempty"`        //发给群id
	FromName        string                  `json:"from_name,omitempty"`          //谁发送的名字
	FromUid         string                  `json:"from_id,omitempty"`            //谁发送的id
	ToUid           string                  `json:"to_id,omitempty"`              //发给谁的id
	Time            string                  `json:"time,omitempty"`               //时间
	Msg             string                  `json:"msg,omitempty"`                //内容
	File            []string                `json:"file_url,omitempty"`           //上传的文件地址
}

//发送给客户端自己包含数据和状态码的结构体
type SendToMeMessageShelf struct {
	Code          string              `json:"code,omitempty"`                //状态码（1无异常，）
	Msg           string              `json:"msg,omitempty"`                 //信息
	Data          *SendToMeMessage    `json:"data,omitempty"`                //数据
}

//发送给客户端自己的信息结构
type SendToMeMessage struct {
	SuccessUid       string                  `json:"success_uid,omitempty"`         //创建群时成功初始化的uid
	FailUid          string                  `json:"fail_uid,omitempty"`            //创建群时初始化失败的uid
	GroupId          string                  `json:"group_id,omitempty"`            //创建群成功后返回的群id
}

//群成员保存到redis的结构体
type GroupMembers struct {
	Uid      string      `json:"uid,omitempty"`                //用户ID
	Level    string      `json:"level,omitempty"`              //级别(1:群主2:普通成员)
}

//api接口返回
type ApiResponse struct {
	Code    int             `json:"code,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
	Page    string          `json:"page,omitempty"`
}