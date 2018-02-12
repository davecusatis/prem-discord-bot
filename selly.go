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
	appSecret    string
	discordToken string
	channelID    string
	guildID      string
	roleID       string
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

	if !validOrder(order) {
		log.Printf("Invalid order %s for user: %s ", order.ID, order.Email)
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
	err = h.db.addOrUpdateUser(userToAdd)
	if err != nil {
		log.Printf("Error %s adding user to database: %#v", err, userToAdd)
	}

	err = addPaidRoleToUser(h.discordSession, userToAdd.discordTag)
	if err != nil {
		log.Printf("Error adding paid role to user %s: %s", userToAdd.discordTag, err)
		_, _ = h.discordSession.ChannelMessageSend(
			channelID,
			fmt.Sprintf("@davethecust#6318 @Tower#6969, these rat basterds didn't do the discord thing properly: %#v", userToAdd.discordTag))
	}

	_, err = h.discordSession.ChannelMessageSendEmbed(channelID, embedFromOrder(order))
	if err != nil {
		log.Printf("Error sending message to channel %s: %s", channelID, err)
		return
	}
}

func addPaidRoleToUser(session *discordgo.Session, user string) error {
	members, err := session.GuildMembers(guildID, "", 1000)
	if err != nil {
		return fmt.Errorf("Discord error while trying to search guild members: %s", err)
	}

	userID := ""
	for _, member := range members {
		if fmt.Sprintf("%s#%s", member.User.Username, member.User.Discriminator) == user {
			userID = member.User.ID
		}
	}
	if userID == "" {
		return fmt.Errorf("Unable to assign Paid role to user %s", userID)
	}

	err = session.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return fmt.Errorf("Discord error while trying to add role to user: %s", err)
	}

	return nil
}

func validOrder(order Order) bool {
	if order.Status != 100 {
		return false
	}
	return true
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
