package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Struct --> Defines what a modular command looks like
type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Slice --> The central list of all commands
var List []Command

// Add --> Function used by modules to "sign up" for the list
func Add(cmd Command) {
	List = append(List, cmd)
}

// Register --> Logic to send the list to Discord and setup the listener
func Register(s *discordgo.Session, guildID string) {

	// Loop --> Prepare definitions for Discord API
	definitions := make([]*discordgo.ApplicationCommand, len(List))
	for i, cmd := range List {
		definitions[i] = cmd.Definition
	}

	// API --> Bulk overwrite commands for instant testing in your Guild
	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, definitions)
	if err != nil {
		log.Printf("❌ Failed to register commands: %v", err)
	}

	// Handler --> Point Discord to our Switchboard function (No more nested mess!)
	s.AddHandler(handleInteraction)
}

// handleInteraction --> The "Switchboard" that routes clicks to the right command logic
func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Filter --> Ignore anything that isn't a Slash Command (Early Return)
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Data --> Grab the name of the command the user actually typed
	commandName := i.ApplicationCommandData().Name

	// Loop --> Find the matching command in our List
	for _, cmd := range List {
		if cmd.Definition.Name == commandName {
			// Logic --> Run the specific handler for this command
			cmd.Handler(s, i)
			return
		}
	}
}
