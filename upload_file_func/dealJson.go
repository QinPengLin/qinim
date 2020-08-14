package upload_file_func

import (
	"encoding/json"
)

func JsonToFromMessageStruct(jsonStr string) (*FromMessage,error)  {
	var s FromMessage

	err:=json.Unmarshal([]byte(jsonStr), &s)
	if err!=nil{
		return &s,err
	}
	return &s,err
}

func JsonToApiResponseStruct(jsonStr string) (*ApiResponse,error)  {
	var s ApiResponse

	err:=json.Unmarshal([]byte(jsonStr), &s)
	if err!=nil{
		return &s,err
	}
	return &s,err
}

func StructToFromMessageJson(structs *FromMessage) (string,error) {
	data, err := json.Marshal(structs)
	return string(data),err
}
