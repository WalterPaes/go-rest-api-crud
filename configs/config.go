package configs

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	errLoadingFile = "Error loading .env file"
)

const appDevEnv = "development"

type Configs struct {
	ApiPort           string
	LogOutput         string
	LogLevel          string
	MongoDBUri        string
	MongoDBDatabase   string
	MongoDBCollection string
	MongoDBTimeout    int
	JwtSecret         string
	JwtExpTime        int
}

func Load(filenames ...string) (*Configs, error) {
	err := godotenv.Load(filenames...)
	if err != nil && os.Getenv("APP_ENV") == appDevEnv {
		return nil, errors.New(errLoadingFile)
	}

	mongoDbTimeout, err := strconv.Atoi(os.Getenv("MONGODB_TIMEOUT_IN_SECONDS"))
	if err != nil {
		return nil, err
	}

	expTime, err := strconv.Atoi(os.Getenv("JWT_EXP_TIME"))
	if err != nil {
		return nil, err
	}

	return &Configs{
		ApiPort:           os.Getenv("API_PORT"),
		LogOutput:         os.Getenv("LOG_OUTPUT"),
		LogLevel:          os.Getenv("LOG_LEVEL"),
		MongoDBUri:        os.Getenv("MONGODB_URI"),
		MongoDBDatabase:   os.Getenv("MONGODB_DATABASE"),
		MongoDBCollection: os.Getenv("MONGODB_COLLECTION"),
		MongoDBTimeout:    mongoDbTimeout,
		JwtSecret:         os.Getenv("JWT_SECRET"),
		JwtExpTime:        expTime,
	}, nil
}
