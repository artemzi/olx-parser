package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"strings"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net"
	"time"
)

// Defaults (will be changed in future)
type Cfg struct {
	BASEURL     string
	CATEGORY    string
	SUBCATEGORY string
	REGION      string
}

type Adverts struct {
	URL    string
	Title  string
	Place  string
	Meta   string
	Images []string
	Text   string
}

func main() {
	var (
		adverts []Adverts
	)

	cfg := &Cfg{
		"https://www.olx.ua",
		"transport",
		"legkovye-avtomobili",
		"donetsk",
	}

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	// Instantiate default collector
	collector := colly.NewCollector(
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"),
		colly.Async(true),
		colly.CacheDir("./.cache"),
	)

	collector.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   40 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	// Instantiate collector for details page
	detailCollector := collector.Clone()

	// visit each advert
	collector.OnHTML(".wrap .offer .space a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		detailCollector.Visit(link)
	})

	// follow available pagination link
	collector.OnHTML(".next a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// parse information from details page
	detailCollector.OnHTML(".offercontent", func(e *colly.HTMLElement) {
		title := e.ChildText(".offer-titlebox h1")
		// example: Донецк, Донецкая область, Калининский
		detailsPlace := e.ChildText(".offer-titlebox__details .show-map-link strong")
		// example: Опубликовано с мобильного в 23:55, 20 июня 2018, Номер объявления: 540309546
		detailsMeta := strings.Join(strings.Fields(e.ChildText(".offer-titlebox__details em")), " ")
		var images []string
		e.ForEach(".img-item", func(_ int, el *colly.HTMLElement) {
			image := e.ChildAttr(".photo-glow img", "src")
			images = append(images, image)
		})
		text := e.ChildText("#textContent p")

		adverts = append(adverts, Adverts{
			e.Request.URL.String(),
			title,
			detailsPlace,
			detailsMeta,
			images,
			strings.TrimSpace(text),
		})
	})

	// 500 - num pages to visit, there is no need in more than current value
	for i := 1; i <= 500; i++ {
		q.AddURL(fmt.Sprintf("%s/%s/%s/%s?page=%d", cfg.BASEURL, cfg.CATEGORY, cfg.SUBCATEGORY, cfg.REGION, i))
	}
	q.Run(collector)

	collector.Wait()

	// debug
	fmt.Println(len(adverts))
	fmt.Printf("%s\n%s\n%s\n%s\n%s\n%v\n",
		adverts[0].Title, adverts[0].URL, adverts[0].Place, adverts[0].Meta, adverts[0].Text, adverts[0].Images)
}
