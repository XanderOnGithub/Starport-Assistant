package main

import (
	"log"
	"os"
	"os/signal"
	"starport-assistant/modules/bot"
	"starport-assistant/modules/commands"
	_ "starport-assistant/modules/commands/list"
	"starport-assistant/modules/commands/list/lobby"
	"starport-assistant/modules/watcher"
	"syscall"
	"time" // Added for sleep

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	guildID := os.Getenv("DISCORD_GUILD_ID")
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is empty")
	}

	session := bot.Boot(token)
	if session == nil {
		log.Fatal("bot.Boot returned nil session")
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("🛰️ Starport Systems Online.")

		if err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name:  "custom",
					Type:  discordgo.ActivityTypeCustom,
					State: "🛰️ Scanning for New Life...",
				},
			},
			Status: "online",
		}); err != nil {
			log.Println("failed to set status:", err)
		}

		// SAFETY PATCH: Run initialization in the background
		go func() {
			// Give Discord time to populate its internal cache/state
			time.Sleep(2 * time.Second)

			log.Println("🧹 Cleaning up old lobbies...")
			lobby.CleanAllLobbies(s, guildID)

			log.Println("🛰️ Starting Watcher...")
			watcher.Start(s)

			log.Println("✅ Starport Assistant: Systems Nominal.")
		}()
	})

	if err := session.Open(); err != nil {
		log.Fatalf("error opening Discord session: %v", err)
	}

	commands.Register(session, guildID)

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	log.Println("📡 Shutting down Starport Assistant...")
}
