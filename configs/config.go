package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	errLoadingFile = "Error loading .env file"
	errToParseEnv  = "Error to parse env '%s': %s"
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
		log.Println("FATAAAAAAAAAAAL")
		return nil, errors.New(errLoadingFile)
	}

	fmt.Println(os.Getenv("API_PORT"), os.Getenv("MONGODB_TIMEOUT_IN_SECONDS"))

	mongoDbTimeout, err := parseEnvToInt("MONGODB_TIMEOUT_IN_SECONDS")
	if err != nil {
		return nil, err
	}

	expTime, err := parseEnvToInt("JWT_EXP_TIME")
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

func parseEnvToInt(key string) (int, error) {
	var value int
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return value, fmt.Errorf(errToParseEnv, key, err.Error())
	}
	return value, nil
}
