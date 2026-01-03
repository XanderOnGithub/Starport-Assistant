package main

import (
	"log"
	"os"
	"os/signal"
	"starport-assistant/modules/bot"
	_ "starport-assistant/modules/commands/list"
	"starport-assistant/modules/commands/list/lobby"
	"starport-assistant/modules/watcher"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	guildID := os.Getenv("DISCORD_GUILD_ID")

	session := bot.Boot(os.Getenv("DISCORD_TOKEN"), guildID)

	// Trigger watcher on ready
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("🛰️ Starport Systems Online. Initializing Persistence & Watcher...")
		watcher.Start(s)
	})

	lobby.CleanAllLobbies(session, guildID)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("Shutting down Starport Assistant...")
	session.Close()
}
