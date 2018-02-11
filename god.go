package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// GodHandler is the handler for god mode messages
type GodHandler struct {
	godChannelID   string
	discordSession *discordgo.Session
}

// GodText is the body we expect being passed to the god mode endpoint
type GodText struct {
	Text string `json:"text"`
}

func (h *GodHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	secret := req.URL.Query().Get("secret")
	if secret != appSecret {
		log.Print("Wrong secret provided")
		http.Error(w, "wrong secret", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}

	var godMessage GodText
	err = json.Unmarshal(body, &godMessage)
	if err != nil {
		log.Printf("couldn't parse god text: %s", err)
		http.Error(w, "bad god text", http.StatusBadRequest)
		return
	}
	log.Printf("god text: %#v", godMessage)
	_, err = h.discordSession.ChannelMessageSend(h.godChannelID, godMessage.Text)
	if err != nil {
		log.Printf("Error sending message to channel %s: %s", h.godChannelID, err)
		return
	}
}
