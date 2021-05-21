package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/prplecake/jphotos/auth"
	"github.com/prplecake/jphotos/db"
)

func createUser(dbs *db.PGStore) {
	username, password := credentials()

	err := auth.AddUser(username, password, dbs)
	if err != nil {
		log.Fatal("An error occured attempting to add user:", err)
	}

	log.Printf("User '%s' created.", username)
}

func deleteUser(username string, dbs *db.PGStore) {
	if len(username) == 0 {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Username: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}

	err := auth.RemoveUser(username, dbs)
	if err != nil {
		log.Fatalf("An error occured attempting to remove user [%s]: %v", username, err)
	}
	log.Printf("User '%s' deleted.", username)
}

func credentials() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}
