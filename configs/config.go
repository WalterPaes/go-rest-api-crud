package configs

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	errLoadingFile = "Error loading .env file"
)

const appDevEnv = "development"

type Configs struct {
	ApiPort        string
	LogOutput      string
	LogLevel       string
	MongoDBUri     string
	MongoDBTimeout int
}

func Load(filenames ...string) (*Configs, error) {
	err := godotenv.Load(filenames...)
	if err != nil && os.Getenv("APP_ENV") == appDevEnv {
		return nil, errors.New(errLoadingFile)
	}

	mongoDbTimeout, err := strconv.Atoi(os.Getenv("MONGODB_TIMEOUT_IN_SECONDS"))
	if err != nil {
		log.Fatal(err)
	}

	return &Configs{
		ApiPort:        os.Getenv("API_PORT"),
		LogOutput:      os.Getenv("LOG_OUTPUT"),
		LogLevel:       os.Getenv("LOG_LEVEL"),
		MongoDBUri:     os.Getenv("MONGODB_URI"),
		MongoDBTimeout: mongoDbTimeout,
	}, nil
}
