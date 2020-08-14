package upload_file_func

//获取客户端结构
type FromMessage struct {
	Method                string         `json:"method,omitempty"`              //方法
	Uid                   string         `json:"uid,omitempty"`                 //发送者id(当方法为modify_uid是修改uid，add_uid是新增uid，发送文件时需要传入发送者uid,普通消息不需要发送uid)
	ToUid                 string         `json:"to_uid,omitempty"`              //发送给谁(为all是广播给所有)
	ToGroupId             string         `json:"to_group_id,omitempty"`         //发送给那个群ID
	GroupInitMembers      string         `json:"group_init_members,omitempty"`  //创建群初始成员用户uid,如果为空就将只有群主，初始成员个数上限为群人数上限
	GroupName             string         `json:"group_name,omitempty"`          //群名称，创建或者修改群的时候需要
	Content               string         `json:"content,omitempty"`             //内容
	File                  []string       `json:"file_url,omitempty"`            //上传的文件地址
}

//api接口返回
type ApiResponse struct {
	Code    int             `json:"code,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Message string          `json:"message,omitempty"`
	Page    string          `json:"page,omitempty"`
}
