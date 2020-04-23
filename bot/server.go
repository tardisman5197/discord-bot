package bot

// DiscordServer has all of the functions that the
// not can perform for a server.
type DiscordServer struct {
	// guildID is the unique ID given to a discord
	// server
	guildID string

	// lists stores a collection of items assigned
	// to a list id.
	lists map[string][]string
}

// NewDiscordServer returns a pointer to a DiscordServer
// instance.
func NewDiscordServer(guildID string) *DiscordServer {
	return &DiscordServer{
		guildID: guildID,
		lists:   make(map[string][]string, 0),
	}
}
