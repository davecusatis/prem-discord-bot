# Discord Crypto/Stock price checker

## Usage
1. edit env file with your discord, and alpha vantage api tokens.
2. `go build`
3. `./prem-discord-bot`

## Supported commands
1. `.price <crypto ticker>` query coinmarketcap api for crypto data.
2. `.stock <stock ticker>` query alpha vantage for [daily stock quote](https://www.alphavantage.co/documentation/#daily)
3. `.top` query coinmarketcap for top crypto currencies by market cap
