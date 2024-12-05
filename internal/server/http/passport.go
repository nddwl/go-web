package http

import (
	"go-web/utils/app"
	"go-web/utils/ecode"
)

type Passport struct {
	*Server
	Group *app.RouterGroup
}

func NewPassport(server *Server) *Passport {
	passport := &Passport{Server: server}
	return passport
}

func (t *Passport) initGroup() {
	t.Group = t.App.Group("/passport")
	t.Group.POST("/create", t.Create)
	t.Group.POST("/delete", t.Delete)
	t.Group.POST("/find", t.Find)
}

func (t *Passport) Create(ctx *app.Context) {
	//第三方登录
	ctx.JSON("暂未实现", nil)
}

func (t *Passport) Delete(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	result := ctx.ParseRequestResource("device_id")
	err := t.Service.Passport.Delete(ctx.Passport.Uid, result.String())
	ctx.JSON(nil, err)
}

func (t *Passport) Find(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	m, err := t.Service.Passport.Find(ctx.Passport.Uid)
	ctx.JSON(&m, err)
}

func (t *Passport) GetPassport(ctx *app.Context) {
	if ctx.Passport.Token != "" {
		dto, err := t.Service.Passport.GetPassport(ctx.Passport.Token)
		if err != nil {
			ctx.JSON(nil, err)
			ctx.Abort()
			return
		}
		ctx.Passport.Uid = dto.Uid
		ctx.Passport.DeviceId = dto.DeviceId
	}
}
