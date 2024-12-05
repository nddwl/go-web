package service

import (
	"go-web/internal/dao"
	"go-web/internal/rdb"
	"go-web/internal/rmq"
)

type Service struct {
	Dao *dao.Dao
	Rdb *rdb.Rdb
	Rmq *rmq.Rmq
	Group
}

type Group struct {
	User     *User
	Passport *Passport
	Activity *Activity
}

func New() *Service {
	service := &Service{
		Dao: dao.New(),
		Rdb: rdb.New(),
		Rmq: rmq.New(),
	}
	service.initGroup()
	return service
}

func (t *Service) initGroup() {
	t.Group = Group{
		User:     NewUser(t),
		Passport: NewPassport(t),
		Activity: NewActivity(t),
	}
}
