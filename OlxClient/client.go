package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	"strings"
	"net/url"
	"strconv"
)

// Defaults (will be changed in future)
type Cfg struct {
	BASEURL string
	CATEGORY string
	SUBCATEGORY string
	REGION string
}

type OlxHttpClient struct {
	http *http.Client
}

type Node struct {
	URL     string
	Title   string
	Time    string
}

func (c *OlxHttpClient) OlxGet(url string) *goquery.Document {
	log.Debug("OlxGet ", url)

	res, err := c.http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return doc
}

func (c *OlxHttpClient) OlxGetPages(u string) []Node {
	var wg sync.WaitGroup
	log.Debug("OlxGetPages ", u)

	doc := c.OlxGet(u)
	var nodes []Node

	// Load the HTML document
	doc.Find("#offers_table").Each(func(i int, s *goquery.Selection) {
		s.Find(".wrap").Each(func(i int, s *goquery.Selection) {
			node := Node{}
			title := s.Find("h3")
			u, _ := title.Find("a").Attr("href")
			time := s.Find(".space").Eq(2).Find("p").Last().Text()

			node.Title = strings.TrimSpace(title.Text())
			node.URL = strings.TrimSpace(u)
			node.Time = strings.TrimSpace(time)

			nodes = append(nodes, node)
		})
	})

	pageUrl, ok := doc.Find(".pager").Find(".next").Find("a").Attr("href")
	// reduce size for debug
	urlString, _ := url.Parse(pageUrl)
	page, _ := strconv.Atoi(urlString.Query().Get("page"))
	// ===
	if ok && page <= 10 {
		wg.Add(1)
		go func(d *goquery.Document) {
			defer wg.Done()
		}(doc)
	}

	wg.Wait()
	return nodes
}

func main() {
	cfg := &Cfg{
		"https://www.olx.ua",
		"transport",
		"legkovye-avtomobili",
		"donetsk",
	}

	u := fmt.Sprintf("%s/%s/%s/%s", cfg.BASEURL, cfg.CATEGORY, cfg.SUBCATEGORY, cfg.REGION)
	cli := &OlxHttpClient{&http.Client{}}
	data := cli.OlxGetPages(u)

	fmt.Printf("%v\n", data)
}