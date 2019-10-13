package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"git.sr.ht/~mjorgensen/jphotos"
	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

func main() {
	log.Print("Initializing...")
	config := defaultConfig()
	configFile := "cmd/jphotos-server/jphotos.yaml"
	cf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panic(err)
	}
	err = yaml.Unmarshal(cf, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(config)

	err = os.MkdirAll(config.Uploads.Path, 0755)
	if err != nil {
		log.Panic(err)
	}

	app.InitTemplates(config.Templates.Path + "/**")
	postgres, err := db.NewPGStore(config.DB.Username, config.DB.Password, config.DB.DBName)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(
		":"+config.App.Port,
		jphotos.NewServer(postgres)))

}
