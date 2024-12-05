package rdb

import (
	"context"
	"github.com/mojocn/base64Captcha"
	"go-web/utils"
	"image/color"
	"strings"
	"time"
)

type Captcha struct {
	*Rdb
	driver     *driver
	redisStore *redisStore
}

type driver struct {
	*base64Captcha.DriverString
}

func (d *driver) GenerateIdQuestionAnswer() (id, content, answer string) {
	id = utils.GenerateToken()
	content = base64Captcha.RandText(d.Length, d.Source)
	return id, content, content
}

type redisStore struct {
	*Rdb
}

func (t *redisStore) Set(id string, name string, value string, expr time.Duration) error {
	return t.db.SetNX(context.Background(), "captcha_"+name+":"+id, strings.ToLower(value), expr).Err()
}

func (t *redisStore) Get(id string, name string) string {
	value, err := t.db.GetDel(context.Background(), "captcha_"+name+":"+id).Result()
	if err != nil {
		return ""
	}
	return value
}

func (t *redisStore) Verify(id, name, answer string) bool {
	if answer == t.Get(id, name) {
		return true
	}
	return false
}

func NewCaptcha(rdb *Rdb) *Captcha {
	c := &Captcha{
		Rdb: rdb,
		driver: &driver{
			DriverString: base64Captcha.NewDriverString(
				80,  // 高度
				240, // 宽度
				10,  // 噪声点数量
				1,   // 显示干扰线
				4,   // 生成4位验证码
				"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", // 字符集
				&color.RGBA{R: 255, G: 255, B: 255, A: 255},                      // 背景色
				nil, // 字体存储
				nil, // 字体文件
			),
		},
		redisStore: &redisStore{rdb},
	}
	return c
}

func (t *Captcha) Generate(name string, expr time.Duration) (id string, b64s string, answer string, err error) {
	id, content, answer := t.driver.GenerateIdQuestionAnswer()
	item, err := t.driver.DrawCaptcha(content)
	if err != nil {
		return "", "", "", err
	}
	err = t.redisStore.Set(id, name, answer, expr)
	if err != nil {
		return "", "", "", err
	}
	b64s = item.EncodeB64string()
	return
}

func (t *Captcha) Verify(id string, name string, answer string) (match bool) {
	return t.redisStore.Verify(id, name, answer)
}
