package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadDBConfig() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	return os.Getenv("POSTGRES_URL"), nil
}
