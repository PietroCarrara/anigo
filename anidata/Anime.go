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

	if a.Completed > a.Chapters && a.Chapters > 0 {
		a.Completed = a.Chapters
	}

	if a.Status == Completed {
		a.Completed = a.Chapters
	} else if a.Chapters == a.Completed {
		a.Status = Completed
	}
}

// Merges anime and return true
// if the anime was modified
func (a *Anime) Merge(b Anime) bool {

	modified := false

	if a.Chapters < b.Chapters {
		a.Chapters = b.Chapters
		modified = true
	}

	if a.Completed < b.Completed {
		a.Completed = b.Completed
		modified = true
	}

	if a.Status < b.Status {
		a.Status = b.Status
		modified = true
	}

	return modified
}
