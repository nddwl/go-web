package model

import "gorm.io/gorm"

type Activity struct {
	Model
	UUID   int64  `json:"uuid" validate:"uid"`
	Name   string `json:"name" validate:"required,noHTML,safaInput,min=3,max=50"`
	Url    string `json:"url" validate:"omitempty,url"`
	Type   uint8  `json:"type"`
	Status uint8  `json:"status"`
	Cost   uint   `json:"cost"`
	Info   string `json:"info" validate:"omitempty,noHTML,max=255"`
}

type Prize struct {
	Model
	ActivityUUID int64  `json:"activity_uuid" validate:"uid"`
	UUID         int64  `json:"uuid" validate:"uid"`
	Name         string `json:"name" validate:"required,noHTML,safaInput,min=3,max=50"`
	Type         uint8  `json:"type"`
	Value        string `json:"value" validate:"required,noHTML,safaInput,min=3,max=50"`
	InitialStock uint   `json:"initial_stock" validate:"min=1"`
	Stock        uint   `json:"stock" validate:"ltefield=InitialStock"`
	Score        uint   `json:"score" validate:"gt=0"`
}

type ActivityRecord struct {
	Model
	Uid          int64  `json:"uid"`
	ActivityUUID int64  `json:"activity_uuid"`
	PrizeUUID    int64  `json:"prize_uuid"`
	Remark       string `json:"remark"`
}

type Pagination struct {
	Current  int   `json:"current"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func (t *Pagination) Sql() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Current <= 0 {
			t.Current = 1
		}

		switch {
		case t.PageSize > 100:
			t.PageSize = 100
		case t.PageSize <= 0:
			t.PageSize = 20
		}

		return db.Offset((t.Current - 1) * t.PageSize).Limit(t.PageSize)
	}
}

func (t *Pagination) Copy(p Pagination) {
	t.Current = p.Current
	t.PageSize = p.PageSize
	t.Total = p.Total
}
