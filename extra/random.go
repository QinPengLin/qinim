package extra

import (
	"math/rand"
	"strconv"
	"time"
)

//获取随机数
func GetRand(s,e int)  int {
	if e<s{
		e=s
	}
	//step1: 设置种子数
	rand.Seed(time.Now().UnixNano())
	//step2：获取随机数
	num := rand.Intn(e-s+1) + s
	return num
}

//创建唯一KEY
func CreateOnlKey() string {
	key:=time.Now().UnixNano()
	keyStr:=strconv.FormatInt(key,10)+strconv.Itoa(GetRand(100,999))
	return keyStr
}
