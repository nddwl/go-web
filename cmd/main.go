package main

import (
	"fmt"
	"go-web/internal/server"
)

func main() {
	err := server.New().Run()
	if err != nil {
		fmt.Println(err)
	}
}
