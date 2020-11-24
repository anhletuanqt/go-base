package main

import (
	"base/app"
	"base/config"
	"base/database"
	"base/server"
	"fmt"
)

func main() {
	conf := config.New()
	database.Connect(conf)

	server := server.Setup()

	app.InitRoute(server)

	port := conf.Server.Port
	fmt.Println("Server is running at port:", port)
	server.Listen(":" + port)
}
