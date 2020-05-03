package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Bot handels all of the incoming messages for the bot
// and stores all of the information for each of the
// discord servers
type Bot struct {
	// discordBot contains a pointer to a session
	// implemented in the discordgo module
	discordBot *discordgo.Session

	// token is the developer token provided by discord
	// to allow the bot to connect to the discord API
	token string

	// mongoClient is used to connect to the mongo database
	mongoClient *mongo.Client

	// mongoURI is the address of the mongo database which
	// the bot stores its data in
	mongoURI string

	// collection stores a pointer to the mongo collection
	// storing the discord servers data structures.
	collection *mongo.Collection

	logger *log.Entry
}

// NewBot returns a pointer to a new instance of a Bot
func NewBot(token, mongoURI string) *Bot {
	return &Bot{
		token:    token,
		mongoURI: mongoURI,
		logger: log.WithFields(log.Fields{
			"package": "bot",
		}),
	}
}

// Setup creates a new session with the token set when the bot
// was created. Setup returns an error if the Bot could not
// connect to the discord API.
func (b *Bot) Setup(databaseName, collectionName string) error {
	b.logger.Debug("Setting up bot")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + b.token)
	if err != nil {
		b.logger.WithError(err).Error("Error creating discord session")
		return err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(b.handleMessage)

	b.discordBot = dg

	ctx, _ := context.WithTimeout(context.Background(), MongoConnectionTime*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(b.mongoURI))
	if err != nil {
		b.logger.WithError(err).Error("Error connecting to Mongo Database")
		return err
	}

	b.mongoClient = client
	b.collection = b.mongoClient.Database(databaseName).Collection(collectionName)

	b.logger.Debug("Bot setup complete")
	return nil
}

// Start the discord bot listening for messsages from all of
// the servers the bot is connected to. This function runs
// forever in a goroutine (thread) and returns a channel. The
// channel has a true placed in it once this function has
// finished (when an error has occurred).
func (b *Bot) Start() chan bool {
	b.logger.Debug("Starting bot")

	done := make(chan bool)

	// Start listening in a new goroutine
	go func() {
		// Open a websocket connection to Discord and begin listening.
		err := b.discordBot.Open()
		if err != nil {
			b.logger.WithError(err).Error("Error opening connection")
			done <- true
		}

		b.logger.Info("Bot Started")
	}()

	return done
}

// Shutdown gracefully stops the discord websocket.
func (b *Bot) Shutdown() {
	b.logger.Debug("Shutting down bot")
	b.discordBot.Close()
}

// handleMessage is called when a new message is received by the Bot.
// This function checks for the command prefix, then splits up the
// message into the command keyword and arguments. Then the corresponding
// command function is called for that server.
func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Split up the received message

	// Do nothing if the message starts with a tilda
	if m.Content[0] != '~' {
		return
	}

	logger := b.logger.WithFields(log.Fields{
		"guildID": m.GuildID,
	})

	// Split the message into cmd and arguments
	// Expected message format: ~COMMAND ARG1 ARG2 ...
	args := strings.Split(m.Content, " ")

	// Get the command key word
	cmd := args[0]

	// No command key word given
	if len(cmd) <= 1 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.displayHelp()))
		logger.Debug("Invalid Command")
		return
	}

	// Remove the tilda from the command keyword
	cmd = cmd[1:]

	// Remove the command keyword from the arg list
	args = args[1:]

	// Check command keyword against the known commands
	// and run the correct server function
	switch cmd {
	case "add":
		logger.Debug("Add Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.add(m.GuildID, args)))
	case "remove":
		logger.Debug("Remove Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.removeItem(m.GuildID, args)))
	case "removeList":
		logger.Debug("Remove List Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.removeList(m.GuildID, args)))
	case "pick":
		logger.Debug("Pick Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.pick(m.GuildID, args)))
	case "list":
		logger.Debug("List Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.getList(m.GuildID, args)))
	case "lists":
		logger.Debug("Lists Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.getLists(m.GuildID)))
	default:
		logger.Debug("Help Command")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", b.displayHelp()))
	}
	return
}

// displayHelp returns the help text for the Bot
func (b *Bot) displayHelp() string {
	return `
	Commands:
		~add <list> <item(s)>
		~remove <list> <item(s)>
		~removeList <list>
		~pick <list>
		~list <list>
		~lists


	Usage:
		Add Command:

			~add <list> <item(s)>

			The add command will place a new item in a list.
			Example Usage: ~add prompts villain hero
			If a list name does not exist, it will be created.
			Multiple items can be added to a list by appending them
			to the end of the command, seperated by a space.


		Remove Command:

			~remove <list> <item(s)>

			The remove command will remove an item from a list.
			Example Usage: ~remove prompts hero
			To remove multiple items from a list, append the items
			to the end of the command separated by spaces.


		Remove List Command:

			~removeList <list>

			The remove list command removes an entire list.
			Example Usage: ~removeList prompts


		Pick Command:

			~pick <list>

			The pick command randomly chooses an item from a list.
			Example Usage: ~pick prompts
			If you would like the selected item to be removed when it
			is picked, append a word after the base command.
			Example Usage: ~pick prompts true


		List Command:

			~list <list>

			The list command displays all of the items in a given list.
			Example Usage: ~list prompts

		Lists Command:

			~lists
			
			The lists command displays all of the lists that are stored.
			Example Usage: ~lists
	`
}

// MonitorMongoConnection constantly pings the mongo database to ensure
// the Bot is till connected.
func (b *Bot) MonitorMongoConnection(ctx context.Context) chan error {
	b.logger.Debug("Starting Monitoring Mongo Connection")

	errChannel := make(chan error)

	ticker := time.NewTicker(PingInterval * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				pingCTX, _ := context.WithTimeout(context.Background(), MongoPingTime*time.Second)
				err := b.mongoClient.Ping(pingCTX, readpref.Primary())
				if err != nil {
					errChannel <- err
				}
			}
		}
	}()

	return errChannel
}
