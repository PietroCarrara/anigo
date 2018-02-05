package util

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"log"
	"os"
	"strings"
)

type argument struct {
	Name        string
	Usage       []string
	Description string
	UsesNext    []bool
	Default     string
}

var arguments []argument = []argument{
	{Name: "debug", Usage: []string{"--debug"}, Description: "Activate debug mode", Default: "true"},

	{Name: "title", Usage: []string{"-t", "--title", "title"}, Description: "Set title", UsesNext: []bool{true, false, false}},
	{Name: "chapters", Usage: []string{"-c", "--chapters", "chapters"}, Description: "Set total number of chapters", UsesNext: []bool{true, false, false}},
	{Name: "completed", Usage: []string{"-w", "--watched", "watched"}, Description: "Set number of watched chapters", UsesNext: []bool{true, false, false}},

	{Name: "set-title", Usage: []string{"--set-title"}, Description: "Set title for edit mode", UsesNext: []bool{false}},
	{Name: "set-chapters", Usage: []string{"--set-chapters"}, Description: "Set total number of chapters for edit mode", UsesNext: []bool{false}},
	{Name: "set-completed", Usage: []string{"--set-completed"}, Description: "Set number of watched chapters for edit mode", UsesNext: []bool{false}},
	{Name: "set-status", Usage: []string{"--set-status"}, Description: "Set status for edit mode", UsesNext: []bool{false}},

	{Name: "user", Usage: []string{"-u", "--user", "user"}, Description: "Define the user for MyAnimeList.net", UsesNext: []bool{true, false, false}},

	// Statuses
	{Name: "status", Usage: []string{"-P"}, Description: "Set status as Plan To Watch", Default: "p"},
	{Name: "status", Usage: []string{"-W"}, Description: "Set status as Watching", Default: "w"},
	{Name: "status", Usage: []string{"-C"}, Description: "Set status as Completed", Default: "c"},
	{Name: "status", Usage: []string{"-D"}, Description: "Set status as Dropped", Default: "d"},

	// Command
	{Name: "command", Usage: []string{"-Q", "search", "query"}, Description: "Do a search", Default: "search"},
	{Name: "command", Usage: []string{"-A", "add"}, Description: "Add an entry to the database", Default: "add"},
	{Name: "command", Usage: []string{"pull"}, Description: "Pull your entries from MyAnimeList.net", Default: "pull"},
	{Name: "command", Usage: []string{"-E", "edit"}, Description: "Edit all entries mathcing the criteria. To set values use --set-value", Default: "edit"},
	{Name: "command", Usage: []string{"--spit"}, Description: "Spit autocomplete", Default: "spit"},
}

var Args map[string]string = map[string]string{}

func init() {

	// Default operation is search
	Args["command"] = "search"

	// Parse command line arguments
	for i := 1; i < len(os.Args); i++ {
		parts := strings.Split(os.Args[i], "=")
		for _, val := range arguments {
			index := in(parts[0], val)
			if index >= 0 {
				arg := strings.Join(parts[1:], "")
				if val.Default != "" {
					arg = val.Default
				} else if val.UsesNext[index] {
					i++
					arg = os.Args[i]
				}
				Args[val.Name] = arg
				break
			}
		}
	}
}

func in(value string, arg argument) int {

	for index, val := range arg.Usage {
		if val == value {
			return index
		}
	}

	return -1
}

func SpitAutoComplete(home string) {

	template := home + "/template.sh"

	vars := map[string]string{}

	vars["database"] = home + "/database.json"

	for _, val := range whereName("command") {
		for _, usage := range val.Usage {
			vars["commands"] += fmt.Sprintf("'%s:%s'\n\t", usage, val.Description)
		}
	}

	for _, val := range arguments {
		for i, usage := range val.Usage {
			cmd := usage
			if len(val.UsesNext) > i && !val.UsesNext[i] {
				cmd += "="
			}
			log.Printf("Usage: %s: ", cmd)
			vars["all"] += fmt.Sprintf("'%s:%s'\n\t", cmd, val.Description)
		}
	}

	res, err := mustache.RenderFile(template, vars)
	if err != nil {
		log.Fatalf("Error while processing auto-complete template: %s", err.Error())
	} else {
		fmt.Println(res)
	}
}

func whereName(n string) []argument {
	res := []argument{}

	for _, val := range arguments {
		if val.Name == n {
			res = append(res, val)
		}
	}

	return res
}
