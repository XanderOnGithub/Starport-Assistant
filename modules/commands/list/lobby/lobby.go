package lobby

import (
	"fmt"
	"starport-assistant/modules/commands"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	var maxSlots float64 = 12
	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "lobby",
			Description: "Start a gaming lobby for the squad",
			Options: []*discordgo.ApplicationCommandOption{
				{Type: discordgo.ApplicationCommandOptionString, Name: "game", Description: "Game name", Required: true},
				{Type: discordgo.ApplicationCommandOptionInteger, Name: "slots", Description: "Total slots", Required: true, MaxValue: maxSlots},
				{Type: discordgo.ApplicationCommandOptionString, Name: "time", Description: "When?", Required: true},
				{
					Type:         discordgo.ApplicationCommandOptionChannel,
					Name:         "channel",
					Description:  "Meet-up voice channel",
					Required:     true,
					ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildVoice},
				},
			},
		},
		Handler: handleLobbyCommand,
	})

	commands.AddComponentHandler("lobby_join", HandleButtonClick)
	commands.AddComponentHandler("lobby_leave", HandleButtonClick)
	commands.AddComponentHandler("lobby_start", HandleButtonClick)
	commands.AddComponentHandler("lobby_delete", HandleButtonClick)
}

// Global Cleanup function to run on bot startup
func CleanAllLobbies(s *discordgo.Session, guildID string) {
	if s == nil {
		return
	}

	channels, err := s.GuildChannels(guildID)
	if err != nil {
		fmt.Printf("⚠️ Could not fetch channels: %v\n", err)
		return
	}

	selfID := ""
	// SAFETY PATCH: Check if state is nil before accessing User
	if s.State != nil && s.State.User != nil {
		selfID = s.State.User.ID
	} else {
		// Fallback: Fetch bot info from API if State isn't cached yet
		user, err := s.User("@me")
		if err == nil {
			selfID = user.ID
		}
	}

	for _, ch := range channels {
		if ch.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		msgs, err := s.ChannelMessages(ch.ID, 20, "", "", "")
		if err != nil {
			continue
		}

		for _, m := range msgs {
			// Check if message author exists before checking ID
			if m.Author != nil && m.Author.Bot && (selfID == "" || m.Author.ID == selfID) && len(m.Embeds) > 0 {
				if m.Embeds[0].Footer != nil && strings.Contains(m.Embeds[0].Footer.Text, "Starport Assistant") {
					s.ChannelMessageDelete(ch.ID, m.ID)
				}
			}
		}
	}
}

func handleLobbyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	opts := i.ApplicationCommandData().Options
	gameName := opts[0].StringValue()

	hostName := i.Member.User.GlobalName
	if hostName == "" {
		hostName = i.Member.User.Username
	}

	Mutex.Lock()
	for msgID, lobby := range ActiveLobbies {
		if lobby.HostID == i.Member.User.ID {
			s.ChannelMessageDelete(lobby.ChannelID, msgID)
			delete(ActiveLobbies, msgID)
		}
	}
	Mutex.Unlock()

	var vcID string
	if len(opts) > 3 && opts[3].Value != nil {
		vcID = fmt.Sprintf("%v", opts[3].Value)
	}

	newLobby := &GameLobby{
		ChannelID:       i.ChannelID,
		HostID:          i.Member.User.ID,
		HostDisplayName: hostName,
		GameName:        gameName,
		MaxSlots:        int(opts[1].IntValue()),
		StartTime:       opts[2].StringValue(),
		VoiceChannelID:  vcID,
		Players:         []string{i.Member.User.ID},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{renderLobbyEmbed(newLobby)},
			Components: []discordgo.MessageComponent{renderButtons()},
		},
	})

	resp, _ := s.InteractionResponse(i.Interaction)
	if resp != nil {
		Mutex.Lock()
		ActiveLobbies[resp.ID] = newLobby
		Mutex.Unlock()
	}
}

func HandleButtonClick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	userID := i.Member.User.ID
	messageID := i.Message.ID

	Mutex.Lock()
	l, exists := ActiveLobbies[messageID]
	if !exists {
		Mutex.Unlock()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: MsgLobbyError, Flags: discordgo.MessageFlagsEphemeral},
		})
		return
	}

	switch customID {
	case "lobby_join":
		if len(l.Players) < l.MaxSlots && !contains(l.Players, userID) {
			l.Players = append(l.Players, userID)
		}
	case "lobby_leave":
		if userID != l.HostID {
			l.Players = remove(l.Players, userID)
		}
	case "lobby_start":
		if userID == l.HostID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})

			mentions := ""
			for _, pID := range l.Players {
				mentions += fmt.Sprintf("<@%s> ", pID)
			}

			s.ChannelMessageDelete(i.ChannelID, messageID)
			delete(ActiveLobbies, messageID)
			Mutex.Unlock()

			liveEmbed := &discordgo.MessageEmbed{
				Title:       LiveEmbedTitle,
				Description: fmt.Sprintf(LiveEmbedDesc, strings.ToUpper(l.GameName), l.VoiceChannelID, mentions),
				Color:       LiveEmbedColor,
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: CustomThumbnail},
				Footer:      &discordgo.MessageEmbedFooter{Text: LiveFooterText},
			}

			// 3. Send Ping + Embed
			liveMsg, _ := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
				Content: "Squad: " + mentions + " Your lobby has started!",
				Embed:   liveEmbed,
			})

			if liveMsg != nil {
				time.AfterFunc(1*time.Minute, func() {
					s.ChannelMessageDelete(i.ChannelID, liveMsg.ID)
				})
			}
			return
		}
	case "lobby_delete":
		if userID == l.HostID {
			s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			delete(ActiveLobbies, messageID)
			Mutex.Unlock()
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseDeferredMessageUpdate})
			return
		}
	}
	Mutex.Unlock()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{renderLobbyEmbed(l)},
			Components: []discordgo.MessageComponent{renderButtons()},
		},
	})
}

func renderLobbyEmbed(l *GameLobby) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: EmbedTitlePrefix + strings.ToUpper(l.GameName),
		Description: fmt.Sprintf(
			EmbedDescLine1+"\n\n"+EmbedTimeLabel+" `%s` \n"+EmbedVoiceLabel+" <#%s>\n\n"+
				"%s", l.HostDisplayName, l.StartTime, l.VoiceChannelID, l.RenderSlots()),
		Color:     EmbedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: CustomThumbnail},
		Footer:    &discordgo.MessageEmbedFooter{Text: EmbedFooterText},
	}
}

func renderButtons() discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    BtnJoin,
				Style:    discordgo.SuccessButton,
				CustomID: "lobby_join",
				Emoji:    &discordgo.ComponentEmoji{Name: EmojiJoin},
			},
			discordgo.Button{
				Label:    BtnLeave,
				Style:    discordgo.SecondaryButton,
				CustomID: "lobby_leave",
				Emoji:    &discordgo.ComponentEmoji{Name: EmojiLeave},
			},
			discordgo.Button{
				Label:    BtnStart,
				Style:    discordgo.PrimaryButton,
				CustomID: "lobby_start",
				Emoji:    &discordgo.ComponentEmoji{Name: EmojiStart},
			},
			discordgo.Button{
				Label:    BtnDelete,
				Style:    discordgo.DangerButton,
				CustomID: "lobby_delete",
				Emoji:    &discordgo.ComponentEmoji{Name: EmojiDelete},
			},
		},
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func remove(slice []string, val string) []string {
	for i, item := range slice {
		if item == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
