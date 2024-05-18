package main

import (
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/router"
)

func main() {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")
	gptApiUrl := os.Getenv("GPT_API_URL")
	gptApiKey := os.Getenv("GPT_API_KEY")
	port := os.Getenv("PORT")

	router := &router.Router{Port: port}

	router.Set(channelSecret, channelToken, gptApiUrl, gptApiKey)
	router.Run()
}
