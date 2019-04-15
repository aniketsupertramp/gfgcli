package process

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"log"
	"errors"
	"fmt"
)

func parseDocument(url string) (*goquery.Document, error) {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("gfgCLI:: Error loading geting link response: link: %s, error: %s", url, err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("gfgCLI::status code error: %d %s", res.StatusCode, res.Status)
		return nil, errors.New(ErrNon200)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("gfgCLI::load error: " + err.Error())
	}
	return doc, err
}
