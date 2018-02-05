package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	appSecret      string
	discordToken   string
	channelID      string
	discordSession *discordgo.Session
)

// SellyHandler is the type representing the handler type
type SellyHandler struct {
	db             *Database
	discordSession *discordgo.Session
}

// Order is the tpe for selly orders
type Order struct {
	ID            string            `json:"id"`
	ProductID     string            `json:"product_id"`
	Email         string            `json:"email"`
	IPAddress     string            `json:"ip_address"`
	CountryCode   string            `json:"country_code"`
	UserAgent     string            `json:"user_agent"`
	Value         string            `json:"value"`
	Currency      string            `json:"currency"` // The ISO 4217 currency code
	Gateway       string            `json:"gateway"`
	RiskLevel     int               `json:"risk_level"`
	Status        int               `json:"status"`
	Delivered     string            `json:"delivered"`
	CryptoValue   interface{}       `json:"crypto_value"`
	CryptoAddress interface{}       `json:"crypto_address"`
	Referral      string            `json:"referral"`
	USDValue      string            `json:"usd_value"`
	ExchangeRate  string            `json:"exchange_rate"`
	Custom        map[string]string `json:"custom"`
	WebhookType   int               `json:"webhook_type"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
}

func (h *SellyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	secret := req.URL.Query().Get("secret")
	if secret != appSecret {
		log.Print("Wrong secret provided")
		return
	}

	var order Order
	err := json.NewDecoder(req.Body).Decode(&order)
	if err != nil {
		log.Printf("Error parsing body: %s", err)
		return
	}

	_, err = h.discordSession.ChannelMessageSendEmbed(channelID, embedFromOrder(order))
	if err != nil {
		log.Printf("Error sending message to channel %s: %s", channelID, err)
		return
	}

	duration, err := h.db.getDurationByProductID(order.ProductID)
	if err != nil {
		log.Printf("Error getting duration for product id %s", order.ProductID)
		return
	}

	now := time.Now().UTC().UnixNano()
	// todo: add OR update
	userToAdd := &User{
		email:      order.Email,
		product:    order.ProductID,
		discordTag: order.Custom["0"],
		startDate:  now,
		endDate:    now + duration.Nanoseconds(),
	}
	err = h.db.addUser(userToAdd)
	if err != nil {
		log.Printf("Error %s adding user to database: %#v", err, userToAdd)
	}
}

func embedFromOrder(order Order) *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, 5)

	fields[0] = new(discordgo.MessageEmbedField)
	fields[0].Name = "ID"
	fields[0].Value = order.ID
	fields[0].Inline = true

	fields[1] = new(discordgo.MessageEmbedField)
	fields[1].Name = "Email"
	fields[1].Value = order.Email
	fields[1].Inline = true

	fields[2] = new(discordgo.MessageEmbedField)
	fields[2].Name = "Value"
	fields[2].Value = order.Value
	fields[2].Inline = true

	fields[3] = new(discordgo.MessageEmbedField)
	fields[3].Name = "Discord"
	fields[3].Value = order.Custom["0"]
	fields[3].Inline = true

	fields[4] = new(discordgo.MessageEmbedField)
	fields[4].Name = "Product"
	fields[4].Value = order.ProductID
	fields[4].Inline = true

	embed := new(discordgo.MessageEmbed)
	thumbnail := new(discordgo.MessageEmbedThumbnail)
	thumbnail.URL = "https://selly.gg/images/apple-touch-icon-180x180.png"
	embed.Thumbnail = thumbnail
	embed.Description = fmt.Sprintf("Completed Order from %s", order.Custom["0"])
	embed.Fields = fields
	return embed
}
