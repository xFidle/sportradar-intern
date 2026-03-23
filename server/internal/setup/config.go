package setup

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/xFidle/sportradar-intern/server/internal/db"
)

type config struct {
	port           string
	fileserverAddr string
	database       db.Config
}

func loadConfig() config {
	return config{
		port:           os.Getenv("PORT"),
		fileserverAddr: os.Getenv("FILESERVER_ADDR"),
		database: db.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
	}
}
