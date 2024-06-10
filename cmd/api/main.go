package main

import (
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/router"
)

func main() {
	port := os.Getenv("PORT")

	router := &router.Router{Port: port}

	router.Set()
	router.Run()
}
