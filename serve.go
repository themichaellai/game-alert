package main

import (
	"encoding/xml"
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"time"
)

type Guid struct {
	Guid        string `xml:",chardata"`
	IsPermalink bool   `xml:"isPermaLink,attr"`
}

type FeedElement struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	Guid        Guid     `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
}

func eventToFeedElement(e *Event, g *Game) *FeedElement {
	const rfc2822 = "Mon, 02 Jan 2006 15:04:05 -0700"
	guid := Guid{
		Guid:        fmt.Sprintf("%s-%d", g.PositionId, e.Datetime),
		IsPermalink: false,
	}
	res := FeedElement{
		Title:       fmt.Sprintf("%s - %s vs %s", g.Status, g.Home, g.Away),
		Description: fmt.Sprintf("%s - %s vs %s", g.Status, g.Home, g.Away),
		Link:        "http://google.com",
		Guid:        guid,
		PubDate:     time.Unix(e.Datetime, 0).Format(rfc2822),
	}
	return &res
}

func rssTemplate(inner []byte) string {
	return fmt.Sprintf(`
	<rss version="2.0">
	<channel>
	<description>bball</description>
	<link>http://michael.pizza</link>
	<title>ball is life</title>
	%s
	</channel>
	</rss>
	`, inner)
}
func ServeFeed(dbFilename string) {
	m := martini.Classic()
	m.Get("/", func(res http.ResponseWriter) string {
		db := DBHandle(dbFilename)
		defer db.Close()
		eventgames, err := GetLatestEvents(db)
		if err != nil {
			log.Println(err)
			res.WriteHeader(500)
		}
		elements := make([]*FeedElement, len(eventgames))
		for i, eg := range eventgames {
			elements[i] = eventToFeedElement(eg.Event, eg.Game)
		}
		xml, err := xml.MarshalIndent(elements, " ", "  ")
		if err != nil {
			log.Println(err)
			res.WriteHeader(500)
		}
		return rssTemplate(xml)
	})
	m.Run()
}
