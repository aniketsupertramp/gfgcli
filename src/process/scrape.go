package process

import (
	"net/http"
	"log"
	"errors"
	"github.com/gocolly/colly"
)

func GetDocumentStatus(url string) int {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("gfgCLI:: Error geting link response: link: %s, error: %s", url, err.Error())
	}
	defer res.Body.Close()
	return res.StatusCode
}

func parsDocument(url, locator string, callbackFunction func(element *colly.HTMLElement)) error {
	if status := GetDocumentStatus(url); status != http.StatusOK {
		return errors.New(ErrNon200)
	}
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: geeksforgeeks.org,www.geeksforgeeks.org
		colly.AllowedDomains("geeksforgeeks.org", "www.geeksforgeeks.org"),
	)
	c.OnHTML(locator, callbackFunction)
	// Start scraping
	c.Visit(url)
	return nil
}
