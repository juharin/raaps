package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Course stores information about a coursera course
type Car struct {
	Brand        string
	Model        string
	URL          string
	Year         string
	Mileage      string
	Type         string
	FullPrice    string
	MonthlyPrice string
	MonthlyLease string
}

func main() {
	fName := "cars.json"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains("aldcarmarket.fi"),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./aldcarmarket_cache"),
	)

	// Create another collector to scrape car details
	detailCollector := c.Clone()

	cars := make([]Car, 0, 200)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// On every <a> element with block-link class call callback
	c.OnHTML(`a.block-link`, func(e *colly.HTMLElement) {
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		log.Println("Link: ", courseURL)
		if strings.Index(courseURL, "aldcarmarket.fi/fi/autot") != -1 {
			detailCollector.Visit(courseURL)
		}
	})

	// Extract details of the car
	detailCollector.OnHTML(`div.car-page-header`, func(e *colly.HTMLElement) {
		log.Println("Car found", e.Request.URL)
		brand := e.ChildText(".brand")
		if brand == "" {
			log.Println("No title found", e.Request.URL)
		}
		car := Car{
			Brand: brand,
			//Model: e.ChildText(""),
			URL:          e.Request.URL.String(),
			Year:         e.ChildText("span.year"),
			Mileage:      e.ChildText("span.mileage"),
			Type:         e.ChildText("div.typename"),
			FullPrice:    e.ChildText("div.car-page-price"),
			MonthlyPrice: e.ChildText("a.price-montly"),
			MonthlyLease: e.ChildText("a.price-leasing"),
		}
		cars = append(cars, car)
	})
	log.Println("Starting")
	// Start scraping on
	c.Visit("https://www.aldcarmarket.fi/fi/vaihtoautot/kaikki/?kayttovoima=hybridi,sahko&leasattavat=1")
	log.Println("Done?")
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(cars)
}
