package main

import (
	"fmt"
	"go-web/internal/server"
)

func main() {
	server := server.New()
	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
