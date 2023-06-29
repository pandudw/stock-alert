package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	telegramBotToken       = "2343r4"
	chatID           int64 = 563
)

type Stock struct {
	Symbol      string  `json:"symbol"`
	TargetPrice float64 `json:"target_price"`
}

func main() {
	stocks := []Stock{
		{Symbol: "AMRT", TargetPrice: 2500.0},
		{Symbol: "TOWR", TargetPrice: 2500.0},
	}

	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		for _, stock := range stocks {
			price, err := getStockPrice(stock.Symbol)
			if err != nil {
				log.Printf("Maaf, tidak bisa mendapatkan harga saham untuk %s: %v\n", stock.Symbol, err)
				continue
			}

			if price >= stock.TargetPrice {
				message := fmt.Sprintf("Harga saham %s telah mencapai %.2f nihh, buruan!", stock.Symbol, price)
				sendTelegramAlert(message)
			}
		}

	}
}

func getStockPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://api.stockbit.com/v2/market/stock/%s", symbol)
	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var stockData struct {
		LastPrice float64 `json:"last"`
	}

	err = json.Unmarshal(data, &stockData)
	if err != nil {
		return 0, err
	}
	return stockData.LastPrice, nil
}

func sendTelegramAlert(message string) {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Println("error bot:", err)
		return
	}
	msg := tgbotapi.NewMessage(chatID, message)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("Error sending Telegram alert:", err)
	}
}
