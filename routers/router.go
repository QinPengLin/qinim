package routers

import (
	"qinim/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WsController{})

	beego.Router("/set_user", &controllers.UserApiController{}, "post:SetUser")  //设置用户
	beego.Router("/verify_up_fileInfo", &controllers.UserApiController{}, "post:VerifyUpFileInfo")  //文件上传验证信息
	beego.Router("/push_msg", &controllers.UserApiController{}, "post:PushMsg")  //推送消息
	beego.Router("/get_group_message", &controllers.UserApiController{}, "post:GetGroupMessage")  //获取群历史记录

	beego.Router("/up_file", &controllers.UploadFileController{}, "post:UploadMsgFile")  //文件上传
}
