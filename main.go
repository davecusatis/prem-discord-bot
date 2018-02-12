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
	godChannelID := mustGetConfigValue("GOD_CHANNEL_ID")

	// TODO: add db config

	discSession, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("Error creating discord session")
	}
	err = discSession.Open()
	defer discSession.Close()
	if err != nil {
		log.Fatal("Error opening discord connection")
	}
	discSession.AddHandler(messageHandler)

	db, err := newDatabase()
	http.Handle("/webhook", &SellyHandler{
		discordSession: discSession,
		db:             db,
	})

	http.Handle("/god", &GodHandler{
		godChannelID:   godChannelID,
		discordSession: discSession,
	})

	log.Printf("Starting bot on port %s", botport)
	log.Fatal(http.ListenAndServe(":"+botport, nil))
}
