package util

import (
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
	{Name: "title", Usage: []string{"-w", "--watched", "watched"}, Description: "Set number of watched chapters", UsesNext: []bool{true, false, false}},

	// Statuses
	{Name: "status", Usage: []string{"-P"}, Description: "Set status as 'Plan To Watch'", Default: "p"},
	{Name: "status", Usage: []string{"-W"}, Description: "Set status as 'Watching'", Default: "w"},
	{Name: "status", Usage: []string{"-C"}, Description: "Set status as 'Completed'", Default: "c"},
	{Name: "status", Usage: []string{"-D"}, Description: "Set status as 'Dropped'", Default: "d"},

	// Command
	{Name: "command", Usage: []string{"-Q", "search", "query"}, Description: "Do a search", Default: "search"},
	{Name: "command", Usage: []string{"-A", "add"}, Description: "Add an entry to the database", Default: "add"},
}

var Args map[string]string = map[string]string{}

func init() {

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
