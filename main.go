//
// main.go
// Copyright (C) 2018 pietro <pietro@the-arch>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/PietroCarrara/anigo/anidata"
	"github.com/PietroCarrara/anigo/aniutil"
	"github.com/PietroCarrara/anigo/util"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var Home string

var Database []anidata.Anime

func main() {

	if util.Args["debug"] != "true" {
		log.SetOutput(ioutil.Discard)
	}

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

	log.Println(util.Args)

	switch util.Args["command"] {
	case "add":
		add()
	case "search":
		search()
	case "pull":
		pull()
	}

	sort.Slice(Database, func(i, j int) bool { return Database[i].Title < Database[j].Title })

	databaseFile, _ = os.Create(dataBaseFileName)
	enc := json.NewEncoder(databaseFile)
	enc.SetIndent(" ", " ")
	err = enc.Encode(Database)
	if err != nil {
		log.Fatal(err)
	}
	databaseFile.Close()
}

func add() {

	var err error

	var chapters int64
	if util.Args["chapters"] != "" {
		chapters, err = strconv.ParseInt(util.Args["chapters"], 10, 32)
		if err != nil {
			log.Fatalf("Invalid number informed for chapters: %s\n", err.Error())
		}
	}

	var completed int64
	if util.Args["completed"] != "" {
		completed, err = strconv.ParseInt(util.Args["completed"], 10, 32)
		if err != nil {
			log.Fatalf("Invalid number informed for completed: %s\n", err.Error())
		}
	}

	status, err := anidata.StatusFromString(util.Args["status"])
	if err != nil {
		log.Println("Either no status was informed or a invalid one was. Assuming 'Unknow Status'...")
	}

	anim := anidata.Anime{
		Title:     util.Args["title"],
		Chapters:  int(chapters),
		Completed: int(completed),
		Status:    status,
	}

	anim.Fix()

	Database = append(Database, anim)
}

func search() {
	for _, a := range Database {
		if matchCriteria(a) {
			print(a)
		}
	}

}

func pull() {

	if util.Args["user"] == "" {
		fmt.Println("user was not informed. Exiting...")
		return
	}

	// TODO: check for repeated entries and update propperly,
	// rather than just overwrite
	Database = aniutil.FromMAL(util.Args["user"])
}

func matchCriteria(a anidata.Anime) bool {
	if util.Args["title"] != "" && strings.ToLower(util.Args["title"]) != strings.ToLower(a.Title) {
		return false
	}

	if util.Args["chapters"] != "" && fmt.Sprint(a.Chapters) != util.Args["chapters"] {
		return false
	}

	if util.Args["completed"] != "" && fmt.Sprint(a.Completed) != util.Args["completed"] {
		return false
	}

	status, err := anidata.StatusFromString(util.Args["status"])
	if err != nil {
		log.Println("Either no status was informed or a invalid one was. Assuming 'Unknow Status'...")
	} else if status != a.Status {
		return false
	}

	return true
}

func print(a anidata.Anime) {
	completition := 0
	if a.Chapters > 0 {
		completition = a.Completed * 100 / a.Chapters
	}

	fmt.Printf("Title: %s\n", a.Title)
	fmt.Printf("Chapters: %d\tWatched: %d (%d%%)\n", a.Chapters, a.Completed, completition)

	fmt.Println()
}
