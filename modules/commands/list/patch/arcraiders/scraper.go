package arcraiders

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ArcPatch struct {
	Title   string
	URL     string
	Summary string
	Image   string
}

func GetLatestArcPatch() (ArcPatch, error) {
	scraper := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	data := ArcPatch{}

	// 1. Visit the News Index to find the latest Update/Hotfix
	scraper.OnHTML("a.news-article-card_container__xsniv", func(e *colly.HTMLElement) {
		if data.URL != "" {
			return // Already found the latest
		}

		title := strings.TrimSpace(e.ChildText(".news-article-card_title__7LpPs"))

		// Only grab the post if it's an Update or Hotfix
		if strings.Contains(title, "Update") || strings.Contains(title, "Hotfix") {
			data.Title = title
			data.URL = e.Request.AbsoluteURL(e.Attr("href"))
			data.Image = e.ChildAttr("img", "src")
		}
	})

	err := scraper.Visit("https://arcraiders.com/news")
	if err != nil {
		return data, err
	}

	if data.URL == "" {
		return data, fmt.Errorf("could not find a recent Arc Raiders update")
	}

	// 2. Visit the actual Patch Page and build a readable summary
	detailScraper := colly.NewCollector()

	detailScraper.OnHTML(".article_article__Do3j2", func(e *colly.HTMLElement) {
		var lines []string

		// Iterate through children to preserve the flow of the article
		e.ForEach("p, li, ul", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if text == "" {
				return
			}

			switch el.Name {
			case "p":
				// Bold the "Raiders!" intro or game version mentions
				if strings.HasPrefix(text, "Raiders!") {
					lines = append(lines, "**"+text+"**")
				} else if strings.Contains(text, "Update") || strings.Contains(text, "//") {
					lines = append(lines, text)
				} else {
					lines = append(lines, text)
				}
			case "ul":
				// When we hit a list, add a header and a small gap
				lines = append(lines, "\n**Changes**")
			case "li":
				// Standard Discord bullet point
				lines = append(lines, "• "+text)
			}
		})

		// Join everything with newlines for that "spaced out" look
		data.Summary = strings.Join(lines, "\n")
	})

	err = detailScraper.Visit(data.URL)

	// Safety check for Discord character limits
	if len(data.Summary) > 3000 {
		data.Summary = data.Summary[:3000] + "..."
	}

	return data, err
}
