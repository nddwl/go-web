package service

import (
	"errors"
	"go-web/internal/model"
	"go-web/utils"
	"go-web/utils/ecode"
	"go-web/utils/lottery"
	"go-web/utils/validate"
	"strconv"
	"time"
)

type User struct {
	*Service
}

func NewUser(service *Service) *User {
	return &User{service}
}

func (t *User) Create(user model.UserCreate) (m model.User, err error) {
	err = validate.Struct(&user)
	if err != nil {
		err = ecode.FormatError
		return
	}
	user.Uid = utils.GenerateUid()
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return
	}
	m, err = t.Dao.User.Create(user)
	return
}

func (t *User) Delete(uid int64) (m model.User, err error) {
	err = t.Dao.User.Delete(uid)
	return
}

func (t *User) Update(uid int64, key string, value string) (err error) {
	switch key {
	case "name":
		if !validate.Name(value) {
			err = ecode.FormatError
			return
		}
	case "avatar":
		if !validate.Avatar(value) {
			err = ecode.FormatError
			return
		}
	case "email":
		if !validate.Email(value) {
			err = ecode.FormatError
			return
		}
	case "phone":
		if !validate.Phone(value) {
			err = ecode.FormatError
			return
		}
	default:
		err = ecode.FormatError
		return
	}
	err = t.Dao.User.Update(uid, key, value)
	return
}

func (t *User) Find(name string) (m model.User, err error) {
	if !validate.Name(name) {
		err = ecode.FormatError
		return
	}
	m, err = t.Dao.User.Find(name)
	return
}

func (t *User) Login(dto model.PasswordDto) (uid int64, err error) {
	err = validate.Struct(&dto)
	if err != nil {
		err = ecode.FormatError
		return
	}
	switch dto.Key {
	case "name":
		if !validate.Name(dto.Value) {
			err = ecode.FormatError
			return
		}
	case "email":
		if !validate.Email(dto.Value) {
			err = ecode.FormatError
			return
		}
	default:
		err = ecode.FormatError
		return
	}
	if !t.Service.Rdb.Captcha.Verify(dto.Session, dto.Code) {
		err = ecode.BadRequest
		return
	}
	password, err := t.Dao.User.Login(dto.Key, dto.Value)
	if err != nil {
		err = ecode.UserLoginFailed
		return
	}
	if utils.ComparePassword(password.PwdHash, dto.Password) {
		uid = password.Uid
		return
	} else {
		err = ecode.UserLoginFailed
		return
	}
}

func (t *User) IsNameExists(name string) (status bool, err error) {
	if !validate.Name(name) {
		err = ecode.FormatError
		return
	}
	return t.Dao.User.IsNameExist(name)
}

func (t *User) Sign(uid int64) (m model.UserSign, err error) {
	ok, err := t.Rdb.User.Sign(uid)
	if err != nil {
		return
	}
	if !ok {
		err = ecode.UserSigned
		return
	}
	coin, reward := lottery.Draw()
	sign := model.UserSign{
		Uid:    uid,
		Status: 1,
		Reward: reward,
	}
	m, err = t.Dao.User.Sign(coin, sign)
	if err != nil {
		err = t.Service.Rdb.User.DelSign(uid)
		return
	}
	return
}

func (t *User) GetUser(uid int64) (m model.User, err error) {
	user, err := t.Rdb.User.Get(uid)
	if err != nil {
		if errors.As(err, &ecode.UserNotExist) {
			m, err = t.Dao.User.GetUser(uid)
			if err != nil {
				return
			}
			err = t.Rdb.User.Create(m)
		}
		return
	}
	m = model.User{
		Uid:    uid,
		Name:   user["name"],
		Avatar: user["avatar"],
		Email:  user["email"],
		Phone:  user["phone"],
	}
	m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", user["create_at"])
	m.Exp, _ = strconv.Atoi(user["exp"])
	m.Coin, _ = strconv.Atoi(user["coin"])
	status, _ := strconv.ParseUint(user["status"], 10, 8)
	m.Status = uint8(status)
	role, _ := strconv.ParseUint(user["role"], 10, 8)
	m.Role = uint8(role)
	return
}
