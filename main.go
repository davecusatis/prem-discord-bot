package main

import (
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	av "github.com/fabianbaier/go-alpha-vantage"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
)

var (
	tickerMap           map[string]float64
	mutex               *sync.Mutex
	ac                  *av.Client
	discordToken        string
	coinMarketCapClient *cmc.Client
)

// func updateTickerMap() {
// 	mutex.Lock()
// 	coins, err := coinMarketCapClient.Cryptocurrency.Map(&cmc.MapOptions{
// 		ListingStatus: "active",
// 		Start:         1,
// 	})
// 	if err != nil {
// 		log.Printf("error getting cmc map: %s", err)
// 	}
// 	for _, coin := range coins {
// 		sym := strings.ToLower(coin.Symbol)
// 		tickerMap[sym] = coin.ID
// 	}
// 	mutex.Unlock()
// }

func main() {
	parseConfig()

	mutex = &sync.Mutex{}
	tickerMap = make(map[string]float64)

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

	cmcToken := mustGetConfigValue("CMC_API_KEY")
	coinMarketCapClient = cmc.NewClient(&cmc.Config{
		ProAPIKey: cmcToken,
	})

	// go updateTickerMap()
	ticker := time.NewTicker(3600 * time.Second)
	quit := make(chan struct{})
	func() {
		for ; true; <-ticker.C {
			select {
			case <-ticker.C:
				log.Printf("Syncing with coinmarketcap")
				// go updateTickerMap()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	log.Println("GOT HERE")
}
