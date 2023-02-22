package main

import (
	"restaurants-service/server"
)

func init() {
	server.InitViper(".")
}

func main() {
	server.Start()
}
