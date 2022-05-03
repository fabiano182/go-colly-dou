package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type crawlURL struct {
	TypeNormDay       typeNormDayStruct `json:"typeNormDay"`
	IdPortletInstance string            `json:"idPortletInstance"`
	DateUrl           string            `json:"dateUrl"`
	Section           string            `json:"section"`
	JsonArray         []crawlURLArray   `json:"jsonArray"`
}

type crawlURLArray struct {
	PubName            string   `json:"pubName"`
	UrlTitle           string   `json:"urlTitle"`
	NumberPage         string   `json:"numberPage"`
	SubTitulo          string   `json:"subTitulo"`
	Titulo             string   `json:"titulo"`
	Title              string   `json:"title"`
	PubDate            string   `json:"pubDate"`
	Content            string   `json:"content"`
	EditionNumber      string   `json:"editionNumber"`
	HierarchyLevelSize int      `json:"hierarchyLevelSize"`
	ArtType            string   `json:"artType"`
	PubOrder           string   `json:"pubOrder"`
	HierarchyStr       string   `json:"hierarchyStr"`
	HierarchyList      []string `json:"hierarchyList"`
}

type typeNormDayStruct struct {
	DO2ESP bool `json:"DO2ESP"`
	DO1ESP bool `json:"DO1ESP"`
	DO1A   bool `json:"DO1A"`
	DO3E   bool `json:"DO3E"`
	DO2E   bool `json:"DO2E"`
	DO1E   bool `json:"DO1E"`
}

func main() {

	collector := colly.NewCollector(
		colly.AllowedDomains(
			"www.in.gov.br/",
			"in.gov.br/",
			"https://www.in.gov.br/",
			"www.in.gov.br",
			"in.gov.br",
			"https://www.in.gov.br",
			"https://in.gov.br",
			"https://in.gov.br/",
		),
	)

	collector.SetRequestTimeout(120 * time.Second)

	_, numberOfUrls := getURLsToScrape(*collector)

	fmt.Println(numberOfUrls)
}

func getURLsToScrape(c colly.Collector) ([]string, int) {
	urlsToScrape := []string{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnXML("//script[@id='params']/text()", func(element *colly.XMLElement) {
		jsonURLs := []byte(element.Text)
		var crawledURLs crawlURL
		err := json.Unmarshal(jsonURLs, &crawledURLs)
		if err != nil {
			panic(err)
		}

		for _, urlCrawled := range crawledURLs.JsonArray {
			url := "https://www.in.gov.br/en/web/dou/-/" + urlCrawled.UrlTitle
			urlsToScrape = append(urlsToScrape, url)
		}
	})

	c.Visit("https://www.in.gov.br/leiturajornal?secao=dou3&data=02-05-2022")

	return urlsToScrape, len(urlsToScrape)
}
