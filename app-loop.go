package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/pseudoelement/coin-ping/api"
)

func runLoop(coins []ArgCoinInfo, apiKey string, settngs Settings) {
	for {
		var cmcQuotesResp api.CoinMarketCapCoinQuotesLatestResp
		var errResp api.CoinMarketCapErrorResp
		err, success := api.Get(
			"https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest",
			&cmcQuotesResp,
			&errResp,
			[][2]string{{"symbol", coinsToSymbol(coins)}},
			[][2]string{{"X-CMC_PRO_API_KEY", apiKey}},
		)
		if err != nil {
			showNotification("Critical error", err.Error())
		}
		if !success {
			showNotification("CMC api error", errResp.Status.ErrorMessage)
		}

		checkPrices(cmcQuotesResp, coins)

		time.Sleep(time.Duration(settngs.DelayMinutes) * time.Minute)
	}
}

func coinsToSymbol(coins []ArgCoinInfo) string {
	var buf bytes.Buffer
	for idx, coin := range coins {
		buf.WriteString(coin.Symbol)
		if idx < len(coins)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

func findCoinBySymbol(symbol string, coins []ArgCoinInfo) *ArgCoinInfo {
	for _, coin := range coins {
		if coin.Symbol == symbol {
			return &coin
		}
	}

	return nil
}

func showNotification(title, content string) error {
	command := fmt.Sprintf(`display notification "%s" with title "%s" sound name "default"`, content, title)
	cmd := exec.Command("osascript", "-e", command)
	err := cmd.Run()

	return err
}

func checkPrices(cmcQuotesResp api.CoinMarketCapCoinQuotesLatestResp, coins []ArgCoinInfo) {
	for currencySymbol, matchedSymbols := range cmcQuotesResp.Data {
		firstMatch := matchedSymbols[0]
		currentPrice := firstMatch.Quote["USD"].Price
		argCoin := findCoinBySymbol(currencySymbol, coins)

		if currentPrice > argCoin.Top {
			strPrice := strconv.FormatFloat(currentPrice, 'f', -1, 64)
			showNotification("High price of "+argCoin.Symbol, "Price is "+strPrice+".")
		}
		if currentPrice < argCoin.Bottom {
			strPrice := strconv.FormatFloat(currentPrice, 'f', -1, 64)
			showNotification("Low price of "+argCoin.Symbol, "Price is "+strPrice+".")
		}
	}
}
