package lobby

import (
	"fmt"
	"strings"
	"sync"
)

type GameLobby struct {
	MessageID       string
	ChannelID       string
	HostID          string
	HostDisplayName string
	GameName        string
	MaxSlots        int
	Players         []string
	StartTime       string
	VoiceChannelID  string
}

var (
	ActiveLobbies = make(map[string]*GameLobby)
	Mutex         sync.Mutex
)

func (l *GameLobby) RenderSlots() string {
	var builder strings.Builder
	builder.WriteString(SlotsHeader)

	for i := 0; i < l.MaxSlots; i++ {
		if i < len(l.Players) {
			playerID := l.Players[i]
			icon := IconPlayer
			if playerID == l.HostID {
				icon = IconHost
			}
			builder.WriteString(fmt.Sprintf("│ %s <@%s>\n", icon, playerID))
		} else {
			builder.WriteString(fmt.Sprintf("│ %s *%s*\n", IconEmpty, TextEmpty))
		}
	}

	builder.WriteString(SlotsFooter)
	return builder.String()
}
