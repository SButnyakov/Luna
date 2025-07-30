package main

import (
	"fmt"

	"github.com/SButnyakov/luna/audio-processing/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
