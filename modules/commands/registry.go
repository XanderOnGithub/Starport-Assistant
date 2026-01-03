package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// List of slash commands
var List []Command

// Map of button/component handlers (Initialized to prevent panic)
var ComponentHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))

func Add(cmd Command) {
	List = append(List, cmd)
	fmt.Printf("📦 [Module] Loaded command: /%s\n", cmd.Definition.Name)
}

// AddComponentHandler allows modules to register button logic
func AddComponentHandler(customID string, h func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	ComponentHandlers[customID] = h
}

func Register(s *discordgo.Session, guildID string) {
	fmt.Println("📡 [Registry] Synchronizing with Discord...")

	definitions := make([]*discordgo.ApplicationCommand, len(List))
	for i, cmd := range List {
		definitions[i] = cmd.Definition
	}

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, definitions)
	if err != nil {
		log.Printf("❌ [Registry] Failed to register commands: %v", err)
	} else {
		fmt.Printf("✅ [Registry] Successfully synced %d command(s) to Discord.\n", len(List))
	}

	s.AddHandler(handleInteraction)
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 1. Handle Slash Commands
	if i.Type == discordgo.InteractionApplicationCommand {
		commandName := i.ApplicationCommandData().Name
		for _, cmd := range List {
			if cmd.Definition.Name == commandName {
				cmd.Handler(s, i)
				return
			}
		}
	}

	// 2. Handle Button Clicks (Components)
	if i.Type == discordgo.InteractionMessageComponent {
		customID := i.MessageComponentData().CustomID
		if handler, ok := ComponentHandlers[customID]; ok {
			handler(s, i)
		}
	}
}
