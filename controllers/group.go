package controllers

import (
	"qinim/common"
	"qinim/models"
	"github.com/garyburd/redigo/redis"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strings"
)

//创建群并且维护uid=>Groupids的关系
func AddGroupToRdis(groupId string,groupData []models.GroupMembers,groupName string) (bool,error) {
	//获取createUids在user_to_groups中的所有信息
	d := common.RedisConn.Get()
	defer d.Close()

	//循环获取
	existUserToGroup:=make(map[string]string)
	for _,groupDataV := range groupData {
		ok, _ := redis.String(d.Do(models.REDIS_HGET,models.USER_TO_GROUPS,groupDataV.Uid))
		if ok=="" {
			existUserToGroup[groupDataV.Uid] = groupId
		}else {
			existUserToGroup[groupDataV.Uid] = ok + "|" + groupId
		}
		existUserToGroup[groupDataV.Uid] = strings.Trim(existUserToGroup[groupDataV.Uid], "|")
	}

	//新增群
	groupMemberDataJson, _ := json.Marshal(groupData)
	groupMemberDataJsonStr:=string(groupMemberDataJson)
	_, hsetErr := d.Do(models.REDIS_HSET,models.GROUP_LIST, groupId, groupMemberDataJsonStr)
	if hsetErr != nil {
		logs.Info("AddGroupToRdis保存账号错误, hsetErr:%v", hsetErr)
		return false,hsetErr
	}
	//新增群id=>群名
	_, hsetToNameErr := d.Do(models.REDIS_HSET,models.GROUPID_TO_GROUPNAME, groupId, groupName)
	if hsetToNameErr != nil {
		logs.Info("AddGroupToRdis保存账号错误, hsetToNameErr:%v", hsetToNameErr)
		return false,hsetToNameErr
	}
	//批量修改或者新增
	_, hmsetErr := d.Do(models.REDIS_HMSET, redis.Args{}.Add(models.USER_TO_GROUPS).AddFlat(existUserToGroup)...)
	if hmsetErr!=nil {
		logs.Info("AddGroupToRdis查询redis错误, hmsetErr:%v", hmsetErr)
		return false,hmsetErr
	}
	return true,nil
}

//获取指定群ID是否在redis中
func GetGroupDataToRedis(group_id string) string {
	d := common.RedisConn.Get()
	defer d.Close()
	data, _ := redis.String(d.Do(models.REDIS_HGET,models.GROUP_LIST,group_id))
	return data
}

//获取指uid对应的群id数据字符
func GetUidToGroupDataRedis(uid string) string {
	d := common.RedisConn.Get()
	defer d.Close()
	data, _ := redis.String(d.Do(models.REDIS_HGET,models.USER_TO_GROUPS,uid))
	return data
}