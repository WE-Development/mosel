package main

import (
	"github.com/bluedevel/mosel/server"
)

func main() {
	server := server.MoselServer{}
	server.Run()
}