package app

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"go-web/internal/model"
	"go-web/utils/ecode"
	"math"
	"strings"
)

type Context struct {
	Ctx      *gin.Context
	handlers Handlers
	index    int8

	User     *model.User
	Passport *model.Passport
	Request  *Request
}

type Request struct {
	Resource string
}

type Response struct {
	Code     int
	Message  string
	Success  bool
	Resource string
}

func (ctx *Context) Abort() {
	ctx.index = math.MaxInt8 / 2
}

func (ctx *Context) Next() {
	ctx.index++
	handlers := ctx.handlers
	for ctx.index < int8(len(handlers)) {
		handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) Status(code int) {
	ctx.Ctx.Status(code)
}

func (ctx *Context) JSON(obj interface{}, err error) {
	code := ecode.Cause(err)
	var data []byte
	if err == nil {
		var err1 error
		data, err1 = json.Marshal(&obj)
		if err1 != nil {
			fmt.Printf("app->Context->JSON->json.Marshal Err:%v\n", err1)
			err = err1
		}
	}
	ctx.Ctx.JSON(200, &Response{
		Code:     code.Code(),
		Message:  code.Message(),
		Success:  err == nil,
		Resource: string(data),
	})
}

func (ctx *Context) PostForm(key string) (value string) {
	return ctx.Ctx.PostForm(key)
}

func (ctx *Context) ClientIP() string {
	return ctx.Ctx.ClientIP()
}

func (ctx *Context) Query(key string) (value string) {
	value = ctx.Ctx.Query(key)
	return
}

func (ctx *Context) GetHeader(key string) string {
	return ctx.Ctx.GetHeader(key)
}

func (ctx *Context) SetCookies(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	ctx.Ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (ctx *Context) Cookie(name string) (string, error) {
	return ctx.Ctx.Cookie(name)
}

func (ctx *Context) String(code int, format string, values ...any) {
	ctx.Ctx.String(code, format, values...)
}

func (ctx *Context) ShouldBindJSON(obj interface{}) error {
	err := json.Unmarshal([]byte(ctx.Request.Resource), &obj)
	if err != nil {
		return ecode.FormatError
	}
	return err
}

func (ctx *Context) ParseRequestResource(path string) gjson.Result {
	return gjson.Get(ctx.Request.Resource, path)
}

func (ctx *Context) ParseRequestResourceMany(path ...string) []gjson.Result {
	return gjson.GetMany(ctx.Request.Resource, path...)
}

func (ctx *Context) IsLogin() bool {
	return ctx.Passport.Uid != 0
}

func (ctx *Context) ParseUserAgent() (ua string, deviceType, deviceName string) {
	ua = ctx.Ctx.GetHeader("User-Agent")
	if strings.Contains(ua, "Mobile") || strings.Contains(ua, "Android") || strings.Contains(ua, "iPhone") {
		deviceType = "Mobile"
		if strings.Contains(ua, "Android") {
			deviceName = "Android Device"
		} else if strings.Contains(ua, "iPhone") {
			deviceName = "iPhone"
		}
	} else if strings.Contains(ua, "iPad") {
		deviceType = "Tablet"
		deviceName = "iPad"
	} else {
		deviceType = "Desktop"
		if strings.Contains(ua, "Windows") {
			deviceName = "Windows PC"
		} else if strings.Contains(ua, "Macintosh") {
			deviceName = "Mac"
		} else if strings.Contains(ua, "Linux") {
			deviceName = "Linux PC"
		} else {
			deviceName = "未知设备"
		}
	}
	return
}

// IsAdmin 此方法在ctx.User初始化后可用
func (ctx *Context) IsAdmin() bool {
	return ctx.User.Role >= 2
}
