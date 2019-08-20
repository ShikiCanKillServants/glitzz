package pornhub

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"regexp"
	"time"
)

const (
	pornhubURL  = "https://www.pornhub.com"
	userAgent   = "glitzz, an irc bot"
	httpTimeout = 4
)

// Pornhub structure, populated and returned by calling Init().
type Pornhub struct {
	doc *goquery.Document

	Title      string
	URL        string
	Uploader   string
	Categories []string
	Views      int
	Votes      votes
	Comments   []comment
}

// Result is used to return invididual search results.
type Result struct {
	URL   string
	Title string
}

func fetch(url string) (doc *goquery.Document, err error) {
	client := &http.Client{Timeout: httpTimeout * time.Second}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("User-Agent", userAgent)

	response, err := client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusMovedPermanently {
		return doc, fmt.Errorf("Something went wrong: %d code", response.StatusCode)
	}
	defer response.Body.Close()

	doc, err = goquery.NewDocumentFromResponse(response)
	if err != nil {
		return
	}
	return
}

// Init fetches and load the page.
func Init(url string) (p *Pornhub, err error) {
	p = &Pornhub{}

	p.doc, err = fetch(url)
	if err != nil {
		return
	}

	p.setInfos()
	p.setSocial()

	return
}

// InitRandom fetches and load a random page.
func InitRandom(gay bool) (p *Pornhub, err error) {
	var gayPath string
	if gay {
		gayPath = "/gay"
	}

	return Init(pornhubURL + gayPath + "/random")
}

// Search returns a list of matching videos
func Search(query string, gay bool) (results []Result, err error) {
	var gayPath string
	if gay {
		gayPath = "/gay"
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return
	}

	doc, err := fetch(
		pornhubURL +
			gayPath +
			"/video/search?search=" +
			reg.ReplaceAllString(query, "+"))

	if err != nil {
		return
	}

	doc.Find("#videoSearchResult .title a[href]").Each(func(i int, s *goquery.Selection) {
		var node Result

		URL, exists := s.Attr("href")
		node.URL = pornhubURL + URL
		node.Title = s.Text()
		if exists {
			results = append(results, node)
		}
	})

	return
}
