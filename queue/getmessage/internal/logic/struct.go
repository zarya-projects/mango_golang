package logic

import (
	"main/internal/database"
)

type Connection struct {
	Connect ConnectMethod
	Val     *Values
	DB      database.StructDatabase
}

type ConnectMethod interface {
	Launch() error
	openChannel() error
}

type Values struct {
	Phone string `json:"phone"`
	ID    string `json:"id"`
	Stage string `json:"stage"`
	City  int    `json:"city"`
}
