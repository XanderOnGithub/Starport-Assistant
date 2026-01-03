package overwatch

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func NewHeroEmbed(hero HeroUpdate, patchTitle string) *discordgo.MessageEmbed {
	// Format the changes with an extra leading newline if the first line is a title
	description := hero.Changes
	if !strings.HasPrefix(description, "\n") {
		description = "\n" + description
	}

	return &discordgo.MessageEmbed{
		Title:       "🛡️ HERO CHANGE: " + strings.ToUpper(hero.Name),
		Description: description,
		Color:       0xF99E1A, // OW Orange
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: hero.IconURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Source: " + patchTitle,
		},
	}
}
