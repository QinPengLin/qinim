package common

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
	"time"
)

var RedisConn *redis.Pool

func RegisterRedis() *redis.Pool {
	//获取数据库配置文件
	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   10,
		IdleTimeout: 20 * time.Second,
		Dial: func() (conn redis.Conn, e error) {
			if e!=nil {
				logs.Info("redis连接失败, Err:%v", e)
			}
			return redis.DialURL(beego.AppConfig.String("redis"))
		},
	}
}
