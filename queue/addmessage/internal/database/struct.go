package database

import (
	"database/sql"
	"main/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type StructDatabase struct {
	DB       DB
	Cfg      *config.Config
	BaseOpen *sql.DB
}

type ConnectDB interface {
	DataDB() (string, string)
	DBJoin() (*sql.DB, error)
	DBClose(dp *sql.DB)
}

type DB interface {
	Launch() error
}

// Получение из .env данных по коннекту к БД
func (s *StructDatabase) DataDB() (string, string) {
	return s.Cfg.DB_TYPE, s.Cfg.DB_LOGIN + ":" + s.Cfg.DB_PASS + "@/" + s.Cfg.DB_NAME
}

// открытие базы
func (s *StructDatabase) DBJoin() (*sql.DB, error) {
	dbName, dbData := s.DataDB()
	dbOpen, err := sql.Open(dbName, dbData)
	if err != nil {
		return nil, err
	}

	return dbOpen, nil
}

// Закрытие базы
func (s *StructDatabase) DBClose() {
	defer s.BaseOpen.Close()
}
