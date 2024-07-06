package main

import (
	"fmt"
	"websocket/config"
	"websocket/internal/service"
	"websocket/pkg/MySQL"
	"websocket/route"
)

func main() {
	err := config.LoadConfig("./config.json")
	if err != nil {
		panic(err)
	}

	_, err = MySQL.InitMySQL()
	if err != nil {
		panic(err)
	}
	go service.Start()

	r := route.NewRouter()
	address := fmt.Sprintf("%v:%v", config.AppConfig.Address, config.AppConfig.Port)
	_ = r.Run(address)
}
