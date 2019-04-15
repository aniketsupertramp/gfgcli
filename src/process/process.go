package process

import (
	"github.com/manifoldco/promptui"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"log"
	"fmt"
	"os"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"net/http"
)

type Company struct {
	Name         string
	Href         string
	NoOfArticles int
}

type Article struct {
	Title string
	Href  string
}

const (
	InterviewUrl = "https://www.geeksforgeeks.org/company-interview-corner/"
	MaxPageSize  = 20
	ErrNon200    = "Non-200 response"
)

var (
	Companies       []Company
	Articles        []Article
	Companytemplate = promptui.SelectTemplates{
		Active:   `🔎  {{ .Name | cyan | bold }}`,
		Inactive: `   {{ .Name | cyan }}`,
		Selected: `{{ "✔" | green | bold }} {{ "Company" | bold }}: {{ .Name | cyan }}`,
	}
	ArticleTemplate = promptui.SelectTemplates{
		Active:   `🔎  {{ .Title | cyan | bold }}`,
		Inactive: `   {{ .Title | cyan }}`,
		Selected: `{{ "✔" | green | bold }} {{ "Article" | bold }}: {{ .Title | cyan }}`,
	}
)

func loadCompanies(wow *wow.Wow) {
	// start the spinner
	wow.Start()
	doc, err := parseDocument(InterviewUrl)
	if err == nil && doc != nil {
		doc.Find("div[class=entry-content]>ul:nth-of-type(1) >li").Each(func(i int, s *goquery.Selection) {
			link, found := s.Find("a[href]").Attr("href")

			if found {
				nameSlice := strings.Split(s.Text(), "[")
				name := strings.TrimRight(nameSlice[0], " ")
				noOfArticles, _ := strconv.Atoi(nameSlice[1][0 : len(nameSlice[1])-1])
				company := Company{Name: name, Href: link, NoOfArticles: noOfArticles}
				Companies = append(Companies, company)
			}
		})
	}
	// stop the spinner
	wow.Stop()
}

func loadArticles(company Company, wow *wow.Wow) {
	wow.Start()
	for pageNo := 1; pageNo < MaxPageSize; pageNo++ {

		link := company.Href + "page/" + strconv.Itoa(pageNo)
		res, err := http.Get(link)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			fmt.Println("gfgCLI::status code: %d error: %s", res.StatusCode, res.Status)
			break
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println("gfgCLI::load ERROR: " + err.Error())
			// Error: gzip: invalid header ,just logging and skipping for now
			continue
		}

		doc.Find("article .entry-title").Each(func(i int, s *goquery.Selection) {

			link, found := s.Find("a[href]").Attr("href")
			if found {
				article := Article{Title: strings.TrimSpace(s.Text()), Href: link}
				Articles = append(Articles, article)
			}
		})
	}
	wow.Stop()
}

func displayArticle(article Article) {
	// Request the HTML page
	doc, _ := parseDocument(article.Href)
	if doc != nil {
		doc.Find("div[class=entry-content]>p").Each(func(i int, s *goquery.Selection) {
			fmt.Println(s.Text())
		})
	}
}

func promptCompaniesList() Company {
	prompt := promptui.Select{
		Label:     "Company",
		Items:     Companies,
		Templates: &Companytemplate,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return Companies[idx]
}

func promptArticles() {
	for true {
		prompt := promptui.Select{
			Label:     "Articles",
			Items:     Articles,
			Templates: &ArticleTemplate,
			Size:      10,
		}

		idx, _, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		displayArticle(Articles[idx])
	}
}

func Init() {
	for true {
		// init the wow spinner
		w := wow.New(os.Stdout, spin.Get(spin.Earth), "Fetching Results")
		// loading companies list
		loadCompanies(w)
		// prompt companies and return selected one
		company := promptCompaniesList()
		// load articles for the selected company
		loadArticles(company, w)
		// prompt all articles for the company
		promptArticles()
	}
}
