package extra

import "strings"

//字符分割
func StrPartition(str,need string) ([]string,int)  {
	cleanStr:=strings.Trim(str, need)
	re:=strings.Split(cleanStr,need)
	re=RemoveRepeatedElement(re)
	len:=len(re)
	return  re,len
}

//数组去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
