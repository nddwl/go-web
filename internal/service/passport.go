package service

import (
	"fmt"
	"go-web/internal/model"
	"go-web/utils"
	"go-web/utils/ecode"
	"go-web/utils/validate"
)

type Passport struct {
	*Service
}

func NewPassport(service *Service) *Passport {
	return &Passport{service}
}

func (t *Passport) Create(passport model.Passport) (m model.Passport, err error) {
	passport.Token = utils.GenerateToken()
	passport.DeviceId = utils.GenerateToken()
	dto := model.PassportDto{
		Uid:      passport.Uid,
		DeviceId: passport.DeviceId,
	}
	err = t.Rdb.Passport.Create(passport.Token, dto)
	if err != nil {
		return
	}
	m, err = t.Dao.Passport.Create(passport)
	if err != nil {
		err1 := t.Rdb.Passport.Delete(passport.Token)
		if err1 != nil {
			fmt.Println(err1)
		}
		return
	}
	return
}

func (t *Passport) Delete(uid int64, deviceId string) (err error) {
	return t.Dao.Passport.Delete(uid, deviceId)
}

func (t *Passport) DeleteAll(uid int64) (err error) {
	return t.Dao.Passport.DeleteAll(uid)
}

func (t *Passport) Find(uid int64) (m []model.Passport, err error) {
	return t.Dao.Passport.Find(uid)
}

func (t *Passport) GetPassport(token string) (m model.PassportDto, err error) {
	if !validate.Token(token) {
		err = ecode.FormatError
		return
	}
	m, err = t.Rdb.Passport.GetPassport(token)
	if err != nil {
		return
	}
	return
}
