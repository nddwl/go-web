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
	Phone  string `json:"phone"`
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
	Avatar   string `json:"avatar" validate:"omitempty,url"`
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone" validate:"omitempty,phone"`
	Code     string `json:"code" validate:"code"`
	Session  string `json:"-"`
}

type UserUpdate struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PasswordDto struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Session  string `json:"-"`
	Password string `json:"password" validate:"password"`
	Code     string `json:"code" validate:"code"`
}

type UserSign struct {
	Model
	Uid    int64  `json:"uid"`
	Status uint8  `json:"status"`
	Reward string `json:"reward"`
}
