package server

import (
	"go-web/internal/model"
	"go-web/utils/app"
	"go-web/utils/ecode"
)

type Post struct {
	*Server
	Group *app.RouterGroup
}

func NewPost(server *Server) *Post {
	post := &Post{
		Server: server,
		Group:  nil,
	}
	return post
}

func (t *Post) initGroup() {
	t.Group = t.Server.App.Group("/post")
	t.Group.POST("/create", t.Create)
	t.Group.POST("/delete", t.Delete)
	t.Group.POST("/update", t.Update)
	t.Group.GET("/find", t.Find)
	t.Group.GET("/finds", t.Finds)
	t.Group.GET("/findByUid", t.FindByUid)
}

func (t *Post) Create(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	create := model.PostCreate{}
	err := ctx.ShouldBindJSON(&create)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	create.UID = ctx.Passport.Uid
	err = t.Service.Post.Create(create)
	ctx.JSON(nil, err)
}

func (t *Post) Delete(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	postUUID := ctx.ParseRequestResource("post_uuid").Int()
	err := t.Service.Post.Delete(ctx.Passport.Uid, postUUID)
	ctx.JSON(nil, err)
}

func (t *Post) Update(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	post := model.Post{}
	err := ctx.ShouldBindJSON(&post)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	post.UID = ctx.Passport.Uid
	err = t.Service.Post.Update(post)
	ctx.JSON(nil, err)
}

func (t *Post) Find(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	postUUID := ctx.ParseRequestResource("post_uuid").Int()
	m, tag, err := t.Service.Post.Find(postUUID)
	ctx.JSON(&struct {
		Post model.Post `json:"post"`
		Tag  []string   `json:"tag"`
	}{m, tag}, err)
}

func (t *Post) Finds(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	find := model.PostFind{}
	err := ctx.ShouldBindJSON(&find)
	if err != nil {
		ctx.JSON(nil, ecode.FormatError)
		return
	}
	m, p, err := t.Service.Post.Finds(find)
	ctx.JSON(&struct {
		PostCover  []model.PostCover `json:"post_cover"`
		Pagination model.Pagination  `json:"pagination"`
	}{m, p}, err)
}

func (t *Post) FindByUid(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	uid := ctx.ParseRequestResource("uid").Int()
	page := int(ctx.ParseRequestResource("page").Int())
	m, p, err := t.Service.Post.FindByUid(uid, page)
	ctx.JSON(&struct {
		PostCover  []model.PostCover `json:"post_cover"`
		Pagination model.Pagination  `json:"pagination"`
	}{m, p}, err)
}
