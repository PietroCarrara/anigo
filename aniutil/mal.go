package aniutil

import (
	"fmt"
	"github.com/PietroCarrara/anigo/anidata"
	"github.com/antchfx/xmlquery"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func FromMAL(user string) []anidata.Anime {

	const url = "https://myanimelist.net/malappinfo.php?status=all&type=anime&u="

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

func Update(a anidata.Anime, usr, pwd string) {

	log.Printf("Updating %s\n", a.Title)

	id := malID(a, usr, pwd)

	u := "https://" + url.QueryEscape(usr) + ":" + url.QueryEscape(pwd) + "@myanimelist.net/api/animelist/update/" + id + ".xml"

	// TODO: Post score together
	data := `<?xml version="1.0" encoding="utf-8"?> <entry>
        <episode>` + fmt.Sprint(a.Completed) + `</episode> 
        <status>` + statusToMal(a.Status) + `</status>
        <score></score>
        <storage_type></storage_type>
        <storage_value></storage_value>
        <times_rewatched></times_rewatched>
        <rewatch_value></rewatch_value>
        <date_start></date_start>
        <date_finish></date_finish>
        <priority></priority>
        <enable_discussion></enable_discussion>
        <enable_rewatching></enable_rewatching>
        <comments></comments>
        <tags></tags>                 
</entry>
`

	v := url.Values{}
	v.Add("data", data)

	http.PostForm(u, v)
}

func malID(a anidata.Anime, usr, pwd string) string {

	search := "https://" + url.QueryEscape(usr) + ":" + url.QueryEscape(pwd) + "@myanimelist.net/api/anime/search.xml?q=" + url.QueryEscape(a.Title)

	res, err := http.Get(search)
	if err != nil {
		log.Fatalf("Could not reach %s: %s\n", search, err.Error())
	}

	doc, err := xmlquery.Parse(res.Body)
	if err != nil || doc == nil {
		log.Fatalf("Error while parsing %s: %s\n", search, err.Error())
	}

	return getID(a, doc)
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

func statusToMal(s anidata.Status) string {
	switch s {
	case anidata.PlanToWatch:
		return "6"
	case anidata.Watching:
		return "1"
	case anidata.Completed:
		return "2"
	case anidata.Dropped:
		return "4"
	default:
		return ""
	}
}

func getID(a anidata.Anime, doc *xmlquery.Node) string {

	for _, node := range doc.SelectElement("anime").SelectElements("entry") {
		if a.Title == xmlquery.FindOne(node, "//title").InnerText() {
			return xmlquery.FindOne(node, "//id").InnerText()
		}
	}

	return xmlquery.FindOne(doc, "//id").InnerText()
}
