package config

import (
	"encoding/json"
	"go-web/configs"
	"io"
	"os"
)

var (
	Amqp   configs.Amqp
	Db     configs.Db
	Rdb    configs.Rdb
	Server configs.Server
)

func init() {
	Load("amqp", &Amqp)
	Load("db", &Db)
	Load("rdb", &Rdb)
	Load("server", &Server)
}

func Load(name string, obj interface{}) {
	f, err := os.Open("./configs/" + name + ".json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &obj)
	if err != nil {
		panic(err)
	}
}

func IsLocal() bool {
	return os.Getenv("IS_IDEA") == "true"
}
