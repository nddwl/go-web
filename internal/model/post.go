package model

type Post struct {
	Model
	UUID     int64  `json:"uuid" validate:"uid"`
	Title    string `json:"title" validate:"required,noHTML,safeInput,min=3,max=50"`
	UID      int64  `json:"uid" validate:"uid"`
	Cover    string `json:"cover"`
	Category uint8  `json:"category"`
	Summary  string `json:"summary" validate:"required,noHTML,max=255"`
	Content  string `json:"content" validate:"required"`
	Source   uint8  `json:"source" validate:"omitempty,oneof=0 1"`
}

type PostCover struct {
	UUID     int64  `json:"uuid" validate:"uid"`
	Title    string `json:"title" validate:"required,noHTML,safeInput,min=3,max=50"`
	UID      int64  `json:"uid" validate:"uid"`
	Cover    string `json:"cover"`
	Category uint8  `json:"category"`
}

type PostTag struct {
	Model
	UUID int64  `json:"uuid" validate:"uid"`
	Name string `json:"name" validate:"required,noHTML,safeInput,min=1,max=10"`
}

type PostAssociation struct {
	Model
	PostTagUUID int64 `json:"post_tag_uuid" validate:"uid"`
	PostUUID    int64 `json:"post_uuid" validate:"uid"`
}

type PostFind struct {
	Title    string `json:"title" validate:"omitempty,noHTML,safeInput,min=3,max=20"`
	Category uint8  `json:"category"`
	Page     int    `json:"page"`
}

type PostCreate struct {
	Post
	Tag []string `json:"tag" validate:"dive,required,noHTML,safeInput,min=1,max=10"`
}
