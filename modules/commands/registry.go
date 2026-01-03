package commands

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var List []Command
var ComponentHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

// LogToFile appends a formatted string to starport.log
func LogToFile(message string) {
	f, err := os.OpenFile("starport.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("❌ Could not open log file:", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	fmt.Print(logEntry)     // Print to Terminal
	f.WriteString(logEntry) // Save to File
}

func Add(cmd Command) {
	List = append(List, cmd)
	LogToFile(fmt.Sprintf("📦 [Module] Loaded command: /%s", cmd.Definition.Name))
}

func AddComponentHandler(customID string, h func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	ComponentHandlers[customID] = h
}

func Register(s *discordgo.Session, guildID string) {
	LogToFile("📡 [Registry] Synchronizing with Discord...")

	definitions := make([]*discordgo.ApplicationCommand, len(List))
	for i, cmd := range List {
		definitions[i] = cmd.Definition
	}

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, definitions)
	if err != nil {
		LogToFile(fmt.Sprintf("❌ [Registry] Failed to register commands: %v", err))
	} else {
		LogToFile(fmt.Sprintf("✅ [Registry] Successfully synced %d command(s) to Discord.", len(List)))
	}

	s.AddHandler(handleInteraction)
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.Member.User.Username

	// 1. Handle Slash Commands
	if i.Type == discordgo.InteractionApplicationCommand {
		commandName := i.ApplicationCommandData().Name
		LogToFile(fmt.Sprintf("💬 [Command] /%s triggered by @%s", commandName, user))

		for _, cmd := range List {
			if cmd.Definition.Name == commandName {
				cmd.Handler(s, i)
				return
			}
		}
	}

	// 2. Handle Button Clicks
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID
		LogToFile(fmt.Sprintf("🔘 [Component] %s clicked by @%s", customID, user))

		if handler, ok := ComponentHandlers[customID]; ok {
			handler(s, i)
		}
	}
}
