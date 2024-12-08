package service

import (
	"fmt"
	"go-web/internal/model"
	"go-web/utils"
	"go-web/utils/ecode"
	"go-web/utils/validate"
	"time"
)

type Activity struct {
	*Service
}

func NewActivity(service *Service) *Activity {
	return &Activity{service}
}

func (t *Activity) Create(activity model.Activity) (m model.Activity, err error) {
	activity.UUID = utils.GenerateUid()
	err = validate.Struct(activity)
	if err != nil {
		err = ecode.FormatError
		return
	}
	m, err = t.Dao.Activity.Create(activity)
	return
}

func (t *Activity) Delete(uuid int64) (err error) {
	if !validate.Uid(uuid) {
		err = ecode.FormatError
		return
	}
	list, err := t.Rdb.Activity.IsList(uuid)
	if err != nil {
		return
	}
	if list {
		err = ecode.ActivityIsList
		return
	}
	err = t.Dao.Activity.Delete(uuid)
	return
}

func (t *Activity) Update(activity model.Activity) (err error) {
	err = validate.Struct(activity)
	if err != nil {
		err = ecode.FormatError
		return
	}
	list, err := t.Rdb.Activity.IsList(activity.UUID)
	if err != nil {
		return
	}
	if list {
		err = ecode.ActivityIsList
		return
	}
	err = t.Dao.Activity.Update(activity)
	return
}

func (t *Activity) Find(uuid int64) (m1 model.Activity, m2 []model.Prize, err error) {
	if !validate.Uid(uuid) {
		err = ecode.FormatError
		return
	}
	m1, m2, err = t.Dao.Activity.Find(uuid)
	return
}

func (t *Activity) FindAll() (m []model.Activity, err error) {
	m, err = t.Dao.Activity.FindAll()
	return
}

func (t *Activity) List(uuid int64) (err error) {
	if !validate.Uid(uuid) {
		err = ecode.FormatError
		return
	}
	prize, err := t.Dao.Activity.List(uuid)
	if err != nil {
		return
	}
	return t.Rdb.Activity.List(uuid, prize)
}

func (t *Activity) UnList(uuid int64) (err error) {
	if !validate.Uid(uuid) {
		err = ecode.FormatError
		return
	}
	err = t.Dao.Activity.UnList(uuid)
	if err != nil {
		return
	}
	err = t.Rdb.Activity.UnList(uuid)
	return
}

func (t *Activity) CreatePrize(prize ...model.Prize) (m []*model.Prize, err error) {
	if len(prize) < 1 {
		err = ecode.FormatError
		return
	}
	var activityUUID int64
	for i := 0; i < len(prize); i++ {
		prize[i].UUID = utils.GenerateUid()
		err = validate.Struct(prize[i])
		if err != nil {
			err = ecode.FormatError
			return
		}
		if prize[i].ActivityUUID != activityUUID {
			activityUUID = prize[i].ActivityUUID
			list, err1 := t.Rdb.Activity.IsList(activityUUID)
			if err1 != nil {
				err = err1
				return
			}
			if list {
				err = ecode.ActivityIsList
				return
			}
		}

	}
	m, err = t.Dao.Activity.CreatePrize(prize...)
	return
}

func (t *Activity) DeletePrize(activityUUID int64, uuid ...int64) (err error) {
	if !validate.Uid(activityUUID) || len(uuid) < 1 {
		err = ecode.FormatError
		return
	}
	list, err := t.Rdb.Activity.IsList(activityUUID)
	if err != nil {
		return
	}
	if list {
		err = ecode.ActivityIsList
		return
	}
	for i := 0; i < len(uuid); i++ {
		if !validate.Uid(uuid[i]) {
			err = ecode.FormatError
			return
		}
	}
	err = t.Dao.Activity.DeletePrize(activityUUID, uuid...)
	return
}

func (t *Activity) UpdatePrize(prize model.Prize) (err error) {
	err = validate.Struct(prize)
	if err != nil {
		err = ecode.FormatError
		return
	}
	list, err := t.Rdb.Activity.IsList(prize.ActivityUUID)
	if err != nil {
		return
	}
	if list {
		err = ecode.ActivityIsList
		return
	}
	err = t.Dao.Activity.UpdatePrize(prize)
	return
}

func (t *Activity) UpdatePrizeStock(activityUUID int64) (data map[string]string, err error) {
	data, err = t.Rdb.Activity.GetPrizeStock(activityUUID)
	if len(data) < 1 {
		return
	}
	err = t.Dao.Activity.UpdatePrizeStock(data)
	return
}

func (t *Activity) CreateRecord() (err error) {
	return t.Rmq.Activity.ConsumeRecord(t.Dao.Activity.CreateRecord)
}

func (t *Activity) FindRecord(uid int64, page int) (m []model.ActivityRecord, p model.Pagination, err error) {
	m, p, err = t.Dao.Activity.FindRecord(uid, page)
	return
}

func (t *Activity) Lottery(activityUUID int64, uid int64) (m model.ActivityRecord, err error) {
	if !validate.Uid(activityUUID) {
		err = ecode.FormatError
		return
	}
	prizeUUID, err := t.Rdb.Activity.Lottery(activityUUID)
	if err != nil {
		return
	}
	m = model.ActivityRecord{
		Model: model.Model{
			CreatedAt: time.Now(),
		},
		Uid:          uid,
		ActivityUUID: activityUUID,
		PrizeUUID:    prizeUUID,
		Remark:       "",
	}
	go func(activity *Activity) {
		err1 := activity.Rmq.Activity.PublishRecord(m)
		if err1 != nil {
			fmt.Println(err1)
		}
	}(t)
	return
}
