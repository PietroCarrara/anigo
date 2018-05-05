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
	"github.com/fatih/color"
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

	// Set config home
	Home = os.Getenv("XDG_CONFIG_HOME")
	if Home == "" {
		Home = os.Getenv("HOME") + "/.config"
	}
	Home += "/anigo"

	// Config file
	configFileName := Home + "/config"
	configFile, err := os.Open(configFileName)
	if err != nil {
		util.Parse(os.Args[1:])
	} else {
		bytes, _ := ioutil.ReadAll(configFile)

		contents := string(bytes)

		// Remove break lines and split by spaces
		args := strings.Split(strings.TrimSpace(strings.Replace(contents, "\n", " ", -1)), " ")

		args = append(args, os.Args[1:]...)

		util.Parse(args)
	}

	if util.Args["debug"] != "true" {
		log.SetOutput(ioutil.Discard)
	}

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
				util.Explode("Error decoding database: " + err.Error())
			}
		}
	}
	// Closing since we are done reading
	databaseFile.Close()

	if util.Args["title"] == "" && util.Args["command"] != "search" {

		red := color.New(color.BgRed).SprintFunc()

		fmt.Printf("%s You have provided no title for the operation. This could go very wrong.\nContinue? [y/N]", red("ATTENTION:"))

		var c rune
		fmt.Scanf("%c", &c)

		if c != 'y' && c != 'Y' {
			return
		}
	}

	switch util.Args["command"] {
	case "add":
		add()
	case "delete":
		del()
	case "search":
		search()
	case "pull":
		pull()
	case "push":
		push()
	case "edit":
		edit()
	case "spit":
		util.SpitAutoComplete(Home)
	}

	sort.Slice(Database, func(i, j int) bool { return Database[i].Title < Database[j].Title })

	databaseFile, _ = os.Create(dataBaseFileName)
	defer databaseFile.Close()

	enc := json.NewEncoder(databaseFile)
	enc.SetIndent(" ", " ")
	err = enc.Encode(Database)
	if err != nil {
		util.Explode(err.Error())
	}
}

func add() {

	var err error

	var chapters int64
	if util.Args["chapters"] != "" {
		chapters, err = strconv.ParseInt(util.Args["chapters"], 10, 32)
		if err != nil {
			util.Explode("Invalid number informed for chapters: " + err.Error())
		}
	}

	var completed int64
	if util.Args["completed"] != "" {
		completed, err = strconv.ParseInt(util.Args["completed"], 10, 32)
		if err != nil {
			util.Explode("Invalid number informed for completed: " + err.Error())
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

func del() {

	var newDB = []anidata.Anime{}

	for _, a := range Database {
		if !matchCriteria(a) {
			newDB = append(newDB, a)
		}
	}

	Database = newDB
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
		util.Explode("User not informed!")
		return
	}

	// TODO: check for repeated entries and update propperly,
	// rather than just overwrite
	Database = aniutil.FromMAL(util.Args["user"])
}

func push() {

	if util.Args["user"] == "" || util.Args["password"] == "" {
		util.Explode("Credentials missing!")
		return
	}

	origin := aniutil.FromMAL(util.Args["user"])

	for i, _ := range Database {
		anime := byTitle(Database[i], origin)
		if anime != nil {
			if anime.Merge(Database[i]) {
				anime.Fix()
				aniutil.Update(*anime, util.Args["user"], util.Args["password"])
			}
		} else {
			// TODO: Insert()
		}
	}
}

func edit() {

	var err error

	var chapters int64 = -1
	if util.Args["set-chapters"] != "" {
		chapters, err = strconv.ParseInt(util.Args["set-chapters"], 10, 32)
		if err != nil {
			util.Explode("Invalid number informed for editing chapters: " + err.Error())
		}
	}

	var completed int64 = -1
	if util.Args["set-completed"] != "" {
		completed, err = strconv.ParseInt(util.Args["set-completed"], 10, 32)
		if err != nil {
			util.Explode("Invalid number informed for editing completed: " + err.Error())
		}
	}

	status, err := anidata.StatusFromString(util.Args["set-status"])
	if err != nil {
		log.Println("Either no status was informed or a invalid one was. Assuming 'Unknow Status' for editing...")
	}

	for i := 0; i < len(Database); i++ {
		if matchCriteria(Database[i]) {
			if util.Args["set-title"] != "" {
				Database[i].Title = util.Args["set-title"]
			}
			if chapters >= 0 {
				Database[i].Chapters = int(chapters)
			}
			if completed >= 0 {
				Database[i].Completed = int(completed)
			}
			if status >= 0 {
				Database[i].Status = status
			}
			Database[i].Fix()
		}
	}
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

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	// black := color.New(color.FgBlack).SprintFunc()

	fmt.Printf("Title: %s\n", yellow(a.Title))
	fmt.Printf("Chapters: %s\tWatched: %s (%s%%)\n", cyan(a.Chapters), cyan(a.Completed), magenta(completition))

	var statColor func(...interface{}) string
	switch a.Status {
	case anidata.PlanToWatch:
		statColor = white
	case anidata.Watching:
		statColor = green
	case anidata.Completed:
		statColor = blue
	case anidata.Dropped:
		statColor = red
	}

	fmt.Printf("Status: %s\n", statColor(a.Status))

	fmt.Println()
}

func byTitle(a anidata.Anime, arr []anidata.Anime) *anidata.Anime {

	for _, b := range arr {
		if strings.ToLower(a.Title) == strings.ToLower(b.Title) {
			return &b
		}
	}

	return nil
}
