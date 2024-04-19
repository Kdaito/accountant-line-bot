package service

import (
	"fmt"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/interfaces"
)

type CallbackService struct {
	Drive interfaces.DriveInterface
}

func (c *CallbackService) Callback(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("こんにちは")
}
