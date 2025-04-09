package main

import (
	"main/internal/app"
	"main/internal/config"
)

func main() {
	cfg := config.Get()
	app.Launch(cfg)
}
