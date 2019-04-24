package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	av "github.com/fabianbaier/go-alpha-vantage"
	coinApi "github.com/miguelmota/go-coinmarketcap"
)

var (
	tickerMap    map[string]string
	mutex        *sync.Mutex
	ac           *av.Client
	discordToken string
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

	discordToken = mustGetConfigValue("DISCORD_TOKEN")
	avToken := mustGetConfigValue("ALPHA_VANTAGE_TOKEN")

	ac = av.NewClient(avToken)

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
}
