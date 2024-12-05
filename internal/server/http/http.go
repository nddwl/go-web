package http

import (
	"go-web/internal/service"
	"go-web/middleware"
	"go-web/utils/app"
	"go-web/utils/config"
)

type Server struct {
	Service *service.Service
	Group
	App *app.App
}

type Group struct {
	User     *User
	Passport *Passport
	Activity *Activity
}

func New() *Server {
	server := &Server{
		Service: service.New(),
		App:     app.New(),
	}
	server.init()
	server.initGroup()
	return server
}

func (t *Server) init() {
	t.Group = Group{
		User:     NewUser(t),
		Passport: NewPassport(t),
		Activity: NewActivity(t),
	}
	t.App.Gin.Use(middleware.Cors())
	t.App.Use(t.Passport.GetPassport)
}

func (t *Server) initGroup() {
	t.User.initGroup()
	t.Passport.initGroup()
	t.Activity.initGroup()
}

func (t *Server) Run() error {
	return t.App.Gin.Run(config.Server.Addr)
}
