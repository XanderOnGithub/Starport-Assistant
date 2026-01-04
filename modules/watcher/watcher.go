package watcher

import (
	"fmt"
	"starport-assistant/modules/commands" // Access the shared logger
	"starport-assistant/modules/storage"
	"time"

	"starport-assistant/modules/commands/list/patch/arcraiders"
	"starport-assistant/modules/commands/list/patch/overwatch"

	"github.com/bwmarrin/discordgo"
)

func Start(s *discordgo.Session) {
	ticker := time.NewTicker(time.Hour)
	commands.LogToFile("🛰️ [Watcher] Background service initialized (4h interval).")

	go func() {
		for range ticker.C {
			PerformScan(s)
		}
	}()
}

func PerformScan(s *discordgo.Session) {
	commands.LogToFile("🔍 [Watcher] Beginning scheduled sector scan...")
	data := storage.LoadData()

	if !data.WatchEnabled {
		commands.LogToFile("⚠️ [Watcher] Scan aborted: Watcher is disabled in config.")
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
				commands.LogToFile(fmt.Sprintf("☄️ [Watcher] New ARC Raiders update detected: %s", patch.Title))
				data.LastArcPatch = patch.URL
				storage.SaveData(data)
				s.ChannelMessageSendEmbed(channelID, arcraiders.NewArcEmbed(patch))
			}

		case "overwatch2":
			patch, err := overwatch.GetLatestOWPatch()
			if err == nil && patch.Title != "" && patch.Title != data.LastOWPatch {
				commands.LogToFile(fmt.Sprintf("🚀 [Watcher] New Overwatch 2 patch detected: %s", patch.Title))
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
	commands.LogToFile("✅ [Watcher] Sector scan complete.")
}
