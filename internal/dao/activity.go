package dao

import (
	"errors"
	"fmt"
	"go-web/internal/model"
	"go-web/utils/ecode"
	"gorm.io/gorm"
)

type Activity struct {
	*Dao
}

func NewActivity(dao *Dao) *Activity {
	return &Activity{dao}
}

func (t *Activity) Create(activity model.Activity) (m model.Activity, err error) {
	m = model.Activity{
		UUID:   activity.UUID,
		Name:   activity.Name,
		Url:    activity.Url,
		Type:   activity.Type,
		Status: activity.Status,
		Cost:   activity.Cost,
		Info:   activity.Info,
	}
	err = t.db.Model(&model.Activity{}).Create(&m).Error
	return
}

func (t *Activity) Delete(uuid int64) (err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	err = tx.Model(&model.Activity{}).Where("uuid", uuid).Delete(nil).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Model(&model.Prize{}).Where("activity_uuid", uuid).Delete(nil).Error
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

func (t *Activity) Update(activity model.Activity) (err error) {
	err = t.db.Model(&model.Activity{}).Where("uuid", activity.UUID).Select("name", "url", "type", "status", "cost", "info").Updates(&activity).Error
	return
}

func (t *Activity) Find(uuid int64) (m1 model.Activity, m2 []model.Prize, err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	err = tx.Model(&model.Activity{}).Where("uuid", uuid).Where("status", 1).First(&m1).Error
	if err != nil {
		if errors.As(err, &gorm.ErrRecordNotFound) {
			err = ecode.ActivityIsOver
		}
		tx.Rollback()
		return
	}
	err = tx.Model(&model.Prize{}).Where("activity_uuid", uuid).Find(&m2).Error
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

func (t *Activity) FindAll() (m []model.Activity, err error) {
	err = t.db.Model(&model.Activity{}).Find(&m).Error
	return
}

func (t *Activity) List(uuid int64) (prize []model.Prize, err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return
	}
	err = tx.Model(&model.Activity{}).Where("uuid", uuid).Update("status", 1).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Model(&model.Prize{}).Where("activity_uuid", uuid).Select("uuid", "stock", "score").Find(&prize).Error
	if err != nil || len(prize) < 1 {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
	}
	return
}

func (t *Activity) UnList(uuid int64) (err error) {
	err = t.db.Model(&model.Activity{}).Where("uuid", uuid).Update("status", 0).Error
	return
}

func (t *Activity) CreatePrize(prize ...model.Prize) (m []*model.Prize, err error) {
	m = make([]*model.Prize, len(prize))
	for k, v := range prize {
		m[k] = &model.Prize{
			ActivityUUID: v.ActivityUUID,
			UUID:         v.UUID,
			Name:         v.Name,
			Type:         v.Type,
			Value:        v.Value,
			InitialStock: v.InitialStock,
			Stock:        v.Stock,
			Score:        v.Score,
		}
	}
	err = t.db.Model(&model.Prize{}).Create(m).Error
	return
}

func (t *Activity) DeletePrize(activityUUID int64, uuid ...int64) (err error) {
	err = t.db.Model(&model.Prize{}).
		Where("activity_uuid", activityUUID).
		Where("uuid in (?)", uuid).Delete(nil).Error
	return
}

func (t *Activity) UpdatePrize(prize model.Prize) (err error) {
	err = t.db.Model(&model.Prize{}).
		Where("uuid", prize.UUID).
		Select("name", "type", "value", "initial_stock", "stock", "score").
		Updates(&prize).Error
	return
}

func (t *Activity) UpdatePrizeStock(data map[string]string) (err error) {
	var uuids []string
	cases := "CASE uuid"
	for uuid, stock := range data {
		uuids = append(uuids, uuid)
		cases += fmt.Sprintf(" WHEN '%s' THEN %s", uuid, stock)
	}
	cases += " END"

	sql := fmt.Sprintf("UPDATE prize SET stock = %s WHERE uuid IN (?)", cases)

	return t.db.Exec(sql, uuids).Error
}

func (t *Activity) CreateRecord(record []model.ActivityRecord) (err error) {
	err = t.db.Model(&model.ActivityRecord{}).Create(&record).Error
	return
}

func (t *Activity) FindRecord(uid int64, page int) (m []model.ActivityRecord, p model.Pagination, err error) {
	p = model.Pagination{
		Current:  page,
		PageSize: 20,
		Total:    0,
	}
	tx := t.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	err = tx.Model(&model.ActivityRecord{}).Where("uid", uid).Count(&p.Total).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Model(&model.ActivityRecord{}).Scopes(p.Sql()).Where("uid", uid).Find(&m).Error
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
