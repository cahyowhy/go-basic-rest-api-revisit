package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DbConfig  DbConfig
	AppEnv    string
	JWTSECRET string
	PORT      string
}

type DbConfig struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
}

var config *Config
var onceConfig sync.Once

func GetConfig(envFile ...string) *Config {
	onceConfig.Do(func() {
		err := godotenv.Load(envFile...)

		if err != nil {
			log.Fatalf("Error env file : %s", err.Error())
		}

		config = &Config{
			AppEnv: os.Getenv("APP_ENV"),
			DbConfig: DbConfig{
				Username: os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASSWORD"),
				Name:     os.Getenv("DB"),
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("DB_PORT"),
			},
			JWTSECRET: os.Getenv("JWT_SECRET"),
			PORT:      os.Getenv("PORT"),
		}
	})

	return config
}
