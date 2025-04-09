package app

import (
	"main/internal/config"
	"main/internal/logic"
	sl "main/internal/logs"
)

type App struct {
	Connect logic.Connection
	Cfg     *config.Config
}

func (a App) Launch(cfg *config.Config) {

	log := sl.MustLoad()
	if err := a.Connect.Launch(cfg, log); err != nil {
		log.Error("err with connect.")
	}
}
