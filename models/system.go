package models

const (
	DispersedSwitch             = "yes"          //分布式开关（yes为开启，其他为关闭）
	VerifyUser                  = "yes"          //是否检查设置用户，并且检查设置用户是否合法yes为开启其他为关闭
	AddUid                      = "add_uid"      //新增uid的method参数
	ModifyUid                   = "modify_uid"   //修改uid的method参数
	CreateGroup                 = "create_group" //创建的method参数
	ToGroup                     = "to_group"     //群聊的method参数
	HEART                       = "heart"        //心跳的method参数
	OneToOne                    = "to_one"       //私聊method参数
	Radio                       = "all"          //广播给所有的method参数
)