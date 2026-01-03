package storage

import (
	"encoding/json"
	"os"
)

type BotData struct {
	LastOWPatch  string            `json:"last_ow_patch"`  // Stores Patch Title/Date
	LastArcPatch string            `json:"last_arc_patch"` // Stores Steam/Site URL
	WatchEnabled bool              `json:"watch_enabled"`
	TrackedGames map[string]string `json:"tracked_games"` // GameName -> ChannelID
}

const filePath = "data.json"

func SaveData(data BotData) error {
	file, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(filePath, file, 0644)
}

func LoadData() BotData {
	data := BotData{
		WatchEnabled: true,
		TrackedGames: make(map[string]string),
	}
	file, err := os.ReadFile(filePath)
	if err == nil {
		json.Unmarshal(file, &data)
	}
	if data.TrackedGames == nil {
		data.TrackedGames = make(map[string]string)
	}
	return data
}
