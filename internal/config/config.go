package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MosquittoHost string `envconfig:"MOSQUITTO_HOST"`
}

var Env Config

func StartConfig() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := envconfig.Process("", &Env); err != nil {
		return err
	}
	return nil
}
