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
	// 1. Load Environment
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, relying on system environment variables")
	}

	guildID := os.Getenv("DISCORD_GUILD_ID")
	token := os.Getenv("DISCORD_TOKEN")

	// 2. Initialize the Session
	session := bot.Boot(token, guildID)

	// 3. Define the Ready Handler
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("🛰️ Starport Systems Online. Initializing Persistence & Watcher...")

		// --- NEW: Set Presence/Status ---
		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "with Space & Time",
					Type: discordgo.ActivityTypeGame,
				},
			},
			Status: "online",
			AFK:    false,
		})

		// Run cleanup now that we are actually connected
		log.Println("🧹 Cleaning up old lobbies...")
		lobby.CleanAllLobbies(s, guildID)

		// Start the patch watcher
		watcher.Start(s)

		log.Println("✅ Starport Assistant: Systems Nominal.")
	})

	// 4. Wait for Termination Signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// This blocks the main thread so the bot stays alive
	<-stop

	// 5. Clean Shutdown
	log.Println("📡 Shutting down Starport Assistant...")
	session.Close()
}
