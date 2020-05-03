package main

import (
	"context"
	"flag"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/tardisman5197/discord-bot/bot"
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{
		"package": "main",
	})

	logger.Info("Staring Discord Bot")

	var token, mongoURI, databaseName string
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&mongoURI, "m", "", "Mongo URI")
	flag.StringVar(&databaseName, "d", "dev", "Database Name")

	flag.Parse()

	discordBot := bot.NewBot(token, mongoURI)

	err := discordBot.Setup(databaseName, "servers")
	if err != nil {
		logger.WithError(err).Error("Error Setting up Bot")
	}

	monitorCTX, monitorCancel := context.WithCancel(context.Background())
	monitorError := discordBot.MonitorMongoConnection(monitorCTX)

	done := discordBot.Start()

mainLoop:
	for {
		select {
		case <-done:
			logger.Info("Bot Finished")
			break mainLoop
		case err := <-monitorError:
			logger.WithError(err).Error("Error with Mongo Connection")
			break mainLoop
		}
	}

	monitorCancel()
	discordBot.Shutdown()

	logger.Info("Discord Bot Stopped")
}
