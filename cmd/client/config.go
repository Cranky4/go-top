package main

import (
	"log"

	topclient "github.com/Cranky4/go-top/internal/top-client"
	"github.com/spf13/viper"
)

func NewConfig(path string) topclient.Config {
	viper.SetConfigFile(path)
	var c topclient.Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}

	return c
}
