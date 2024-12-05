package http

import (
	"encoding/json"
	"go-web/internal/model"
	"go-web/utils"
	"go-web/utils/app"
	"go-web/utils/ecode"
)

type User struct {
	*Server
	Group *app.RouterGroup
}

func NewUser(server *Server) *User {
	user := &User{Server: server}
	return user
}

func (t *User) initGroup() {
	t.Group = t.App.Group("/user")
	t.Group.POST("/create", t.Create)
	t.Group.POST("/delete", t.Delete)
	t.Group.POST("/update", t.Update)
	t.Group.GET("/find", t.Find)
	t.Group.GET("/login", t.Login)
	t.Group.POST("/auth", t.Auth)
	t.Group.POST("/logout", t.Logout)
	t.Group.GET("/isNameExists", t.IsNameExists)
	t.Group.POST("/sign", t.Sign)
	t.Group.GET("/findActivityRecord", t.FindActivityRecord)
}

func (t *User) Create(ctx *app.Context) {
	data, err := utils.Decrypt(ctx.Request.Resource)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	user := model.UserCreate{}
	err = json.Unmarshal(data, &user)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	_, err = t.Service.User.Create(user)
	ctx.JSON(nil, err)
}

func (t *User) Delete(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	_, err := t.Service.User.Delete(ctx.Passport.Uid)
	ctx.JSON(nil, err)
}

func (t *User) Update(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	dto := model.UserUpdate{}
	err := ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	err = t.Service.User.Update(ctx.Passport.Uid, dto.Key, dto.Value)
	ctx.JSON(nil, err)
}

func (t *User) Find(ctx *app.Context) {
	result := ctx.ParseRequestResource("name")
	m, err := t.Service.User.Find(result.String())
	ctx.JSON(&m, err)
}

func (t *User) Auth(ctx *app.Context) {
	data, err := utils.Decrypt(ctx.Request.Resource)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	passwordDto := model.PasswordDto{}
	err = json.Unmarshal(data, &passwordDto)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	session, err := ctx.Cookie("SESS__LOGIN")
	if err != nil {
		ctx.JSON(nil, ecode.BadRequest)
	}
	passwordDto.Session = session
	uid, err := t.Service.User.Login(passwordDto)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	m, err := t.Service.Passport.Create(model.Passport{
		Uid: uid,
		Ip:  ctx.ClientIP(),
		Ua:  ctx.GetHeader("User-Agent"),
	})
	ctx.SetCookies("SESS__LOGIN", "", 0, "/user/auth", "127.0.0.1", false, true)
	ctx.SetCookies("__PASSPORT", m.Token, 60*60*24*30, "/", "127.0.0.1", false, true)
	ctx.JSON(nil, err)
}

func (t *User) Login(ctx *app.Context) {
	id, b64s, _, err := t.Service.Rdb.Captcha.Generate()
	if err != nil {
		ctx.JSON(nil, err)
	}
	ctx.SetCookies("SESS__LOGIN", id, -1, "/user/auth", "", false, true)
	ctx.String(200, "%s", b64s)
}

func (t *User) Logout(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	err := t.Service.Passport.Delete(ctx.Passport.Uid, ctx.Passport.DeviceId)
	ctx.SetCookies("__PASSPORT", "", 0, "/", "127.0.0.1", false, true)
	ctx.JSON(nil, err)
}

func (t *User) IsNameExists(ctx *app.Context) {
	result := ctx.ParseRequestResource("name")
	exists, err := t.Service.User.IsNameExists(result.String())
	ctx.JSON(exists, err)
}

func (t *User) Sign(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	m, err := t.Service.User.Sign(ctx.Passport.Uid)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	ctx.JSON(m, nil)
}

func (t *User) GetUser(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		ctx.Abort()
		return
	}
	user, err := t.Service.User.GetUser(ctx.Passport.Uid)
	if err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
	ctx.User = &user
}

func (t *User) FindActivityRecord(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	m, err := t.Service.Activity.FindRecord(ctx.Passport.Uid)
	ctx.JSON(&m, err)
}
