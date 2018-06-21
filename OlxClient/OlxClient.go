package olxclient

import (
	"fmt"
	cfg2 "github.com/artemzi/olx-parser/cfg"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"
	"github.com/artemzi/olx-parser/entities"
)

func GetInstance() (c *colly.Collector) {
	// Instantiate default collector
	c = colly.NewCollector(
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"),
		colly.Async(true),
		colly.CacheDir("./storage/.cache"),
		//colly.Debugger(&debug.LogDebugger{}),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*olx.*",
		RandomDelay: 1 * time.Second,
	})

	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   40 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	return
}

func InitLogs(collector ...*colly.Collector) {
	for _, c := range collector {
		c.OnRequest(func(r *colly.Request) {
			log.Info("Crawler visiting: ", r.URL.String())
		})

		c.OnError(func(r *colly.Response, err error) {
			log.Error("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})
	}
}

func prettifyString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func parse() (adverts []*entities.Adverts) {
	cfg := cfg2.NewRequestCfg()
	c := GetInstance()
	dc := GetInstance()
	InitLogs(c, dc)

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	// visit each advert
	c.OnHTML(".wrap .offer .space a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		dc.Visit(link)
	})

	// follow available pagination link
	c.OnHTML(".next a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// parse information from details page
	dc.OnHTML(".offercontent", func(e *colly.HTMLElement) {
		var (
			images []string
			details []*entities.DetailsItem
		)

		title := e.ChildText(".offer-titlebox h1")
		// example: Донецк, Донецкая область, Калининский
		detailsPlace := e.ChildText(".offer-titlebox__details .show-map-link strong")
		// example: Опубликовано с мобильного в 23:55, 20 июня 2018, Номер объявления: 540309546
		detailsMeta := prettifyString(e.ChildText(".offer-titlebox__details em"))
		e.ForEach(".img-item", func(_ int, el *colly.HTMLElement) {
			image := e.ChildAttr(".photo-glow img", "src")
			images = append(images, image)
		})
		text := e.ChildText("#textContent p")
		e.ForEach(".descriptioncontent .details .item", func(_ int, el *colly.HTMLElement) {
			name := el.ChildText("th")
			value := prettifyString(el.ChildText(".value a"))
			if "" == value {
				value = el.ChildText(".value strong")
			}
			details = append(details, &entities.DetailsItem{name, value})
		})
		price := e.ChildText(".price-label strong")

		adverts = append(adverts, &entities.Adverts{
			URL: e.Request.URL.String(),
			Title: title,
			Place: detailsPlace,
			Meta: detailsMeta,
			Details: details,
			Images: images,
			Text: strings.TrimSpace(text),
			Price: price,
		})
	})

	// 500 - num pages to visit, there is no need in more than current value
	for i := 1; i <= 500; i++ {
		q.AddURL(fmt.Sprintf("%s/%s/%s/%s?page=%d", cfg.BASEURL, cfg.CATEGORY, cfg.SUBCATEGORY, cfg.REGION, i))
	}
	q.Run(c)
	c.Wait()
	dc.Wait()
	return
}
