package models

import (
	"github.com/astaxie/beego/orm"
)

func InsertGroupChatRecord(data *GroupChatRecord)  int64 {
	o := orm.NewOrm()

	id, err := o.Insert(data)
	if err == nil {
		return id
	}
	return 0
}

func GetGroupChatRecord(group_id string,page int,limit int)  []*GroupChatRecord{
	o := orm.NewOrm()

	if page<0 {
		page=0
	}
	offset:=page*limit

	var data []*GroupChatRecord
	_, err :=o.QueryTable(MYWS_GROUP_CHAT_RECORD).Filter("group_id", group_id).Limit(limit, offset).All(&data)
	if err!=nil {
		return nil
	}
	return data
}