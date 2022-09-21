package main

import (
	"log"

	"github.com/Cranky4/go-top/internal/top"
	"github.com/spf13/viper"
)

func NewConfig(path string) top.Config {
	viper.SetConfigFile(path)
	var c top.Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return c
}
