package main

import (
	"fmt"
	"github.com/gocolly/colly"

	"strings"
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

	// Instantiate default collector
	collector := colly.NewCollector(
		colly.MaxDepth(2),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"),
		colly.Async(true),
		colly.CacheDir("./.cache"),
	)

	// Instantiate collector for details page
	detailCollector := collector.Clone()

	// visit each advert
	collector.OnHTML(".wrap .offer .space a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		detailCollector.Visit(link)
	})

	// follow available pagination link
	collector.OnHTML(".next a[href]", func(e *colly.HTMLElement) {
		collector.Visit(e.Attr("href"))
	})

	// parse information from details page
	detailCollector.OnHTML(".offercontent", func(e *colly.HTMLElement) {
		title := e.ChildText(".offer-titlebox h1")
		detailsPlace := e.ChildText(".offer-titlebox__details .show-map-link strong")
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

	collector.Visit(fmt.Sprintf("%s/%s/%s/%s", cfg.BASEURL, cfg.CATEGORY, cfg.SUBCATEGORY, cfg.REGION)) // start url

	collector.Wait()

	// debug
	fmt.Println(len(adverts))
	for _, advert := range adverts {
		fmt.Printf("%s\n%s\n%s\n%s\n%s\n%v\n",
			advert.Title, advert.URL, advert.Place, advert.Meta, advert.Text, advert.Images)
		fmt.Println("===")
	}
}
