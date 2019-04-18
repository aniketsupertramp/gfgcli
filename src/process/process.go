package process

import (
	"github.com/manifoldco/promptui"
	"strings"
	"strconv"
	"fmt"
	"os"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/gocolly/colly"
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
	parsDocument(InterviewUrl, "div[class=entry-content]>ul:nth-of-type(1) >li", func(element *colly.HTMLElement) {
		element.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			nameSlice := strings.Split(e.Text, "[")
			name := strings.TrimRight(nameSlice[0], " ")
			noOfArticles, _ := strconv.Atoi(nameSlice[1][0 : len(nameSlice[1])-1])
			company := Company{Name: name, Href: link, NoOfArticles: noOfArticles}
			Companies = append(Companies, company)

		})
	})
	// stop the spinner
	wow.Stop()
}

func loadArticles(company Company, wow *wow.Wow) {
	wow.Start()
	for pageNo := 1; pageNo < MaxPageSize; pageNo++ {

		url := company.Href + "page/" + strconv.Itoa(pageNo)
		err := parsDocument(url, "article .entry-title", func(element *colly.HTMLElement) {
			element.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
				article := Article{Title: strings.TrimSpace(e.Text), Href: e.Attr("href")}
				Articles = append(Articles, article)
			})
		})
		if err != nil {
			break
		}
	}
	wow.Stop()
}

func displayArticle(article Article) {
	locator := "div[class=entry-content]>p,div[class=entry-content]>ol,div[class=entry-content]>ul,div[class=entry-content]>pre,div[class=entry-content]>blockquote>p,div[class=entry-content] .code-container"
	_ = parsDocument(article.Href, locator, func(element *colly.HTMLElement) {
		fmt.Println(element.Text)
	})

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
