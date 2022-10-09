package main

import (
	"log"
	"turnstile/configs"
	app "turnstile/internal"
)

func main() {
	if err := configs.InitConfig(); err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run()
}
