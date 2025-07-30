package main

import (
	"os"

	"github.com/SButnyakov/luna/audio-processing/internal/app"
	"github.com/joho/godotenv"
)

const envFileName = ".env"

func main() {
	if _, err := os.Stat(envFileName); !os.IsNotExist(err) {
		godotenv.Load(envFileName)
	}

	app.Run()
}
