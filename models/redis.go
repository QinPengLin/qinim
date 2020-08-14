package models

// redis命令
const (
	REDIS_RPOPLPUSH = "RPOPLPUSH"
	REDIS_LPUSH     = "LPUSH"
	REDIS_RPUSH     = "RPUSH"
	REDIS_LLEN      = "LLEN"
	REDIS_HSET      = "HSET"
	REDIS_HLEN      = "HLEN"
	REDIS_HGET      = "HGET"
	REDIS_HMGET     = "HMGET"
	REDIS_HMSET     = "HMSET"
	REDIS_HEXISTS   = "HEXISTS"
	REDIS_SET       = "SET"
	REDIS_GET       = "GET"
	REDIS_EXISTS    = "EXISTS"
	REDIS_SADD      = "SADD"
	REDIS_SISMEMBER      = "SISMEMBER"
)

// redis对应的key
const (

	//用户列表
	USER_LIST                = "user_list"
	//群列表
	GROUP_LIST               = "group_list"
	//用户拥有的群
	USER_TO_GROUPS           = "user_to_groups"
	//群id=>群名
	GROUPID_TO_GROUPNAME     = "groupid_to_groupname"

)
