package rmq

type User struct {
	*Rmq
}

func NewUser(rmq *Rmq) *User {
	return &User{rmq}
}
