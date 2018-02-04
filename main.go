//
// main.go
// Copyright (C) 2018 pietro <pietro@the-arch>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"encoding/json"
	"github.com/PietroCarrara/anigo/anidata"
	"log"
	"os"
)

var Home string

var Database []anidata.Anime

func main() {

	// Set config home
	Home = os.Getenv("XDG_CONFIG_HOME")
	if Home == "" {
		Home = os.Getenv("HOME") + "/.config"
	}
	Home += "/anigo"

	// Loading previous entries
	dataBaseFileName := Home + "/database.json"
	databaseFile, err := os.Open(dataBaseFileName)

	if err != nil {
		log.Printf("%s could not be open, initializing new database...\n", dataBaseFileName)
		// Empty database
		Database = []anidata.Anime{}
	} else {
		log.Printf("Found database in %s, loading...\n", dataBaseFileName)

		dec := json.NewDecoder(databaseFile)
		err := dec.Decode(&Database)
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("Database empty, initializing new...")
				Database = []anidata.Anime{}
			} else {
				log.Fatalf("Error decoding database: %s", err.Error())
			}
		}
	}
	// Closing since we are done reading
	databaseFile.Close()

	anim := anidata.Anime{Name: "Title", Chapters: 25, Status: anidata.PlanToWatch}

	Database = append(Database, anim)

	databaseFile, _ = os.Create(dataBaseFileName)
	enc := json.NewEncoder(databaseFile)
	enc.SetIndent(" ", " ")
	err = enc.Encode(Database)
	if err != nil {
		log.Fatal(err)
	}
	databaseFile.Close()
}
