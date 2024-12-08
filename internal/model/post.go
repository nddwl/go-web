package model

type Post struct {
	Model
	UUID        int64  `json:"uuid" validate:"uid"`
	Title       string `json:"title" validate:"required,noHTML,safaInput,min=3,max=50"`
	UID         int64  `json:"uid" validate:"uid"`
	Category    uint8  `json:"category"`
	TopCategory uint8  `json:"top_category"`
	Summary     string `json:"summary" validate:"required,noHTML,max=255"`
	Content     string `json:"content" validate:"required"`
	Source      uint8  `json:"source" validate:"omitempty,oneof=0 1"`
}

type PostTag struct {
	Model
	UUID int64  `json:"uuid" validate:"uid"`
	Name string `json:"name" validate:"required,noHTML,safaInput,min=1,max=10"`
}

type PostAssociation struct {
	Model
	PostTagUUID int64 `json:"post_tag_uuid" validate:"uid"`
	PostUUID    int64 `json:"post_uuid" validate:"uid"`
}
