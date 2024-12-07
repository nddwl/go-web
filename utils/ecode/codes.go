package ecode

var (
	OK        = New(0, "成功")
	ServerErr = New(-500, "服务器繁忙")

	BadRequest   = New(-400, "无效请求")
	TokenExpired = New(-401, "登录过期")
	Forbidden    = New(-403, "无权访问")

	UserEmailRegistered = New(301001, "邮箱已注册")
	UserNameExisted     = New(301002, "名称已存在")
	UserLoginFailed     = New(301003, "用户名或密码错误")
	FormatError         = New(301004, "格式错误")
	UserSigned          = New(301005, "今日已签到")
	UserNotExist        = New(301006, "用户不存在")

	ActivityIsOver = New(401001, "活动不存在或已经结束")
	ActivityIsList = New(401002, "活动已上架暂无法修改")
)
