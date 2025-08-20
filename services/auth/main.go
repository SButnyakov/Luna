package main

import (
	"log"
	"os"

	"github.com/SButnyakov/luna/id/internal/app"
	"github.com/joho/godotenv"
)

const (
	envFile = ".env"
)

func main() {
	if _, err := os.Stat(envFile); err == nil {
		log.Println("loading envs")
		if err := godotenv.Load(envFile); err != nil {
			log.Fatal("failed to load envs", err)
		}
	}
	app.Run()
}
