package main

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Postgres struct {
		Host              string
		Database          string
		Port              int
		User              string
		Password          string
		ConnectionTimeout int
	}
	Smart struct {
		Path string
	}
}

func getConfig() (*Config, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath(homeDir + "/.config/")
	viper.SetConfigName("mac-stats")

	viper.SetDefault("Postgres.Host", "localhost")
	viper.SetDefault("Postgres.Database", "macstats")
	viper.SetDefault("Postgres.Port", "5432")
	viper.SetDefault("Postgres.connectionTimeout", 10)

	// assuming smartctl was brewed
	viper.SetDefault("Smart.path", "/usr/local/bin")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config not found in current dir or in ~/.config")
		} else {
			return nil, err
		}
	}
	var conf Config

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}
	return &conf, nil
}
