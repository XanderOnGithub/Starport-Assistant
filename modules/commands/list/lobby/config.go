package lobby

const (
	// Visual Assets
	CustomThumbnail = "https://media.discordapp.net/attachments/581641208170938389/1457113193020264468/ServerIcon2026Attempt4.gif?ex=695ad1e9&is=69598069&hm=829b5056c3e75b3c67cfaa24229f22dc49863a92c9cf6052fdc5d349a5a5c849&="
	EmbedColor      = 0x18176c

	// Embed Text (Recruitment Phase)
	EmbedTitlePrefix = "🚀 **LOBBY CREATED:** "
	EmbedDescLine1   = "Join **%s** for some games!"
	EmbedTimeLabel   = "🕒 **STARTING:**"
	EmbedVoiceLabel  = "🔊 **VOICE:**"
	EmbedFooterText  = "Starport Assistant • New Lobby"

	// Live Message Config (The "Cool" Embed)
	LiveEmbedTitle = "🚀 LOBBY: BLAST OFF"
	LiveEmbedDesc  = "The lobby for **%s** is now active!\n\n**Gathering at:** <#%s>\n\n**Squad:** %s"
	LiveEmbedColor = 0x18176c
	LiveFooterText = "Starport Assistant • Active Lobby"

	// Slot List Styling
	SlotsHeader = "╭─── ▼ **CURRENT PARTY** ▼ ───\n"
	SlotsFooter = "╰────────────────────────"
	IconHost    = "👑"
	IconPlayer  = "⭐"
	IconEmpty   = "🌑"
	TextEmpty   = "Empty Slot"

	// Interaction Messages
	MsgLobbyLive  = "🚀 **LOBBY STARTED!**\n%s\nGather here: <#%s>"
	MsgLobbyError = "❌ This lobby session is no longer active."

	// Button Labels
	BtnJoin   = "Join"
	BtnLeave  = "Leave"
	BtnStart  = "Start"
	BtnDelete = "Close"

	// Button Emojis
	EmojiJoin   = "🎮"
	EmojiLeave  = "🚪"
	EmojiStart  = "🚀"
	EmojiDelete = "🗑️"
)
