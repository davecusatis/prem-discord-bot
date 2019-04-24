package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	coinApi "github.com/miguelmota/go-coinmarketcap"
)

const (
	// PriceCheckCMD is the command to check coinmarketcap price
	PriceCheckCMD = ".price"

	// PriceCheckStockCMD is the command to check coinmarketcap price
	PriceCheckStockCMD = ".stock"

	// TopTenCMD is the command to check the top 10 on cmc
	TopTenCMD = ".top"

	// SpaceDelimiter is the delimiter to split commands on
	SpaceDelimiter = " "
)

func messageHandler(disc *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == disc.State.User.ID {
		return
	}
	if !strings.Contains(m.Content, PriceCheckCMD) && !strings.Contains(m.Content, TopTenCMD) && !strings.Contains(m.Content, PriceCheckStockCMD) {
		return
	}

	cmdArray := strings.Split(m.Content, SpaceDelimiter)
	if cmdArray[0] == PriceCheckCMD {
		if len(cmdArray) < 2 {
			_, _ = disc.ChannelMessageSend(m.ChannelID, "Please provide a valid coin.")
			return
		}
		err := sendPriceCheckMessage(disc, m.ChannelID, strings.ToLower(cmdArray[1]))
		if err != nil {
			log.Printf("Error checking price %s", err)
		}
		return
	}

	if cmdArray[0] == PriceCheckStockCMD {
		if len(cmdArray) < 2 {
			_, _ = disc.ChannelMessageSend(m.ChannelID, "Please provide a valid stock ticker.")
			return
		}
		err := sendStockCheckMessage(disc, m.ChannelID, strings.ToUpper(cmdArray[1]))
		if err != nil {
			log.Printf("Error checking price %s", err)
		}
		return
	}

	if cmdArray[0] == TopTenCMD {
		err := sendTopTenInfoMessage(disc, m.ChannelID)
		if err != nil {
			log.Printf("Error checking top ten %s", err)
		}
		return
	}
}

func checkPrice(ticker string) (coinApi.Coin, error) {
	if val, ok := tickerMap[ticker]; ok {
		return coinApi.GetCoinData(val)
	}
	return coinApi.GetCoinData(ticker)
}

func sendStockCheckMessage(disc *discordgo.Session, channelID, ticker string) error {
	quoteValue, err := ac.StockQuote(ticker)
	if err != nil {
		_, _ = disc.ChannelMessageSend(channelID, fmt.Sprintf("Error getting info for %s from alpha vantage: %s.", ticker, err.Error()))
		log.Printf("error stockquote response: %s", err)
		return err
	}

	fields := make([]*discordgo.MessageEmbedField, 7)

	fields[0] = new(discordgo.MessageEmbedField)
	fields[0].Name = "Current Price"
	fields[0].Value = fmt.Sprintf("%.2f", quoteValue.Price)
	fields[0].Inline = true

	fields[1] = new(discordgo.MessageEmbedField)
	fields[1].Name = "Open"
	fields[1].Value = fmt.Sprintf("%.2f", quoteValue.Open)
	fields[1].Inline = true

	fields[2] = new(discordgo.MessageEmbedField)
	fields[2].Name = "High"
	fields[2].Value = fmt.Sprintf("%.2f", quoteValue.High)
	fields[2].Inline = true

	fields[3] = new(discordgo.MessageEmbedField)
	fields[3].Name = "Low"
	fields[3].Value = fmt.Sprintf("%.2f", quoteValue.Low)
	fields[3].Inline = true

	fields[4] = new(discordgo.MessageEmbedField)
	fields[4].Name = "Change"
	fields[4].Value = fmt.Sprintf("%.2f", quoteValue.Change)
	fields[4].Inline = true

	fields[5] = new(discordgo.MessageEmbedField)
	fields[5].Name = "Change %"
	fields[5].Value = fmt.Sprintf("%s", quoteValue.ChangePercent)
	fields[5].Inline = true

	fields[6] = new(discordgo.MessageEmbedField)
	fields[6].Name = "Volume"
	fields[6].Value = fmt.Sprintf("%d", quoteValue.Volume)
	fields[6].Inline = true

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

func sendPriceCheckMessage(disc *discordgo.Session, channelID, ticker string) error {
	coinInfo, err := checkPrice(ticker)
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

func generateTopTenMessage(coins map[string]coinApi.Coin) string {
	rankMap := make(map[int]string)
	for _, coin := range coins {
		rankMap[coin.Rank] = fmt.Sprintf("%d. %s $%.2f\n", coin.Rank, coin.Name, coin.PriceUsd)
	}

	i := 1
	var buffer bytes.Buffer
	for i <= 10 {
		buffer.WriteString(rankMap[i])
		i++
	}
	return buffer.String()
}

func sendTopTenInfoMessage(disc *discordgo.Session, channelID string) error {
	topCoins, err := coinApi.GetAllCoinData(10)
	if err != nil {
		_, _ = disc.ChannelMessageSend(channelID, fmt.Sprint("Error getting top coins from coinmarketcap. Please try again later"))
		return err
	}
	_, err = disc.ChannelMessageSend(channelID, generateTopTenMessage(topCoins))
	if err != nil {
		log.Printf("Error sending message to channel %s: %s", channelID, err)
		return err
	}
	return nil
}
