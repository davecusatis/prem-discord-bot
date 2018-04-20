package main

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	coinApi "github.com/miguelmota/go-coinmarketcap"
)

var (
	tickerMap map[string]string
	mutex     *sync.Mutex
)

func updateTickerMap() {
	mutex.Lock()
	res, err := coinApi.GetMarketData()
	if err != nil {
		log.Print(err)
	}
	coins, err := coinApi.GetAllCoinData(res.ActiveCurrencies)
	for _, coin := range coins {
		sym := strings.ToLower(coin.Symbol)
		tickerMap[sym] = coin.ID
	}
	mutex.Unlock()
}

func main() {
	parseConfig()

	mutex = &sync.Mutex{}
	tickerMap = make(map[string]string)

	appSecret = mustGetConfigValue("APP_SECRET")
	discordToken = mustGetConfigValue("DISCORD_TOKEN")
	channelID = mustGetConfigValue("CHANNEL_ID")
	botport := getConfigValue("BOT_PORT", "8000")
	godChannelID := mustGetConfigValue("GOD_CHANNEL_ID")
	guildID = mustGetConfigValue("GUILD_ID")
	roleID = mustGetConfigValue("ROLE_ID")

	discSession, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Error creating discord session: %s", err)
	}
	err = discSession.Open()
	defer discSession.Close()
	if err != nil {
		log.Fatalf("Error opening discord connection: %s", err)
	}
	discSession.AddHandler(messageHandler)

	http.Handle("/god", &GodHandler{
		godChannelID:   godChannelID,
		discordSession: discSession,
	})

	go updateTickerMap()
	ticker := time.NewTicker(3600 * time.Second)
	quit := make(chan struct{})
	go func() {
		for ; true; <-ticker.C {
			select {
			case <-ticker.C:
				log.Printf("Syncing with coinmarketcap")
				go updateTickerMap()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	log.Printf("Starting bot on port %s", botport)
	log.Fatal(http.ListenAndServe(":"+botport, nil))
}
