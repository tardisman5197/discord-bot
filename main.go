package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tardisman5197/discord-bot/bot"
)

func main() {
	var token, mongoURI string
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&mongoURI, "m", "", "Mongo URI")

	flag.Parse()

	discordBot := bot.NewBot(token, mongoURI)

	err := discordBot.Setup("dev", "servers")
	if err != nil {
		fmt.Printf("Error setting up bot - %v\n", err)
	}

	monitorCTX, monitorCancel := context.WithCancel(context.Background())
	monitorError := discordBot.MonitorMongoConnection(monitorCTX)

	done := discordBot.Start()

mainLoop:
	for {
		select {
		case <-done:
			fmt.Println("Bot Finished")
			break mainLoop
		case err := <-monitorError:
			fmt.Printf("Error with mongo connection - %v\n", err)
			break mainLoop
		}
	}

	monitorCancel()
	discordBot.Shutdown()
}
