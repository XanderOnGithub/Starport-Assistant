package bot

import (
	"log"
	"starport-assistant/modules/commands" // "modules/commands" --> To access Register function

	"github.com/bwmarrin/discordgo"
)

// Boot --> Creates session, opens connection, and registers commands
func Boot(token string, guildID string) *discordgo.Session {

	// "discordgo" --> Create new session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Critical Error: %v", err)
	}

	// "discordgo" --> Open connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Connection Error: %v", err)
	}

	// "modules/commands" --> Push all commands to Discord
	commands.Register(dg, guildID)

	log.Println("🛰️ Starport Assistant: Systems Nominal.")
	return dg
}
