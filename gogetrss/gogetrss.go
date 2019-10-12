package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// URL where to read Atom feed
var URL string

// Feed represents the Atom feed
type Feed struct {
	Entries []FeedEntry `xml:"entry"`
}

// FeedEntry represent single atom entry with data
type FeedEntry struct {
	Title   string   `xml:"title" json:"title"`
	Updated string   `xml:"updated"`
	Link    LinkItem `xml:"link"`
}

type LinkItem struct {
	Link string `xml:"href,attr"`
}

// enmtrypoint to the program
func main() {
	log.Printf("Started parsing %s", URL)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatalf("Cannot create request %s", err.Error())
	}

	req.Header.Add("User-Agent", "Golang bot 1.0")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Cannot get response %s", err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Cannot read response %s", err.Error())
	}

	feed := Feed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		log.Fatalf("Cannot parse response: %s", err.Error())
	}

	for index, item := range feed.Entries {
		fmt.Printf("%d: %s, %s, %s\n", index, item.Updated, item.Title, item.Link.Link)
	}
}

// initialization
func init() {
	flag.StringVar(&URL, "URL", "", "Atom feed url")
	flag.Parse()
}
