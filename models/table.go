package models

import (
	"github.com/astaxie/beego/orm"
)

const (
	MYWS_GROUP_CHAT_RECORD  string = "group_chat_record"  // 群聊消息记录表名
)

func init() {
	orm.RegisterModel(new(GroupChatRecord),)
}

func Initializes()  {
	CheckGroupChatRecord()
}

func (gcr *GroupChatRecord) TableName() string {
	return MYWS_GROUP_CHAT_RECORD
}

type GroupChatRecord struct {
	Id             int32     `json:"id,omitempty"`
	GroupId        string    `orm:"column(group_id)" json:"group_id,omitempty"`
	FormUid        string    `orm:"column(form_uid)" json:"form_uid,omitempty"`
	Content        string    `orm:"column(content)" json:"content,omitempty"`
	CreationTime   string    `orm:"column(creation_time)" json:"creation_time,omitempty"`
}

//检测某个表是否存在
func CheckMywsTable(tableName string) int {
	o := orm.NewOrm()

	sql := "SHOW TABLES LIKE '" + tableName + "'"
	var a []orm.Params
	o.Raw(sql).Values(&a)

	return len(a)
}

//检测group_chat_record表是否存在，不存在就新增（程序每次启动都会执行一次）
func CheckGroupChatRecord()  {
	if  CheckMywsTable(MYWS_GROUP_CHAT_RECORD) == 0 {
		o := orm.NewOrm()
		create :="CREATE TABLE `"+MYWS_GROUP_CHAT_RECORD+"` ("+
			"`id`  bigint(20) NOT NULL AUTO_INCREMENT ,"+
			"`group_id`  varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '群id' ,"+
			"`form_uid`  varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '发送者uid' ,"+
			"`content`  text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '发送的消息' ,"+
			"`creation_time`  bigint(14) NULL COMMENT '创建时间' ,"+
			"PRIMARY KEY (`id`),"+
			"INDEX `group_id` (`group_id`) USING HASH ,"+
			"INDEX `form_uid` (`form_uid`) USING HASH)"
		o.Raw(create).Exec()
	}
}