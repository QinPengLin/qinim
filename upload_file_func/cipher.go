package upload_file_func

import "encoding/base64"

// Base64EncodeString 编码
func Base64EncodeString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64DecodeString 解码
func Base64DecodeString(str string) (string, []byte) {
	resBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(resBytes), resBytes
}
