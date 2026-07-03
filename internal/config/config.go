package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
	POSTGRES_CONN string
	REDIS_URL string
}

func LoadConfig()(*Config, error){
	err:= godotenv.Load()
	if err != nil {
		return nil , err
	}

	return &Config{
       PORT: os.Getenv("PORT"),
	   POSTGRES_CONN: os.Getenv("PORT"),
	   REDIS_URL: os.Getenv("REDIS_URL"),
	}, nil
}