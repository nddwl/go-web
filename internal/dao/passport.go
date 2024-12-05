package dao

import "go-web/internal/model"

type Passport struct {
	*Dao
}

func NewPassport(dao *Dao) *Passport {
	return &Passport{dao}
}

func (t *Passport) Create(passport model.Passport) (m model.Passport, err error) {
	m = model.Passport{
		Uid:      passport.Uid,
		Token:    passport.Token,
		Ip:       passport.Ip,
		DeviceId: passport.DeviceId,
		Ua:       passport.Ua,
	}
	err = t.db.Model(&model.Passport{}).Create(&m).Error
	return
}

func (t *Passport) Delete(uid int64, deviceId string) (err error) {
	err = t.db.Model(&model.Passport{}).Where("uid", uid).Where("device_id", deviceId).Delete(nil).Error
	return
}

func (t *Passport) DeleteAll(uid int64) (err error) {
	err = t.db.Model(&model.Passport{}).Where("uid", uid).Delete(nil).Error
	return
}

func (t *Passport) Find(uid int64) (m []model.Passport, err error) {
	err = t.db.Model(&model.Passport{}).Where("uid", uid).Omit("token").Find(&m).Error
	return
}
