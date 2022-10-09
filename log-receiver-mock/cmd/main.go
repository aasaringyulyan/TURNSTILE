package main

import (
	"log"
	"log-receiver-mock/configs"
	app "log-receiver-mock/internal"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run()
}
