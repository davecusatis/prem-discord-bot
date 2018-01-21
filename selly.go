package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	appSecret      string
	discordToken   string
	channelID      string
	discordSession *discordgo.Session
)

// Order is the tpe for selly orders
type Order struct {
	ID            string      `json:"id"`
	ProductID     string      `json:"product_id"`
	Email         string      `json:"email"`
	IPAddress     string      `json:"ip_address"`
	CountryCode   string      `json:"country_code"`
	UserAgent     string      `json:"user_agent"`
	Value         string      `json:"value"`
	Currency      string      `json:"currency"` // The ISO 4217 currency code
	Gateway       string      `json:"gateway"`
	RiskLevel     int         `json:"risk_level"`
	Status        int         `json:"status"`
	Delivered     string      `json:"delivered"`
	CryptoValue   interface{} `json:"crypto_value"`
	CryptoAddress interface{} `json:"crypto_address"`
	Referral      string      `json:"referral"`
	USDValue      float32     `json:"usd_value"`
	ExchangeRate  float32     `json:"exchange_rate"`
	Custom        interface{} `json:"custom"`
	WebhookType   int         `json:"webhook_type"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
}

func sellyHandler(w http.ResponseWriter, req *http.Request) {
	// get secret
	secret := req.URL.Query().Get("secret")
	if secret != appSecret {
		log.Print("Wrong secret provided")
	}

	var order Order
	err := json.NewDecoder(req.Body).Decode(&order)
	if err != nil {
		log.Printf("Error parsing body: %s", err)
	}

	// post discord message

	_, err = discordSession.ChannelMessageSendComplex(channelID, messageFromOrder(order))
	if err != nil {
		log.Printf("Error sending message to channel %s: %s", channelID, err)
	}
}

func messageFromOrder(order Order) *discordgo.MessageSend {
	var embed *discordgo.MessageEmbed
	var message *discordgo.MessageSend
	fields := make([]*discordgo.MessageEmbedField, 7)

	fields[0].Name = "ID"
	fields[0].Value = order.ID
	fields[0].Inline = true

	fields[1].Name = "Email"
	fields[1].Value = order.Email
	fields[1].Inline = true

	fields[2].Name = "Value"
	fields[2].Value = order.Value
	fields[2].Inline = true

	fields[3].Name = "Status"
	fields[3].Value = string(order.Status)
	fields[3].Inline = true

	fields[4].Name = "USDValue"
	fields[4].Value = fmt.Sprintf("%f", order.USDValue)
	fields[4].Inline = true

	fields[5].Name = "Custom"
	fields[5].Value = fmt.Sprint("%#v", order.Custom)
	fields[5].Inline = true

	embed.Fields = fields
	message.Embed = embed
	return message
}
