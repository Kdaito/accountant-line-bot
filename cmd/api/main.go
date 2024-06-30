package main

import (
	"flag"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/router"
)

func main() {
	port := os.Getenv("PORT")

	isSkipGpt := flag.Bool("isSkipGpt", false, "Whether to skip receipt analysis process using sample data to save GPT resources")

	flag.Parse()

	router := &router.Router{Port: port}

	router.Set(*isSkipGpt)
	router.Run()
}
