package main

import (
	yaac_mvvm "github.com/DHBW-SE-2023/yaac-go-prototype/internal/mvvm"
	yaac_shared "github.com/DHBW-SE-2023/yaac-go-prototype/internal/shared"
	"log"
	"os"
)

func main() {
	// create database if not exists
	if _, err := os.Stat("./data/data.db"); err == nil {
		yaac_shared.CreateDatabase()
	} else {
		log.Println("Found existing database")
	}

	db := yaac_shared.ConnectDatabase()

	defer yaac_shared.DisconnectDatabase(db)

	// TODO: remove this it is a test
	err := yaac_shared.InsertStudent(db, "Mustermann", "Max", true, "TIK22")
	if err != nil {
		log.Fatal("TEST: Could not insert student")
	}

	mvvm := yaac_mvvm.New()
	mvvm.OpenMainWindow()
}
