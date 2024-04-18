package handler

import (
	"log"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/pkg/processer"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Handler struct {
	processer *processer.Processer
}

func NewHandler(processer *processer.Processer) *Handler {
	return &Handler{
		processer: processer,
	}
}

func (h *Handler) HandleCallback(w http.ResponseWriter, req *http.Request) {
	cb := h.processer.ParseRequest(w, req)
	if cb == nil {
		return
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				if message.Text == "hello" {
					h.processer.ExportSheet()
				}
				responseMessage := h.processer.AnalyzeTextMessage(message)
				_ = h.processer.ReplyMessage(responseMessage, e.ReplyToken)
			default:
				_ = h.processer.ReplyMessage("bye", e.ReplyToken)
			}
		default:
			log.Printf("Unsupported message: %T\n", event)
		}
	}
}
