package sens

type GameInfo struct {
	DisplayName string
	Ratio       float64 // Based on Overwatch 2 Sens 1.0
}

const (
	SensTitle  = "🚀 STARPORT: SENSITIVITY CONVERTER"
	SensFooter = "Starport Sensitivity Converter • Conversion Complete"
	SensColor  = 0x18176c
)

var GameData = map[string]GameInfo{
	"overwatch2":   {DisplayName: "Overwatch 2", Ratio: 1.0},
	"thefinals":    {DisplayName: "THE FINALS", Ratio: 6.6},
	"cs2":          {DisplayName: "Counter-Strike 2", Ratio: 0.3},
	"valorant":     {DisplayName: "Valorant", Ratio: 0.09434},
	"marvelrivals": {DisplayName: "Marvel Rivals", Ratio: 0.37714},
	"cod":          {DisplayName: "Call of Duty", Ratio: 1.0},
	"tarkov":       {DisplayName: "Escape From Tarkov", Ratio: 0.053},
	"apex":         {DisplayName: "Apex Legends", Ratio: 0.3},
	"deadlock":     {DisplayName: "Deadlock", Ratio: 0.150},
	"fortnite":     {DisplayName: "Fortnite", Ratio: 1.188},
	"arcraiders":   {DisplayName: "ARC Raiders", Ratio: 4.849},
	"battlefield6": {DisplayName: "Battlefield 6", Ratio: 2.468},
	"valheim":      {DisplayName: "Valheim", Ratio: 0.132},
	"roblox":       {DisplayName: "Roblox", Ratio: 0.017},
	"rust":         {DisplayName: "Rust", Ratio: 0.059},
	"splitgate2":   {DisplayName: "Splitgate 2", Ratio: 0.591},
	"tf2":          {DisplayName: "Team Fortress 2", Ratio: 0.3},
}
