package dao

type Post struct {
	*Dao
}

func NewPost(dao *Dao) *Post {
	return &Post{dao}
}
