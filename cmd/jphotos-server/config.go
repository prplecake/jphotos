package main

import (
	"os"

	"github.com/prplecake/jphotos/app"
)

func defaultConfig() app.Configuration {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	config := app.Configuration{
		App: app.Config{
			Port: port,
		},
		DB: app.DatabaseConfig{
			Username: os.Getenv("USER"),
			Hostname: "localhost",
			Port:     5432,
		},
		Templates: app.TemplateConfig{
			"templates",
		},
		Uploads: app.UploadConfig{
			"data/uploads/photos/",
			"data/thumbnails/",
		},
	}
	return config
}
