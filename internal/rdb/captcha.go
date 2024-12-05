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
	captcha    *base64Captcha.Captcha
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

func (t *redisStore) Set(id string, value string) error {
	return t.db.SetNX(context.Background(), "captcha:"+id, strings.ToLower(value), time.Minute*15).Err()
}

func (t *redisStore) Get(id string, clear bool) string {
	value, err := t.db.GetDel(context.Background(), "captcha:"+id).Result()
	if err != nil {
		return ""
	}
	return value
}

func (t *redisStore) Verify(id, answer string, clear bool) bool {
	if answer == t.Get(id, clear) {
		t.db.Del(context.Background(), id)
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
	c.captcha = base64Captcha.NewCaptcha(c.driver, c.redisStore)
	return c
}

func (t *Captcha) Generate() (id string, b64s string, answer string, err error) {
	return t.captcha.Generate()
}

func (t *Captcha) Verify(id string, answer string) (match bool) {
	return t.captcha.Verify(id, answer, true)
}
