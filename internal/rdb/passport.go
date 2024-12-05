package rdb

import (
	"context"
	"encoding/json"
	"go-web/internal/model"
	"go-web/utils/ecode"
	"time"
)

type Passport struct {
	*Rdb
}

func NewPassport(rdb *Rdb) *Passport {
	return &Passport{rdb}
}

func (t *Passport) GetPassport(token string) (m model.PassportDto, err error) {
	result, err := t.db.Get(context.Background(), "passport:"+token).Result()
	if err != nil {
		return
	}
	if result == "" {
		err = ecode.TokenExpired
		return
	}
	err = json.Unmarshal([]byte(result), &m)
	return
}

func (t *Passport) Create(token string, dto model.PassportDto) (err error) {
	data, err := json.Marshal(&dto)
	if err != nil {
		return
	}
	return t.db.SetNX(context.Background(), "passport:"+token, string(data), time.Hour*24*30).Err()
}

func (t *Passport) Delete(token string) error {
	script := `return redis.call("del", KEYS[1])`
	return t.db.Eval(context.Background(), script, []string{"passport:" + token}, nil).Err()
}
