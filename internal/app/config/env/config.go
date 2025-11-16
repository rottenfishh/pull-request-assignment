package env

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type ConfigDb struct {
	DbUsername string `env:"DB_USERNAME,required"`
	DbPassword string `env:"DB_PASSWORD,required"`
	DbURL      string `env:"DB_URL,required"`
}

func LoadConfigEnv() (*ConfigDb, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file %w", err)
	}

	configDb := ConfigDb{}

	err = env.Parse(&configDb)
	if err != nil {
		return nil, err
	}

	return &configDb, nil
}
