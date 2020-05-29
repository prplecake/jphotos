package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"

	"git.sr.ht/~mjorgensen/jphotos"
	"git.sr.ht/~mjorgensen/jphotos/app"
	"git.sr.ht/~mjorgensen/jphotos/db"
)

var (
	config app.Configuration
)

func processArgs() []string {
	argsWithoutProg := os.Args[1:]
	log.Print(argsWithoutProg)
	return argsWithoutProg
}

func main() {
	log.Print("Initializing...")
	config = defaultConfig()
	configFile := "cmd/jphotos-server/jphotos.yaml"
	cf, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panic(err)
	}
	err = yaml.Unmarshal(cf, &config)
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := db.NewPGStore(config.DB.Username, config.DB.Password, config.DB.Name)
	if err != nil {
		log.Fatal(err)
	}

	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	migrateUp := migrateCmd.Bool("up", false, "migrate up")
	migrateDown := migrateCmd.Bool("down", false, "migrate down")
	migrateForce := migrateCmd.Bool("force", false, "force migration")
	migrateVersion := migrateCmd.Int("version", 0, "migrations to run")

	userCmd := flag.NewFlagSet("user", flag.ExitOnError)
	userCreate := userCmd.Bool("create", false, "create user")
	userDelete := userCmd.Bool("delete", false, "delete user")
	userName := userCmd.String("username", "", "username")

	var direction string
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.Username, config.DB.Password, config.DB.Hostname,
		config.DB.Port, config.DB.Name)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			migrateCmd.Parse(os.Args[2:])
			log.Print("subcommand 'migrate'")
			log.Print("\tup:", *migrateUp)
			log.Print("\tdown:", *migrateDown)
			log.Print("\tforced?", *migrateForce)
			log.Print("\tver:", *migrateVersion)
			if *migrateForce {
				direction = "force"
			} else {
				if *migrateUp && *migrateDown {
					log.Fatal("Please only choose up or down.")
				} else {
					if *migrateUp {
						direction = "up"
					}
					if *migrateDown {
						direction = "down"
					}
				}
			}
			dbMigrate(direction, dbURL, *migrateVersion)
		case "user":
			userCmd.Parse(os.Args[2:])
			if *userCreate {
				createUser(postgres)
			}
			if *userDelete {
				deleteUser(*userName, postgres)
			}
		default:
			log.Print("Subcommand not understood.")
		}
	} else {
		runServer()
	}
}

func runServer() {
	err := os.MkdirAll(config.Uploads.Path, 0755)
	if err != nil {
		log.Panic(err)
	}
	err = os.MkdirAll(config.Uploads.ThumbnailsPath, 0755)
	if err != nil {
		log.Panic(err)
	}

	app.InitTemplates(config.Templates.Path + "/**")

	// Get current git version
	app.CurrentVersion = getGitTag()

	postgres, err := db.NewPGStore(config.DB.Username, config.DB.Password, config.DB.Name)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(
		":"+config.App.Port,
		jphotos.NewServer(postgres, config)))
}

func getGitTag() string {
	gitCmd := exec.Command("git", "describe", "--tag")
	outBytes, _ := gitCmd.Output()
	out := string(outBytes)
	return out
}
