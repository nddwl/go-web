package dao

import (
	"go-web/internal/model"
	"go-web/utils/ecode"
	"gorm.io/gorm/clause"
)

type Post struct {
	*Dao
}

func NewPost(dao *Dao) *Post {
	return &Post{dao}
}

func (t *Post) Create(post model.Post, tag []*model.PostTag) (m1 model.Post, m2 []*model.PostTag, err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	m1 = post
	err = tx.Model(&model.Post{}).Create(&m1).Error
	if err != nil {
		tx.Rollback()
		return
	}
	m2 = tag
	err = tx.Model(&model.PostTag{}).Clauses(clause.OnConflict{DoNothing: true}).Create(&tag).Error
	if err != nil {
		tx.Rollback()
		return
	}
	association := make([]*model.PostAssociation, len(tag))
	for i := 0; i < len(tag); i++ {
		association[i].PostUUID = post.UUID
		association[i].PostTagUUID = tag[i].UUID
	}
	err = tx.Model(&model.PostAssociation{}).Create(&association).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

func (t *Post) Delete(uid int64, uuid int64) (err error) {
	tx := t.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var count int64
	err = tx.Model(&model.Post{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("uuid = ? AND uid = ?", uuid, uid).Count(&count).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if count != 1 {
		err = ecode.Forbidden
		tx.Rollback()
		return
	}

	err = tx.Model(&model.PostAssociation{}).Where("post_uuid = ?", uuid).Delete(nil).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&model.Post{}).Where("uuid = ? AND uid = ?", uuid, uid).Delete(nil).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (t *Post) Update(post model.Post) (err error) {
	return t.db.Model(&model.Post{}).Where("uid", post.UID).Where("uuid", post.UUID).Updates(&post).Error
}

func (t *Post) Find(uuid int64) (m model.Post, tag []string, err error) {
	err = t.db.Model(&model.Post{}).
		Joins("INNER JOIN post_association ON post_association.post_uuid = post.uuid").
		Joins("INNER JOIN post_tag ON post_tag.uuid = post_association.post_tag_uuid").
		Where("post.uuid = ?", uuid).
		Pluck("post_tag.name", &tag).
		First(&m).Error
	return
}

func (t *Post) Finds(find model.PostFind) (m []model.PostCover, p model.Pagination, err error) {
	p = model.Pagination{
		Current:  find.Page,
		PageSize: 20,
		Total:    0,
	}
	sql := t.db.Model(&model.Post{}).Scopes(p.Sql())
	if find.Title != "" {
		sql = sql.Where("MATCH(title) AGAINST(?)", find.Title)
	}
	sql = sql.Where("category", find.Category)
	err = sql.Scopes(p.Sql()).Find(&m).Error
	if err != nil {
		return
	}
	err = sql.Count(&p.Total).Error
	return
}

func (t *Post) FindByUid(uid int64, page int) (m []model.PostCover, p model.Pagination, err error) {
	p = model.Pagination{
		Current:  page,
		PageSize: 20,
		Total:    0,
	}
	err = t.db.Model(&model.Post{}).Where("uid", uid).Scopes(p.Sql()).Find(&m).Error
	err = t.db.Model(&model.Post{}).Where("uid", uid).Count(&p.Total).Error
	return
}
