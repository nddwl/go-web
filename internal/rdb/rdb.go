package rdb

import (
	"github.com/redis/go-redis/v9"
	"go-web/utils/config"
	"time"
)

type Rdb struct {
	db *redis.Client
	Group
}

type Group struct {
	Captcha  *Captcha
	User     *User
	Activity *Activity
	Passport *Passport
}

func New() *Rdb {
	client := redis.NewClient(&redis.Options{
		// Redis 服务器的地址，通常是 localhost 或其他指定的主机和端口
		Addr: config.Rdb.Addr, // Redis 服务地址，格式为 "host:port"
		// 如果 Redis 服务器启用了密码保护，设置密码
		Password: config.Rdb.Password, // Redis 密码，空值表示没有密码
		// 默认使用数据库索引 0
		DB: 0, // 默认数据库索引，默认为 0
		// 连接池配置
		PoolSize: 100, // 连接池大小，设置为 10，表示最多可以有 10 个并发连接
		// 设置连接池中的最大空闲连接数
		MaxIdleConns: 10, // 最多空闲 5 个连接
		// 设置连接池的超时时间
		PoolTimeout: 5 * time.Second, // 获取连接的最大等待时间
		// 设置每次连接的超时时间
		DialTimeout: 10 * time.Second, // 连接超时
		// 设置读写超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时
	})
	rdb := &Rdb{db: client}
	rdb.initGroup()
	return rdb
}

func (t *Rdb) initGroup() {
	t.Group = Group{
		Captcha:  NewCaptcha(t),
		User:     NewUser(t),
		Activity: NewActivity(t),
		Passport: NewPassport(t),
	}
}
