package util

import (
	"fmt"
	"github.com/cbroglie/mustache"
)

type argument struct {
	Name        string
	Usage       []string
	Description string
	Default     string
}

var arguments []argument = []argument{
	{Name: "debug", Usage: []string{"--debug"}, Description: "Activate debug mode", Default: "true"},

	{Name: "title", Usage: []string{"-t", "--title"}, Description: "Set title"},
	{Name: "chapters", Usage: []string{"-c", "--chapters"}, Description: "Set total number of chapters"},
	{Name: "completed", Usage: []string{"-w", "--watched"}, Description: "Set number of watched chapters"},

	{Name: "set-title", Usage: []string{"--set-title"}, Description: "Set title for edit mode"},
	{Name: "set-chapters", Usage: []string{"--set-chapters"}, Description: "Set total number of chapters for edit mode"},
	{Name: "set-completed", Usage: []string{"--set-completed"}, Description: "Set number of watched chapters for edit mode"},
	{Name: "set-status", Usage: []string{"--set-status"}, Description: "Set status for edit mode"},

	{Name: "user", Usage: []string{"-u", "--user"}, Description: "Define the user for MyAnimeList.net"},
	{Name: "password", Usage: []string{"-p", "--password"}, Description: "Define the password for MyAnimeList.net"},

	// Statuses
	{Name: "status", Usage: []string{"-P"}, Description: "Set status as Plan To Watch", Default: "p"},
	{Name: "status", Usage: []string{"-W"}, Description: "Set status as Watching", Default: "w"},
	{Name: "status", Usage: []string{"-C"}, Description: "Set status as Completed", Default: "c"},
	{Name: "status", Usage: []string{"-D"}, Description: "Set status as Dropped", Default: "d"},

	// Command
	{Name: "command", Usage: []string{"-Q", "search", "query"}, Description: "Do a search", Default: "search"},
	{Name: "command", Usage: []string{"-A", "add"}, Description: "Add an entry to the database", Default: "add"},
	{Name: "command", Usage: []string{"-R", "del", "delete"}, Description: "Deletes one or more entries matching the params", Default: "delete"},
	{Name: "command", Usage: []string{"pull"}, Description: "Pull your entries from MyAnimeList.net", Default: "pull"},
	{Name: "command", Usage: []string{"push"}, Description: "Push your entries to MyAnimeList.net", Default: "push"},
	{Name: "command", Usage: []string{"-E", "edit"}, Description: "Edit all entries mathcing the criteria. To set values use --set-value", Default: "edit"},
	{Name: "command", Usage: []string{"--spit"}, Description: "Spit autocomplete", Default: "spit"},
}

var Args map[string]string = map[string]string{}

func Parse(args []string) {

	// Default operation is search
	Args["command"] = "search"

	// Parse command line arguments
	for i := 0; i < len(args); i++ {
		for _, val := range arguments {
			if in(args[i], val) {
				arg := args[i]
				if val.Default != "" {
					arg = val.Default
				} else {
					i++
					arg = args[i]
				}
				Args[val.Name] = arg
				break
			}
		}
	}
}

func in(value string, arg argument) bool {

	for _, val := range arg.Usage {
		if val == value {
			return true
		}
	}

	return false
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
		for _, usage := range val.Usage {
			vars["all"] += fmt.Sprintf("'%s:%s'\n\t", usage, val.Description)
		}
	}

	res, err := mustache.RenderFile(template, vars)
	if err != nil {
		Explode("Error while processing auto-complete template: " + err.Error())
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
