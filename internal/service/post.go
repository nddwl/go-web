package service

import (
	"go-web/internal/model"
	"go-web/utils"
	"go-web/utils/ecode"
	"go-web/utils/validate"
)

type Post struct {
	*Service
}

func NewPost(service *Service) *Post {
	return &Post{service}
}

func (t *Post) Create(post model.PostCreate) (err error) {
	post.UUID = utils.GenerateUid()
	err = validate.Struct(&post)
	if err != nil {
		err = ecode.FormatError
		return
	}
	tag := make([]*model.PostTag, len(post.Tag))
	for i := 0; i < len(post.Tag); i++ {
		tag[i].UUID = utils.GenerateUid()
		tag[i].Name = post.Tag[i]
	}
	_, _, err = t.Dao.Post.Create(post.Post, tag)
	return
}

func (t *Post) Delete(uid int64, uuid int64) (err error) {
	if !validate.Uid(uuid) {
		err = ecode.FormatError
		return
	}
	return t.Dao.Post.Delete(uid, uuid)
}

func (t *Post) Update(post model.Post) (err error) {
	err = validate.Struct(&post)
	if err != nil {
		err = ecode.FormatError
		return
	}
	return t.Dao.Post.Update(post)
}

func (t *Post) Find(uuid int64) (m model.Post, tag []string, err error) {
	if !validate.Uid(uuid) {
		err = ecode.FormatError
		return
	}
	return t.Dao.Post.Find(uuid)
}

func (t *Post) Finds(find model.PostFind) (m []model.PostCover, p model.Pagination, err error) {
	err = validate.Struct(&find)
	if err != nil {
		err = ecode.FormatError
		return
	}
	return t.Dao.Post.Finds(find)
}

func (t *Post) FindByUid(uid int64, page int) (m []model.PostCover, p model.Pagination, err error) {
	if !validate.Uid(uid) {
		err = ecode.FormatError
		return
	}
	return t.Dao.Post.FindByUid(uid, page)
}
