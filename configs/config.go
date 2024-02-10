package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var (
	errLoadingFile = "Error loading .env file"
)

const appDevEnv = "development"

type Configs struct {
	ApiPort   string
	LogOutput string
	LogLevel  string
}

func Load(filenames ...string) (*Configs, error) {
	err := godotenv.Load(filenames...)
	if err != nil && os.Getenv("APP_ENV") == appDevEnv {
		return nil, errors.New(errLoadingFile)
	}

	return &Configs{
		ApiPort:   os.Getenv("API_PORT"),
		LogOutput: os.Getenv("LOG_OUTPUT"),
		LogLevel:  os.Getenv("LOG_LEVEL"),
	}, nil
}
