package list

import (
	"starport-assistant/modules/commands"
	"starport-assistant/modules/utils"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "patch",
			Description: "Get latest Overwatch hero balance changes",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// 1. Scrape the patch data from Blizzard
			data, err := utils.GetLatestOWPatch()
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{Content: "❌ Error: Could not reach Blizzard HQ."},
				})
				return
			}

			// 2. Check if hero changes were actually found
			if len(data.Updates) == 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "🛰️ No hero balance changes found in the latest patch.",
					},
				})
				return
			}

			// 3. FIRST MESSAGE: Send the Intro with the link
			// We use InteractionRespond here to acknowledge the slash command immediately
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🚀 **" + data.Title + "**\nFull notes available at: " + data.URL,
				},
			})
			if err != nil {
				return
			}

			// 4. SUBSEQUENT MESSAGES: Send each hero as a follow-up
			for _, hero := range data.Updates {
				embed := utils.NewHeroEmbed(hero, data.Title)

				// FollowupMessageCreate sends additional messages tied to the original interaction
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{embed},
				})
			}
		},
	})
}
