package extra

import (
	"encoding/json"
	"qinim/models"
)

func JsonToFromMessageStruct(jsonStr string) (*models.FromMessage,error)  {
	var s models.FromMessage

	err:=json.Unmarshal([]byte(jsonStr), &s)
	if err!=nil{
		return &s,err
	}
	return &s,err
}

func StructToFromMessageJson(structs *models.SendMessageShelf) (string,error) {
	data, err := json.Marshal(structs)
	return string(data),err
}

func StructToFromMeMessageJson(structs *models.SendToMeMessageShelf) (string,error) {
	data, err := json.Marshal(structs)
	return string(data),err
}