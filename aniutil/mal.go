package aniutil

import (
	"github.com/PietroCarrara/anigo/anidata"
	"github.com/antchfx/xmlquery"
	"log"
	"net/http"
	"strconv"
)

const url = "https://myanimelist.net/malappinfo.php?status=all&type=anime&u="

func FromMAL(user string) []anidata.Anime {
	animeList := []anidata.Anime{}

	req, err := http.Get(url + user)
	if err != nil {
		log.Fatalf("Could not reach %s\n", url+user)
	}

	doc, err := xmlquery.Parse(req.Body)
	if err != nil {
		log.Fatalf("Error while parsing %s: %s\n", url+user, err.Error())
	}

	for _, node := range xmlquery.Find(doc, "//anime") {

		anime := anidata.Anime{}

		chapters, err := strconv.ParseInt(node.SelectElement("series_episodes").InnerText(), 10, 32)
		if err != nil {
			log.Printf("Error while parsing chapter numbers from MAL: %s", err.Error())
		}

		completed, err := strconv.ParseInt(node.SelectElement("my_watched_episodes").InnerText(), 10, 32)
		if err != nil {
			log.Printf("Error while parsing watched chapter numbers from MAL: %s", err.Error())
		}

		anime.Title = node.SelectElement("series_title").InnerText()
		anime.Chapters = int(chapters)
		anime.Completed = int(completed)
		anime.Status = statusFromMal(node.SelectElement("my_status").InnerText())

		animeList = append(animeList, anime)
	}

	return animeList
}

func statusFromMal(s string) anidata.Status {
	switch s {
	case "6":
		return anidata.PlanToWatch
	case "1":
		return anidata.Watching
	case "2":
		return anidata.Completed
	case "4":
		return anidata.Dropped
	default:
		return -1
	}
}
