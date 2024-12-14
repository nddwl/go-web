package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Model
	Uid    int64  `json:"uid"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Exp    int    `json:"exp"`
	Coin   int    `json:"coin"`
	Status uint8  `json:"status"`
	Role   uint8  `json:"role"`
}

type Model struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Password struct {
	Model
	Uid     int64
	PwdHash string
}

type UserCreate struct {
	Uid      int64  `json:"-"`
	Name     string `json:"name" validate:"name"`
	Password string `json:"password" validate:"password"`
	Email    string `json:"email" validate:"email"`
	Code     string `json:"code" validate:"code"`
	Session  string `json:"session" validate:"token"`
}

type UserUpdate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UserLogin struct {
	Type     string `json:"-"`
	Name     string `json:"name"`
	Session  string `json:"session" validate:"token"`
	Password string `json:"password" validate:"password"`
	Code     string `json:"code" validate:"code"`
}

type UserSign struct {
	Model
	Uid    int64  `json:"uid"`
	Status uint8  `json:"status"`
	Reward string `json:"reward"`
}

type Code struct {
	Captcha string `json:"captcha"`
	Session string `json:"session"`
}
