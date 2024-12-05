package app

import (
	"github.com/gin-gonic/gin"
	"go-web/internal/model"
	"go-web/utils/ecode"
)

type App struct {
	Gin *gin.Engine
	RouterGroup
}

func New() *App {
	app := &App{
		Gin: gin.Default(),
	}
	app.RouterGroup = RouterGroup{
		ginRouterGroup: &app.Gin.RouterGroup,
		app:            app,
		root:           true,
	}
	app.Use(app.parse)
	return app
}

func (app *App) parse(ctx *Context) {
	token, err := ctx.Cookie("__PASSPORT")
	if err == nil {
		ctx.Passport = &model.Passport{Token: token}
	}
	ctx.Passport.Ip = ctx.ClientIP()
	ctx.Passport.Ua = ctx.GetHeader("User-Agent")
	if ctx.GetHeader("Content-Type") == "application/json" {
		err = ctx.Ctx.ShouldBindJSON(&ctx.Request)
		if err != nil {
			ctx.JSON(nil, ecode.ServerErr)
			ctx.Abort()
		}
	} else {
		ctx.Request = &Request{}
	}
}
