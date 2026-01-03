package overwatch

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type HeroUpdate struct {
	Name    string
	IconURL string
	Changes string
}

type PatchData struct {
	Title   string
	URL     string
	Updates []HeroUpdate
}

func GetLatestOWPatch() (PatchData, error) {
	scraper := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	data := PatchData{
		URL: "https://overwatch.blizzard.com/en-us/news/patch-notes/",
	}

	// Regex to find and bold values (e.g., 50%, 12s, 10,000, 4.5)
	valueRegex := regexp.MustCompile(`(\d+,\d+|\d+(\.\d+)?%|\d+s|\d+(\.\d+)?)`)

	// Use :first-of-type to grab only the first patch block (the live one)
	scraper.OnHTML(".PatchNotes-patch:first-of-type", func(patch *colly.HTMLElement) {

		// 1. Get Title and Date
		data.Title = patch.ChildText(".PatchNotes-patchTitle")

		// 2. Find every Hero card within THIS patch only
		patch.ForEach(".PatchNotesHeroUpdate", func(_ int, hero *colly.HTMLElement) {
			name := hero.ChildText(".PatchNotesHeroUpdate-name")
			icon := hero.ChildAttr("img.PatchNotesHeroUpdate-icon", "src")

			var changes []string

			// We iterate through the top-level list items (Abilities or standalone changes)
			hero.ForEach(".PatchNotesHeroUpdate-generalUpdates > ul > li", func(_ int, li *colly.HTMLElement) {

				// Check if this item has a sub-list (nested changes)
				subList := li.DOM.Find("ul")
				if subList.Length() == 0 {
					// No sub-list: It's a direct bullet point
					text := strings.TrimSpace(li.Text)
					if text != "" {
						styledText := valueRegex.ReplaceAllString(text, "**$1**")
						changes = append(changes, "• "+styledText)
					}
				} else {
					// Sub-list exists: This LI is a category header (e.g., "Primal Punch - Power")
					// We extract the header text while ignoring the sub-list contents to avoid duplication
					clone := li.DOM.Clone()
					clone.Find("ul").Remove()
					headerText := strings.TrimSpace(clone.Text())

					if headerText != "" {
						// Add a visual gap before categories to make them pop
						if len(changes) > 0 {
							changes = append(changes, "")
						}
						changes = append(changes, "**"+headerText+"**")
					}

					// Now process the specific changes inside the sub-list
					li.ForEach("ul > li", func(_ int, subLi *colly.HTMLElement) {
						subText := strings.TrimSpace(subLi.Text)
						if subText != "" {
							// Bold the numbers/values in the change description
							styledSub := valueRegex.ReplaceAllString(subText, "**$1**")
							changes = append(changes, "  └ "+styledSub)
						}
					})
				}
			})

			if name != "" && len(changes) > 0 {
				data.Updates = append(data.Updates, HeroUpdate{
					Name:    strings.TrimSpace(name),
					IconURL: icon,
					Changes: strings.Join(changes, "\n"),
				})
			}
		})
	})

	err := scraper.Visit(data.URL)
	return data, err
}
