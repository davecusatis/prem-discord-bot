package main

import (
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

const (
	truestr      = "true"
	botPortConfg = "BOT_PORT"
)

func main() {
	parseConfig()
	appSecret = mustGetConfigValue("APP_SECRET")
	discordToken = mustGetConfigValue("DISCORD_TOKEN")
	channelID = mustGetConfigValue("CHANNEL_ID")
	botport := getConfigValue("BOT_PORT", "8000")

	var err error
	discordSession, err = discordgo.New(discordToken)
	if err != nil {
		log.Fatal("Error creating discord session")
	}

	http.HandleFunc("/webook", sellyHandler)
	http.ListenAndServe(":"+botport, nil)
}
