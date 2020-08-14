package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"os"
	"path"
	"qinim/upload_file_func"
	"qinim/upload_file_func/sendhttp"
	"time"
)

type UploadFileController struct {
	beego.Controller
}

//独立的包括消息内容文件上传（上传成功后会通过接口推送给im系统，im系统再转发给需要推送的人）
func (this *UploadFileController) UploadMsgFile() {
	data := this.GetString("data")//必填
	//最多支持一次性上传三个文件
	f, h, err := this.GetFile("file_1")//获取上传的1文件(目前只做了单文件上次，多文件需要考虑文件上传出来事务性)
	//f_2, h_2, err_3 := this.GetFile("file_2")//获取上传的2文件
	//f_3, h_3, err_3 := this.GetFile("file_3")//获取上传的3文件

	if data!="" && err==nil {//先检查data是否为空
		messageStr,_ := upload_file_func.Base64DecodeString(data)
		logs.Info("文件上传信息:",messageStr)
		//在通过接口访问im系统获取发送者是否登陆和如果是私聊或者群聊是否有该用户或者群
		fileToIm:=beego.AppConfig.String("file_to_im")
		verifyUpFileInfoJson:=sendhttp.Post(fileToIm+"/verify_up_fileInfo",data,"application/json; charset=utf-8")
		verifyUpFileInfo,_:=upload_file_func.JsonToApiResponseStruct(verifyUpFileInfoJson)
		if(verifyUpFileInfo.Code!=1){//验证不通过
			f.Close()
			this.Data["json"] = upload_file_func.ApiResponse{
				Code: verifyUpFileInfo.Code,
				Message: verifyUpFileInfo.Message,
			}
			this.ServeJSON()
			return
		}

		fromData,_:=upload_file_func.JsonToFromMessageStruct(messageStr)

		ext := path.Ext(h.Filename)
		//验证后缀名是否符合要求
		var AllowExtMap  = map[string]bool{
			".png":true,
			".mp4":true,
			".pge":true,
			".txt":true,
		}
		if _,ok:=AllowExtMap[ext];!ok{//返回文件不合法
			f.Close()
			this.Data["json"] = upload_file_func.ApiResponse{
				Code: -11,
				Message: "文件不合法",
			}
			this.ServeJSON()
			return
		}
		//创建目录
		uploadDir := "static/upload/" + time.Now().Format("2006/01/02/")
		err := os.MkdirAll( uploadDir , 775)
		if err != nil {//返回文件上传失败
			f.Close()
			this.Data["json"] = upload_file_func.ApiResponse{
				Code: -12,
				Message: "文件上传失败",
			}
			this.ServeJSON()
			return
		}

		//构造文件名称
		str_time :=time.Now().Format("20060102150405")
		fileName := str_time + ext

		fpath := uploadDir + fileName
		err = this.SaveToFile("file_1", fpath)
		if err != nil {//文件上传失败
			f.Close()
			this.Data["json"] = upload_file_func.ApiResponse{
				Code: -13,
				Message: "文件上传失败",
			}
			this.ServeJSON()
			return
		}
		//文件上传成功返回上传成功信息并且通知给IM系统
		f.Close()

		fromData.File=append(fromData.File,beego.AppConfig.String("file_domain")+uploadDir+fileName)
		sendData,_:=upload_file_func.StructToFromMessageJson(fromData)
		encryptSendData:=upload_file_func.Base64EncodeString(sendData)
		pushMsgJson:=sendhttp.Post(fileToIm+"/push_msg",encryptSendData,"application/json; charset=utf-8")
		pushMsg,_:=upload_file_func.JsonToApiResponseStruct(pushMsgJson)
		if(pushMsg.Code!=1){//发送失败
			this.Data["json"] = upload_file_func.ApiResponse{
				Code: pushMsg.Code,
				Message: pushMsg.Message,
			}
			this.ServeJSON()
			return
		}
		this.Data["json"] = upload_file_func.ApiResponse{
			Code: 1,
			Message: "发送成功",
		}
		this.ServeJSON()
		return
	}
	this.Data["json"] = upload_file_func.ApiResponse{
		Code: -14,
		Message: "参数错误",
	}
	this.ServeJSON()
	return
}