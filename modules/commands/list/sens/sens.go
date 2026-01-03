package sens

import (
	"fmt"
	"math"
	"starport-assistant/modules/commands"

	"github.com/bwmarrin/discordgo"
)

func init() {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for slug, info := range GameData {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  info.DisplayName,
			Value: slug,
		})
	}

	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "sens",
			Description: "1:1 Sensitivity Translation",
			Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "from_game", Description: "Starting game", Required: true, Choices: choices},
				{Type: discordgo.ApplicationCommandOptionNumber, Name: "value", Description: "Current sensitivity", Required: true},
				{Type: discordgo.ApplicationCommandOptionString, Name: "to_game", Description: "Target game", Required: true, Choices: choices},
				{Type: discordgo.ApplicationCommandOptionInteger, Name: "from_dpi", Description: "Starting DPI", Required: false},
				{Type: discordgo.ApplicationCommandOptionInteger, Name: "to_dpi", Description: "Target DPI", Required: false},
			},
		},
		Handler: handleSensCommand,
	})
}

func handleSensCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	opts := i.ApplicationCommandData().Options
	fromGame := GameData[opts[0].StringValue()]
	fromSens := opts[1].FloatValue()
	toGame := GameData[opts[2].StringValue()]

	// Handle Optional DPIs
	fromDPI, toDPI := 0, 0
	for _, opt := range opts {
		if opt.Name == "from_dpi" {
			fromDPI = int(opt.IntValue())
		}
		if opt.Name == "to_dpi" {
			toDPI = int(opt.IntValue())
		}
	}

	// 1. Calculate Sensitivity Conversion
	// Logic: (InputSens / FromRatio) * ToRatio
	targetSens := (fromSens / fromGame.Ratio) * toGame.Ratio

	// 2. Adjust for DPI change if both are provided
	if fromDPI > 0 && toDPI > 0 && fromDPI != toDPI {
		targetSens = targetSens * (float64(fromDPI) / float64(toDPI))
	}

	// 3. Formatting
	roundedSens := math.Round(targetSens*1000) / 1000

	embed := &discordgo.MessageEmbed{
		Title:       SensTitle,
		Color:       SensColor,
		Description: fmt.Sprintf("Muscle memory: **%s** ➔ **%s**", fromGame.DisplayName, toGame.DisplayName),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "📥 Input", Value: fmt.Sprintf("`%v`", fromSens), Inline: true},
			{Name: "📤 Result", Value: fmt.Sprintf("`%.3f`", roundedSens), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: SensFooter},
	}

	// 4. Only show Physical Distance if DPI is provided
	if fromDPI > 0 {
		// Using OW base yaw (0.0066) to determine cm/360 for the "Input" side
		cm360 := 360.0 / (fromSens / fromGame.Ratio * 0.0066 * (float64(fromDPI) / 2.54))
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📏 Physical Distance",
			Value:  fmt.Sprintf("**%.2f** cm per 360° (@ %v DPI)", cm360, fromDPI),
			Inline: false,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})
}
