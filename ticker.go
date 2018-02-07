package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	coinApi "github.com/miguelmota/go-coinmarketcap"
)

const (
	// PriceCheckCMD is the command to check coinmarketcap price
	PriceCheckCMD = ".price"

	// TopTenCMD is the command to check the top 10 on cmc
	TopTenCMD = ".topten"

	// SpaceDelimiter is the delimiter to split commands on
	SpaceDelimiter = " "
)

func messageHandler(disc *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == disc.State.User.ID {
		return
	}

	if !strings.Contains(m.Content, PriceCheckCMD) {
		return
	}

	cmdArray := strings.Split(m.Content, SpaceDelimiter)
	if cmdArray[0] == PriceCheckCMD {
		err := sendPriceCheckMessage(disc, m.ChannelID, cmdArray[1])
		if err != nil {
			log.Printf("Error checking price %s", err)
		}
		return
	}

	// if cmdArray[0] == TopTenCMD {
	// 	err := sendTopTenInfoMessage(disc)
	// 	return
	// }
}

func sendPriceCheckMessage(disc *discordgo.Session, channelID, ticker string) error {
	coinInfo, err := coinApi.GetCoinData(ticker)
	if err != nil {
		_, _ = disc.ChannelMessageSend(channelID, fmt.Sprintf("Error getting info for %s from coinmarketcap: %s.", ticker, err.Error()))
		return err
	}

	fields := make([]*discordgo.MessageEmbedField, 6)

	fields[0] = new(discordgo.MessageEmbedField)
	fields[0].Name = "Price USD"
	fields[0].Value = fmt.Sprintf("%.2f", coinInfo.PriceUsd)
	fields[0].Inline = true

	fields[1] = new(discordgo.MessageEmbedField)
	fields[1].Name = "Price BTC"
	fields[1].Value = fmt.Sprintf("%f", coinInfo.PriceBtc)
	fields[1].Inline = true

	fields[2] = new(discordgo.MessageEmbedField)
	fields[2].Name = "Market Cap"
	fields[2].Value = fmt.Sprintf("%.2f", coinInfo.MarketCapUsd)
	fields[2].Inline = true

	fields[3] = new(discordgo.MessageEmbedField)
	fields[3].Name = "Change 1hr"
	fields[3].Value = fmt.Sprintf("%.2f%%", coinInfo.PercentChange1h)
	fields[3].Inline = true

	fields[4] = new(discordgo.MessageEmbedField)
	fields[4].Name = "Change 24hr"
	fields[4].Value = fmt.Sprintf("%.2f%%", coinInfo.PercentChange24h)
	fields[4].Inline = true

	fields[5] = new(discordgo.MessageEmbedField)
	fields[5].Name = "Change 7d"
	fields[5].Value = fmt.Sprintf("%.2f%%", coinInfo.PercentChange7d)
	fields[5].Inline = true

	embed := new(discordgo.MessageEmbed)
	embed.Description = fmt.Sprintf("Grabbing latest data for %s", ticker)
	embed.Fields = fields

	_, err = disc.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Printf("Error sending message to channel %s: %s", channelID, err)
		return err
	}
	return nil
}

// func sendTopTenInfoMessage(disc *discordgo.Session) error {
// 	coinInfo, err := coinApi.GetAllCoinData(10)
// 	if err != nil {
// 		_, _ = disc.ChannelMessageSend(channelID, fmt.Sprintf("Error getting info for %s from coinmarketcap. Please try again later.", ticker))
// 		return err
// 	}

// 	fields := make([]*discordgo.MessageEmbedField, 10)

// 	fields[0] = new(discordgo.MessageEmbedField)
// 	fields[0].Name = "Price USD"
// 	fields[0].Value = fmt.Sprint(coinInfo.PriceUsd)
// 	fields[0].Inline = true

// 	fields[1] = new(discordgo.MessageEmbedField)
// 	fields[1].Name = "Price BTC"
// 	fields[1].Value = fmt.Sprint(coinInfo.PriceBtc)
// 	fields[1].Inline = true

// 	fields[2] = new(discordgo.MessageEmbedField)
// 	fields[2].Name = "Market Cap"
// 	fields[2].Value = fmt.Sprint(coinInfo.MarketCapUsd)
// 	fields[2].Inline = true

// 	fields[3] = new(discordgo.MessageEmbedField)
// 	fields[3].Name = "Percentage Change 1hr"
// 	fields[3].Value = fmt.Sprint(coinInfo.PercentChange1h)
// 	fields[3].Inline = true

// 	fields[4] = new(discordgo.MessageEmbedField)
// 	fields[4].Name = "Percentage Change 24hr"
// 	fields[4].Value = fmt.Sprint(coinInfo.PercentChange24h)
// 	fields[4].Inline = true

// 	fields[5] = new(discordgo.MessageEmbedField)
// 	fields[5].Name = "Percentage Change 7d"
// 	fields[5].Value = fmt.Sprint(coinInfo.PercentChange7d)
// 	fields[5].Inline = true

// 	embed := new(discordgo.MessageEmbed)
// 	embed.Description = fmt.Sprintf("Grabbing latest data for %s", ticker)
// 	embed.Fields = fields

// 	_, err = disc.ChannelMessageSendEmbed(channelID, embed)
// 	if err != nil {
// 		log.Printf("Error sending message to channel %s: %s", channelID, err)
// 		return err
// 	}
// 	return nil
// }