package arcraiders

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func NewArcEmbed(data ArcPatch) *discordgo.MessageEmbed {
	// We'll format the summary here to add extra polish
	formattedDescription := fmt.Sprintf("%s\n\n**Stay warm. Stay alive.**\n\n*// The ARC Raiders Team*", data.Summary)

	return &discordgo.MessageEmbed{
		Title:       "☄️ ARC RAIDERS: " + data.Title,
		URL:         data.URL,
		Description: formattedDescription,
		Color:       0xff6600, // Arc Raiders Orange/Rust color
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://arcraiders.com/favicon.ico",
		},
		Image: &discordgo.MessageEmbedImage{
			URL: data.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Embark Studios • Patch News",
		},
	}
}
