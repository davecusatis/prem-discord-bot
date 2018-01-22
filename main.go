package main

import (
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func main() {
	parseConfig()
	appSecret = mustGetConfigValue("APP_SECRET")
	discordToken = mustGetConfigValue("DISCORD_TOKEN")
	channelID = mustGetConfigValue("CHANNEL_ID")
	botport := getConfigValue("BOT_PORT", "8000")

	var err error
	discordSession, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("Error creating discord session")
	}
	err = discordSession.Open()
	if err != nil {
		log.Fatal("Error opening discord connection")
	}
	http.HandleFunc("/webhook", sellyHandler)
	log.Printf("Starting bot on port %s", botport)
	log.Fatal(http.ListenAndServe(":"+botport, nil))
}