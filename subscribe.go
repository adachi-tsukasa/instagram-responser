package main

import (
	"io/ioutil"
	"net/http"

	"encoding/xml"
	"html/template"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// RSS一覧用
type RssListItem struct {
	Title string
	Url   string
}

// RSS構造体
type Rss2 struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	Title       string   `xml:"channel>title"`
	Link        string   `xml:"channel>link"`
	Description string   `xml:"channel>description"`
	PubDate     string   `xml:"channel>pubDate"`
	ItemList    []Item   `xml:"channel>item"`
}

type Item struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`
	Content     template.HTML `xml:"encoded"`
	PubDate     string        `xml:"pubDate"`
	Comments    string        `xml:"comments"`
}

type RssJsonList struct {
	Title string        `json:"title"`
	List  []RssJsonItem `json:"list"`
}
type RssJsonItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}

func getLatestFeed(w http.ResponseWriter, r *http.Request, url string) RssJsonItem {
	c := urlfetch.Client(appengine.NewContext(r))
	res, err := c.Get(url)
	if err != nil {
		panic(err)
	}

	asText, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		panic(err2)
	}

	rssData := Rss2{}

	err3 := xml.Unmarshal(asText, &rssData)
	if err3 != nil {
		panic(err3)
	}

	var jsonData = RssJsonList{}
	jsonData.Title = rssData.Title
	jsonData.List = []RssJsonItem{}

	for _, value := range rssData.ItemList {
		jsonData.List = append(jsonData.List, RssJsonItem{Title: value.Title, Link: value.Link, Date: value.PubDate})
		break
	}

	return jsonData.List[0]

}
