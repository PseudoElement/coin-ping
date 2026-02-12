package main

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ArgCoinInfo struct {
	Symbol string
	Top    float64
	Bottom float64
}

type Settings struct {
	DelayMinutes int
}

var logs_file io.Writer = nil

func main() {
	godotenv.Load()
	apiKey := os.Getenv("COIN_MARKET_CAP_API_KEY")
	if apiKey == "" {
		panic("COIN_MARKET_CAP_API_KEY var is required.")
	}

	err := openLogsFile()
	if err != nil {
		panic("Failed open logs.txt. Error: " + err.Error())
	}

	argCoins, settings := parseCmdArgs()

	runLoop(argCoins, apiKey, settings)
}

func openLogsFile() error {
	file, err := os.OpenFile("./logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	logs_file = file
	return nil
}

func parseCmdArgs() ([]ArgCoinInfo, Settings) {
	coins := make([]ArgCoinInfo, 0)
	settings := Settings{DelayMinutes: 15}

	for _, arg := range os.Args {
		tokenData, found := strings.CutPrefix(arg, "--token=")
		if found {
			splitted := strings.Split(tokenData, ",")
			if len(splitted) != 3 {
				panic(tokenData + " is invalid --token.")
			}

			symbol := splitted[0]
			if len(symbol) < 1 {
				panic(splitted[0] + " is invalid symbol.")
			}
			bottom, err := strconv.ParseFloat(splitted[1], 64)
			if err != nil {
				panic(splitted[1] + " is invalid low amount.")
			}
			top, err := strconv.ParseFloat(splitted[2], 64)
			if err != nil {
				panic(splitted[2] + " is invalid low amount.")
			}

			coins = append(coins, ArgCoinInfo{
				Symbol: symbol,
				Top:    top,
				Bottom: bottom,
			})
			continue
		}

		intervalMinutes, found := strings.CutPrefix(arg, "--interval=")
		if found {
			minutes, err := strconv.Atoi(intervalMinutes)
			if err != nil {
				panic(intervalMinutes + " is invalid --interval.")
			}
			settings.DelayMinutes = minutes
		}
	}

	if len(coins) == 0 {
		panic("--token param not provided")
	}

	return coins, settings
}
