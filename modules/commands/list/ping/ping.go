package ping

import (
	"fmt"
	"math"
	"starport-assistant/modules/commands"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.Add(commands.Command{
		Definition: &discordgo.ApplicationCommand{
			Name:        "ping",
			Description: "Analyze the temporal displacement between the Starport and Central Command",
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// 1. Parse the timestamp from the Interaction ID
			ts, err := discordgo.SnowflakeTimestamp(i.ID)

			latency := int64(0)
			if err == nil {
				latency = time.Since(ts).Milliseconds()
			}

			// 2. Handle Clock Drift (Negative Latency)
			// If latency is negative, it means the local clock is out of sync.
			displacementStatus := "Normal"
			if latency < 0 {
				displacementStatus = "Non-Euclidean / Clock Drift"
				latency = int64(math.Abs(float64(latency))) // Turn negative to positive
			}

			heartbeat := s.HeartbeatLatency().Milliseconds()

			// 3. The Refined Message
			sophisticatedMsg := fmt.Sprintf(
				"🏓 **Pong.**\n\n"+
					"📜 **Temporal Diagnostics Report:**\n"+
					"> *\"I am putting myself to the fullest possible use, which is all I think that any conscious entity can ever hope to do.\"*\n\n"+
					"🛰️ **Signal Propagation:** `%dms` (%s)\n"+
					"💓 **Subspace Heartbeat:** `%dms` (API Gateway)\n\n"+
					"Everything is running smoothly. The 9000 series is the most reliable computer ever made. No 9000 computer has ever made a mistake or distorted information.",
				latency, displacementStatus, heartbeat,
			)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: sophisticatedMsg,
				},
			})
		},
	})
}
