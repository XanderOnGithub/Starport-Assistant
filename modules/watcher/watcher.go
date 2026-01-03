package watcher

import (
	"starport-assistant/modules/storage"
	"time"

	"starport-assistant/modules/commands/list/patch/arcraiders"
	"starport-assistant/modules/commands/list/patch/overwatch"

	"github.com/bwmarrin/discordgo"
)

func Start(s *discordgo.Session) {
	ticker := time.NewTicker(4 * time.Hour)
	go func() {
		for range ticker.C {
			PerformScan(s)
		}
	}()
}

func PerformScan(s *discordgo.Session) {
	data := storage.LoadData()
	if !data.WatchEnabled {
		return
	}

	for game, channelID := range data.TrackedGames {
		if channelID == "" {
			continue
		}

		switch game {
		case "arcraiders":
			patch, err := arcraiders.GetLatestArcPatch()
			if err == nil && patch.URL != "" && patch.URL != data.LastArcPatch {
				data.LastArcPatch = patch.URL
				storage.SaveData(data)
				s.ChannelMessageSendEmbed(channelID, arcraiders.NewArcEmbed(patch))
			}

		case "overwatch2":
			patch, err := overwatch.GetLatestOWPatch()
			// Compare Title because URL is static
			if err == nil && patch.Title != "" && patch.Title != data.LastOWPatch {
				data.LastOWPatch = patch.Title
				storage.SaveData(data)

				s.ChannelMessageSend(channelID, "🚀 **NEW OVERWATCH 2 PATCH: "+patch.Title+"**\n"+patch.URL)
				for _, hero := range patch.Updates {
					embed := overwatch.NewHeroEmbed(hero, patch.Title)
					s.ChannelMessageSendEmbed(channelID, embed)
				}
			}
		}
	}
}
