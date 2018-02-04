package anidata

import (
	"errors"
	"strings"
)

type Status int

const (
	PlanToWatch Status = iota
	Watching
	Completed
	Dropped
)

func (s Status) String() string {
	switch s {
	case PlanToWatch:
		return "Plan To Watch"
	case Watching:
		return "Watching"
	case Completed:
		return "Completed"
	case Dropped:
		return "Dropped"
	default:
		return "Status Unknown"
	}
}

func StatusFromString(s string) (Status, error) {
	switch strings.ToLower(s) {
	case "p", "ptw", "plan to watch":
		return PlanToWatch, nil
	case "w", "wtc", "watching":
		return Watching, nil
	case "c", "cpl", "completed":
		return Completed, nil
	case "d", "drp", "dropped":
		return Dropped, nil
	default:
		return -1, errors.New("Unknown status")
	}
}
