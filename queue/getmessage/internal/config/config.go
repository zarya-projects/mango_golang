package config

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

func Get() *Config {

	if err := godotenv.Load("/home/root666/go/project/mango/.env"); err != nil {
		panic(err)
	}

	return &Config{
		RABBIT_HOST:  os.Getenv("RABBIT_HOST"),
		RABBIT_LOGIN: os.Getenv("RABBIT_LOGIN"),
		RABBIT_NAME:  os.Getenv("RABBIT_NAME"),
		RABBIT_PORT:  os.Getenv("RABBIT_PORT"),
		RABBIT_PASS:  os.Getenv("RABBIT_PASSWORD"),
		DB_TYPE:      os.Getenv("DB_TYPE"),
		DB_NAME:      os.Getenv("DB_NAME"),
		DB_HOST:      os.Getenv("DB_HOST"),
		DB_LOGIN:     os.Getenv("DB_LOGIN"),
		DB_PASS:      os.Getenv("DB_PASS"),
		DB_PORT:      os.Getenv("DB_PORT"),
	}
}
