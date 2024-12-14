package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-web/internal/model"
	"go-web/utils/ecode"
	"strconv"
	"time"
)

type User struct {
	*Rdb
}

func NewUser(rdb *Rdb) *User {
	return &User{rdb}
}

func (t *User) Sign(uid int64) (bool, error) {
	return t.db.SetNX(context.Background(), "sign:"+strconv.FormatInt(uid, 10), "1", time.Hour*24).Result()
}

func (t *User) DelSign(uid int64) error {
	return t.db.Del(context.Background(), "sign:"+strconv.FormatInt(uid, 10)).Err()
}

func (t *User) Create(user model.User) (err error) {
	values := map[string]interface{}{
		"create_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"name":      user.Name,
		"avatar":    user.Avatar,
		"email":     user.Email,
		"exp":       user.Exp,
		"coin":      user.Coin,
		"status":    user.Status,
		"role":      user.Role,
	}
	ctx := context.Background()

	err = t.db.Watch(ctx, func(tx *redis.Tx) error {
		cmder, err1 := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, "user:"+strconv.FormatInt(user.Uid, 10), values)
			pipe.Expire(ctx, "user:"+strconv.FormatInt(user.Uid, 10), 7*time.Hour*24)
			return nil
		})
		if err1 != nil {
			return err1
		}
		return cmder[0].Err()
	}, "user:"+strconv.FormatInt(user.Uid, 10))
	return
}

func (t *User) Get(uid int64) (user map[string]string, err error) {
	user, err = t.db.HGetAll(context.Background(), "user:"+strconv.FormatInt(uid, 10)).Result()
	if err != nil {
		return
	}
	if len(user) == 0 {
		err = ecode.UserNotExist
		return
	}
	return
}

func (t *User) Set(uid int64, values ...interface{}) error {
	ctx := context.Background()
	luaScript := `return redis.call('HSET', KEYS[1], unpack(ARGV))`
	_, err := t.db.Eval(ctx, luaScript, []string{"user:" + strconv.FormatInt(uid, 10)}, values...).Result()
	return err
}
