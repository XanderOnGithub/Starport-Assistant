package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Boot -> Creates session but does NOT open it or register commands.
// This allows callers to add handlers before opening the connection.
func Boot(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Critical Error: %v", err)
	}

	return dg
}
