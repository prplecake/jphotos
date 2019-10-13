package main

import (
	"os"
)

type configuration struct {
	App       appConfig
	DB        databaseConfig `yaml:"database"`
	Templates templateConfig
}

type appConfig struct {
	Port string
}

type databaseConfig struct {
	Username, Password string
	DBName             string `yaml:"name"`
}

type templateConfig struct {
	Path string
}

func defaultConfig() configuration {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	config := configuration{
		App: appConfig{
			Port: port,
		},
		DB: databaseConfig{
			Username: os.Getenv("USER"),
		},
		Templates: templateConfig{
			"templates",
		},
	}
	return config
}
