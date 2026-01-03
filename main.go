package main

import (
	"os"
	"os/signal"
	"starport-assistant/modules/bot"
	_ "starport-assistant/modules/commands/list" // "_" --> Runs all init() functions in that folder automatically
	"syscall"

	"github.com/joho/godotenv"
)

func main() {

	// ".env" --> Load environment variables
	godotenv.Load()

	// "modules/bot" --> Initialize and Login (Everything handled in one go)
	session := bot.Boot(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_GUILD_ID"))

	// Action --> Wait for CTRL-C to shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	// Action --> Clean Shutdown
	session.Close()
}
