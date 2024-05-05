package main

import (
	"log"
	"logpuller/server"
)

func main() {
	app := server.NewApp()
	if err := app.Run("8080"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
