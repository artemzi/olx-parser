package olxclient

import (
	"fmt"
	cfg2 "github.com/artemzi/olx-parser/cfg"
	"github.com/artemzi/olx-parser/entities"
	"github.com/artemzi/olx-parser/helpers"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"
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

	rp, err := proxy.RoundRobinProxySwitcher("socks5://s5.citadel.cc:61080")
	if err != nil {
		log.Fatal(err)
	}
	c.SetProxyFunc(rp)

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

func parse() (adverts []*entities.Adverts) {
	cfg := cfg2.NewRequestCfg()
	c := GetInstance()
	dc := GetInstance()
	InitLogs(c, dc)

	// create the redis storage
	storage := &redisstorage.Storage{
		Address:  "0.0.0.0:8081",
		Password: "",
		DB:       0,
		Prefix:   "olx_test",
	}

	// add storage to the collector
	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	// delete previous data from storage
	if err := storage.Clear(); err != nil {
		log.Fatal(err)
	}

	// close redis client
	defer storage.Client.Close()

	// create a new request queue with redis storage backend
	q, _ := queue.New(2, storage)

	// visit each advert
	c.OnHTML(".wrap .offer .space a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = strings.Split(link, "#")[0] // remove anchor tag
		dc.Visit(link)
	})

	// follow available pagination link
	c.OnHTML(".next a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	// parse information from details page
	dc.OnHTML(".offerbody", func(e *colly.HTMLElement) {
		var (
			images  []string
			details []*entities.DetailsItem
		)

		title := e.ChildText(".offer-titlebox h1")
		detailsPlace := e.ChildText(".offer-titlebox__details .show-map-link strong")
		price := e.ChildText(".price-label .xxxx-large")
		t, id := helpers.ParseMeta(helpers.PrettifyString(e.ChildText(".offer-titlebox__details em")))
		text := e.ChildText("#textContent p")

		e.ForEach(".img-item", func(_ int, el *colly.HTMLElement) {
			image := el.ChildAttr(".photo-glow img", "src")
			images = append(images, image)
		})
		e.ForEach(".descriptioncontent .details .item", func(_ int, el *colly.HTMLElement) {
			name := el.ChildText("th")
			value := helpers.PrettifyString(el.ChildText(".value a"))
			if "" == value {
				value = el.ChildText(".value strong")
			}
			details = append(details, &entities.DetailsItem{name, value})
		})

		adverts = append(adverts, &entities.Adverts{
			Id:        id,
			URL:       e.Request.URL.String(),
			Title:     title,
			Place:     detailsPlace,
			CreatedAt: t,
			Details:   details,
			Images:    images,
			Text:      strings.TrimSpace(text),
			Price:     price,
		})
	})

	// ?500? - num pages to visit, there is no need in more than current value
	for i := 1; i <= 10; i++ { // TODO
		q.AddURL(fmt.Sprintf("%s/%s/%s/%s?page=%d", cfg.BASEURL, cfg.CATEGORY, cfg.SUBCATEGORY, cfg.REGION, i))
	}
	q.Run(c)
	c.Wait()
	dc.Wait()
	return
}
