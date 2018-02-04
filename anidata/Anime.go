package anidata

import (
	"strings"
)

type Anime struct {
	Title     string
	Chapters  int
	Completed int
	Status    Status
}

// Small fixes
func (a *Anime) Fix() {

	a.Title = strings.TrimSpace(a.Title)

	if a.Status == Completed {
		a.Completed = a.Chapters
	}
}
