package bot

// DiscordServer has all of the functions that the
// not can perform for a server.
type DiscordServer struct {
	// guildID is the unique ID given to a discord
	// server
	GuildID string `bson:"guildID"`

	// lists stores a collection of items assigned
	// to a list id.
	Lists map[string][]string `bson:"lists"`
}

// NewDiscordServer returns a pointer to a DiscordServer
// instance.
func NewDiscordServer(guildID string) *DiscordServer {
	return &DiscordServer{
		GuildID: guildID,
		Lists:   make(map[string][]string, 0),
	}
}
