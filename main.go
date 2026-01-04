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

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment
	_ = godotenv.Load()

	guildID := os.Getenv("DISCORD_GUILD_ID")
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is empty")
	}

	// 2. Initialize the Session (do NOT open yet)
	session := bot.Boot(token)
	if session == nil {
		log.Fatal("bot.Boot returned nil session")
	}

	// 3. Add handlers BEFORE opening the connection
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("🛰️ Starport Systems Online. Initializing Persistence & Watcher...")

		// Set Presence/Status and check error
		if err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name:  "custom",
					Type:  discordgo.ActivityTypeCustom,
					State: "🛰️ Scanning for New Life...",
				},
			},
			Status: "online",
			AFK:    false,
		}); err != nil {
			log.Println("failed to set status:", err)
		}

		// Run cleanup now that we are actually connected
		log.Println("🧹 Cleaning up old lobbies...")
		lobby.CleanAllLobbies(s, guildID)

		// Start the patch watcher
		watcher.Start(s)

		log.Println("✅ Starport Assistant: Systems Nominal.")
	})

	// 4. Open the session
	if err := session.Open(); err != nil {
		log.Fatalf("error opening Discord session: %v", err)
	}
	// Register commands after opening (if your Register uses the open session/REST)
	commands.Register(session, guildID)

	// Ensure session is closed on exit
	defer session.Close()

	// 5. Wait for Termination Signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	// 6. Clean Shutdown
	log.Println("📡 Shutting down Starport Assistant...")
}
