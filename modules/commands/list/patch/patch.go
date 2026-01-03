package patch

import (
	"starport-assistant/modules/commands"
	"starport-assistant/modules/commands/list/patch/arcraiders" // Import Arc Raiders package
	"starport-assistant/modules/commands/list/patch/overwatch"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "patch",
			Description: "Get latest patch notes for a specific game",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "game",
					Description: "Choose the game",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Overwatch 2", Value: "ow2"},
						{Name: "Arc Raiders", Value: "arc"}, // Added Arc Raiders Choice
					},
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			gameChoice := options[0].StringValue()

			if gameChoice == "ow2" {
				handleOverwatch(s, i)
			} else if gameChoice == "arc" {
				handleArcRaiders(s, i) // Added Arc Raiders Router
			}
		},
	})
}

// --- OVERWATCH HANDLER ---
func handleOverwatch(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data, err := overwatch.GetLatestOWPatch()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Error connecting to Blizzard."},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "🚀 **" + data.Title + "**\nFull notes: " + data.URL,
		},
	})

	for _, hero := range data.Updates {
		embed := overwatch.NewHeroEmbed(hero, data.Title)

		// Fixed pointer issue by assigning to a variable first
		params := &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{embed},
		}

		s.FollowupMessageCreate(i.Interaction, true, params)
	}
}

// --- ARC RAIDERS HANDLER ---
func handleArcRaiders(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data, err := arcraiders.GetLatestArcPatch()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "❌ Error fetching Arc Raiders news."},
		})
		return
	}

	embed := arcraiders.NewArcEmbed(data)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
