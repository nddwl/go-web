package main

import (
	"fmt"
	"go-web/internal/server/http"
)

func main() {
	server := http.New()
	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
