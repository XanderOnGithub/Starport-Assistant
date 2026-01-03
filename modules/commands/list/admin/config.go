package admin

import (
	"encoding/json"
	"fmt"
	"starport-assistant/modules/commands"
	"starport-assistant/modules/storage"
	"starport-assistant/modules/watcher"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var gameChoices = []*discordgo.ApplicationCommandOptionChoice{
	{Name: "ARC Raiders", Value: "arcraiders"},
	{Name: "Overwatch 2", Value: "overwatch2"},
}

func init() {
	var adminPerms int64 = discordgo.PermissionAdministrator

	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:                     "config",
			Description:              "Manage Starport Assistant routing and raw data",
			DefaultMemberPermissions: &adminPerms,
			Options: []*discordgo.ApplicationCommandOption{
				{Name: "watcher", Description: "Toggle watcher state", Type: discordgo.ApplicationCommandOptionSubCommand, Options: []*discordgo.ApplicationCommandOption{{Type: discordgo.ApplicationCommandOptionBoolean, Name: "enabled", Description: "Set state", Required: true}}},
				{Name: "add_game", Description: "Route a game to a channel", Type: discordgo.ApplicationCommandOptionSubCommand, Options: []*discordgo.ApplicationCommandOption{{Type: discordgo.ApplicationCommandOptionString, Name: "game", Description: "Game", Choices: gameChoices, Required: true}, {Type: discordgo.ApplicationCommandOptionChannel, Name: "channel", Description: "Channel", Required: true}}},
				{Name: "remove_game", Description: "Stop tracking a game", Type: discordgo.ApplicationCommandOptionSubCommand, Options: []*discordgo.ApplicationCommandOption{{Type: discordgo.ApplicationCommandOptionString, Name: "game", Description: "Game", Choices: gameChoices, Required: true}}},
				{Name: "check", Description: "Force manual scan", Type: discordgo.ApplicationCommandOptionSubCommand},
				{Name: "data", Description: "Show raw storage object", Type: discordgo.ApplicationCommandOptionSubCommand},
			},
		},
		Handler: handleConfigCommand,
	})
}

func handleConfigCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := storage.LoadData()
	cmd := i.ApplicationCommandData().Options[0]

	switch cmd.Name {
	case "data":
		raw, _ := json.MarshalIndent(data, "", "  ")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("📂 **Internal Storage State:**\n```json\n%s\n```", string(raw)),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return

	case "check":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: "📡 **Manual Scan Initialized.**", Flags: discordgo.MessageFlagsEphemeral},
		})
		go watcher.PerformScan(s)
		return

	case "watcher":
		data.WatchEnabled = cmd.Options[0].BoolValue()
	case "add_game":
		game := cmd.Options[0].StringValue()
		data.TrackedGames[game] = cmd.Options[1].ChannelValue(nil).ID
	case "remove_game":
		delete(data.TrackedGames, cmd.Options[0].StringValue())
	}

	storage.SaveData(data)

	// Build Embed
	status := "🔴 DISABLED"
	if data.WatchEnabled {
		status = "🟢 ENABLED"
	}

	routes := []string{}
	for g, c := range data.TrackedGames {
		routes = append(routes, fmt.Sprintf("**%s** ➔ <#%s>", g, c))
	}
	routeText := strings.Join(routes, "\n")
	if routeText == "" {
		routeText = "None"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{{
				Title: "⚙️ CONFIGURATION",
				Color: 0x18176c,
				Fields: []*discordgo.MessageEmbedField{
					{Name: "Watcher", Value: status, Inline: true},
					{Name: "Routes", Value: routeText, Inline: false},
				},
			}},
		},
	})
}
