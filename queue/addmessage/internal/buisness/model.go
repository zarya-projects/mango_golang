package buisness

import (
	"log/slog"
	"main/internal/config"
	"main/internal/database"
	"net/http"
)

type AddMessage interface {
	Add(w http.ResponseWriter, r *http.Request)
	getData(id string) error
	addFromRabbitMQ() error
}

type Values struct {
	ID       string `json:"id"`
	Phone    string `json:"phone"`
	Stage    string `json:"stage"`
	City     int    `json:"city"`
	Source   string `json:"source"`
	Region   int    `json:"region"`
	Napr     int    `json:"napr"`
	CityName string `json:"cityname"`
}

type Message struct {
	Val     Values
	AddMess AddMessage
	Config  *config.Config
	DB      database.StructDatabase
	Logger  *slog.Logger
}
