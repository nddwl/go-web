package http

import (
	"go-web/internal/model"
	"go-web/utils/app"
	"go-web/utils/ecode"
)

type Activity struct {
	*Server
	Group *app.RouterGroup
}

func NewActivity(server *Server) *Activity {
	activity := &Activity{
		Server: server,
		Group:  nil,
	}
	return activity
}

func (t *Activity) initGroup() {
	t.Group = t.Server.App.Group("/activity", t.User.GetUser)
	t.Group.POST("/create", t.Create)
	t.Group.POST("/delete", t.Delete)
	t.Group.POST("/update", t.Delete)
	t.Group.GET("/find", t.Find)
	t.Group.GET("/findAll", t.FindAll)
	t.Group.POST("/list", t.List)
	t.Group.POST("/unList", t.UnList)
	t.Group.POST("/createPrize", t.CreatePrize)
	t.Group.POST("/deletePrize", t.DeletePrize)
	t.Group.POST("/updatePrize", t.UpdatePrize)
	t.Group.POST("/updatePrizeStock", t.UpdatePrizeStock)
	t.Group.POST("/lottery", t.Lottery)
	t.Group.POST("/createRecord", t.CreateRecord)
	t.Group.GET("/findRecord", t.FindRecord)
}

func (t *Activity) Create(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	activity := model.Activity{}
	err := ctx.ShouldBindJSON(&activity)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	m, err := t.Service.Activity.Create(activity)
	ctx.JSON(&m, err)
}

func (t *Activity) Delete(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	result := ctx.ParseRequestResource("activity_uuid")
	err := t.Service.Activity.Delete(result.Int())
	ctx.JSON(nil, err)
}

func (t *Activity) Update(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	activity := model.Activity{}
	err := ctx.ShouldBindJSON(&activity)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	err = t.Service.Activity.Update(activity)
	ctx.JSON(nil, err)
}
func (t *Activity) Find(ctx *app.Context) {
	result := ctx.ParseRequestResource("activity_uuid")
	m1, m2, err := t.Service.Activity.Find(result.Int())
	ctx.JSON(&struct {
		Activity model.Activity
		Prize    []model.Prize
	}{m1, m2}, err)
}

func (t *Activity) FindAll(ctx *app.Context) {
	m, err := t.Service.Activity.FindAll()
	ctx.JSON(&m, err)
}

func (t *Activity) List(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	result := ctx.ParseRequestResource("activity_uuid")
	err := t.Service.Activity.List(result.Int())
	ctx.JSON(nil, err)
}

func (t *Activity) UnList(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	result := ctx.ParseRequestResource("activity_uuid")
	err := t.Service.Activity.UnList(result.Int())
	ctx.JSON(nil, err)
}

func (t *Activity) CreatePrize(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	var prize []model.Prize
	err := ctx.ShouldBindJSON(&prize)
	if err != nil {
		ctx.JSON(nil, err)
	}
	m, err := t.Service.Activity.CreatePrize(prize...)
	ctx.JSON(&m, err)
}
func (t *Activity) DeletePrize(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	result := ctx.ParseRequestResourceMany("activity_uuid", "prize_uuid")
	var name []int64
	for _, v := range result[1].Array() {
		name = append(name, v.Int())
	}
	err := t.Service.Activity.DeletePrize(result[0].Int(), name...)
	ctx.JSON(nil, err)
}
func (t *Activity) UpdatePrize(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	var prize []model.Prize
	err := ctx.ShouldBindJSON(&prize)
	if err != nil {
		ctx.JSON(nil, err)
		return
	}
	err = t.Service.Activity.UpdatePrize(prize...)
	ctx.JSON(nil, err)
}

func (t *Activity) UpdatePrizeStock(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	result := ctx.ParseRequestResource("activity_uuid")
	data, err := t.Service.Activity.UpdatePrizeStock(result.Int())
	ctx.JSON(&data, err)
}

func (t *Activity) Lottery(ctx *app.Context) {
	uuid := ctx.ParseRequestResource("activity_uuid")
	m, err := t.Service.Activity.Lottery(uuid.Int(), ctx.Passport.Uid)
	ctx.JSON(&m, err)
}

func (t *Activity) CreateRecord(ctx *app.Context) {
	if !ctx.IsAdmin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	err := t.Service.Activity.CreateRecord()
	ctx.JSON(nil, err)
}

func (t *Activity) FindRecord(ctx *app.Context) {
	if !ctx.IsLogin() {
		ctx.JSON(nil, ecode.Forbidden)
		return
	}
	page := ctx.ParseRequestResource("page").Int()
	m, p, err := t.Service.Activity.FindRecord(ctx.Passport.Uid, int(page))
	ctx.JSON(&struct {
		Record     []model.ActivityRecord
		Pagination model.Pagination
	}{m, p}, err)
}
