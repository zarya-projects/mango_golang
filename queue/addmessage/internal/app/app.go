package app

import (
	"main/internal/config"
	sl "main/internal/logger"
	"main/internal/transport"
)

func Launch(cfg *config.Config) {
	logger := sl.MustLoad()
	transport.Hudlers(cfg, logger)
}
