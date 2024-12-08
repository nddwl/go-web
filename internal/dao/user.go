package dao

import (
	"errors"
	"go-web/internal/model"
	"go-web/utils/ecode"
	"gorm.io/gorm"
)

type User struct {
	*Dao
}

func NewUser(dao *Dao) *User {
	return &User{dao}
}

func (t *User) Create(user model.UserCreate) (m model.User, err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	m = model.User{
		Uid:    user.Uid,
		Name:   user.Name,
		Avatar: user.Avatar,
		Email:  user.Email,
		Phone:  user.Phone,
		Role:   1,
	}
	password := model.Password{
		Uid:     user.Uid,
		PwdHash: user.Password,
	}
	err = tx.Model(&model.User{}).Create(&m).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Model(&model.Password{}).Create(&password).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
	}
	return
}

func (t *User) Delete(uid int64) (err error) {
	err = t.db.Model(&model.User{}).Where("uid", uid).Delete(nil).Error
	return
}

func (t *User) Update(uid int64, key string, value string) (err error) {
	err = t.db.Model(&model.User{}).Where("uid", uid).Update(key, value).Error
	return
}

func (t *User) Find(name string) (m model.User, err error) {
	err = t.db.Model(&model.User{}).Omit("email", "phone").Where("name", name).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ecode.UserNotExist
	}
	return
}

func (t *User) Login(key string, value string) (m model.Password, err error) {
	err = t.db.Model(&model.User{}).
		Joins("INNER JOIN password ON user.uid = password.uid").
		Where("user."+key, value).
		Select("password.uid", "password.pwd_hash").
		First(&m).Error
	return
}

func (t *User) IsNameExist(name string) (exists bool, err error) {
	err = t.db.Model(&model.User{}).Where("name", name).Select("1").Limit(1).Find(&exists).Error
	return
}

func (t *User) Sign(coin int, sign model.UserSign) (m model.UserSign, err error) {
	m = model.UserSign{
		Uid:    sign.Uid,
		Status: sign.Status,
		Reward: sign.Reward,
	}
	tx := t.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	err = tx.Model(&model.UserSign{}).Create(&m).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Model(&model.User{}).Where("uid", sign.Uid).Update("coin", coin).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
	}
	return
}

func (t *User) GetUser(uid int64) (m model.User, err error) {
	err = t.db.Model(&model.User{}).Where("uid", uid).First(&m).Error
	return
}
