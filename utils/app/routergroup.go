package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterGroup struct {
	ginRouterGroup *gin.RouterGroup
	handlers       Handlers
	app            *App
	root           bool
}

type IRouter interface {
	IRoutes
	Group(string, ...Handler) *RouterGroup
}

type IRoutes interface {
	Use(...Handler) IRoutes

	Handle(string, string, ...Handler) IRoutes
	Any(string, ...Handler) IRoutes
	GET(string, ...Handler) IRoutes
	POST(string, ...Handler) IRoutes
	DELETE(string, ...Handler) IRoutes
	PATCH(string, ...Handler) IRoutes
	PUT(string, ...Handler) IRoutes
	OPTIONS(string, ...Handler) IRoutes
	HEAD(string, ...Handler) IRoutes
	Match([]string, string, ...Handler) IRoutes

	StaticFile(string, string) IRoutes
	StaticFileFS(string, string, http.FileSystem) IRoutes
	Static(string, string) IRoutes
	StaticFS(string, http.FileSystem) IRoutes
}

type Handler func(ctx *Context)
type Handlers []Handler

var _ IRoutes = &RouterGroup{}

func (group *RouterGroup) Use(handlers ...Handler) IRoutes {
	group.handlers = append(group.handlers, handlers...)
	return group.returnObj()
}

func (group *RouterGroup) Handle(httpMethod string, relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.Handle(httpMethod, relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}

func (group *RouterGroup) Any(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.Any(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}

func (group *RouterGroup) GET(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.GET(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) POST(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.POST(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) DELETE(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.DELETE(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}

func (group *RouterGroup) PATCH(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.PATCH(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) PUT(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.PUT(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) OPTIONS(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.OPTIONS(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) HEAD(relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.HEAD(relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) Match(methods []string, relativePath string, handlers ...Handler) IRoutes {
	group.ginRouterGroup.Match(methods, relativePath, group.toGinHandler(handlers...))
	return group.returnObj()
}
func (group *RouterGroup) StaticFile(relativePath string, filepath string) IRoutes {
	group.ginRouterGroup.StaticFile(relativePath, filepath)
	return group.returnObj()
}
func (group *RouterGroup) StaticFileFS(relativePath string, filepath string, fs http.FileSystem) IRoutes {
	group.ginRouterGroup.StaticFileFS(relativePath, filepath, fs)
	return group.returnObj()
}
func (group *RouterGroup) Static(relativePath string, root string) IRoutes {
	group.ginRouterGroup.Static(relativePath, root)
	return group.returnObj()
}
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	group.ginRouterGroup.StaticFS(relativePath, fs)
	return group.returnObj()
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.app
	}
	return group
}

func (group *RouterGroup) Group(relativePath string, handlers ...Handler) *RouterGroup {
	ginRouterGroup := group.ginRouterGroup.Group(relativePath)
	return &RouterGroup{
		ginRouterGroup: ginRouterGroup,
		handlers:       append(group.handlers, handlers...),
		app:            group.app,
		root:           false,
	}
}

func (group *RouterGroup) toGinHandler(handlers ...Handler) gin.HandlerFunc {
	return func(context *gin.Context) {
		ctx := &Context{
			Ctx:      context,
			handlers: append(group.handlers, handlers...),
			index:    -1,
		}
		ctx.Next()
	}
}
