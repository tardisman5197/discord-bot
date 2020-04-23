package main

import (
	"flag"
	"fmt"

	"github.com/tardisman5197/discord-bot/bot"
)

func main() {
	var token string
	flag.StringVar(&token, "t", "", "Bot Token")

	flag.Parse()

	discordBot := bot.NewBot(token)

	err := discordBot.Setup()
	if err != nil {
		fmt.Printf("Error setting up bot - %v\n", err)
	}

	done := discordBot.Start()

	<-done

	discordBot.Shutdown()
}
