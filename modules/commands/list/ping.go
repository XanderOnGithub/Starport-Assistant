package list

import (
	"starport-assistant/modules/commands" // "modules/commands" --> To use the Add function

	"github.com/bwmarrin/discordgo"
)

// init --> Runs automatically when the package is imported by main.go
func init() {

	// "modules/commands" --> Self-register to the central list
	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Check Starport latency",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "🏓 Pong!",
				},
			})
		},
	})
}
