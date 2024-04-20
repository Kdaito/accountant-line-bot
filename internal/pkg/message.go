package pkg

import (
	"errors"
	"log"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/types"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Message struct {
	ChannelSecret string
	Bot           *messaging_api.MessagingApiAPI
}

func (m *Message) ParseRequest(w http.ResponseWriter, req *http.Request) ([]*types.ParsedMessage, error) {
	res := []*types.ParsedMessage{}

	cb, err := webhook.ParseRequest(m.ChannelSecret, req)
	if err != nil {
		log.Printf("Cannot parse request: %+v\n", err)
		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return res, err
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				res = append(res, &types.ParsedMessage{
					MessageType: types.MESSAGE_TYPE_TEXT,
					Text:        message.Text,
					ID:          message.Id,
					ReplyToken:  e.ReplyToken,
				})
			default:
				res = append(res, &types.ParsedMessage{
					MessageType: types.MESSAGE_TYPE_OTHERS,
					Text:        "このメッセージには対応できません。",
					ID:          "",
					ReplyToken:  e.ReplyToken,
				})
			}
		default:
			// 予期しないメッセージが含まれる場合は一旦無視
			log.Printf("Unsupported message: %T\n", event)
		}
	}

	return res, nil
}

func (m *Message) ReplyMessage(message *types.ParsedMessage) error {
	_, err := m.Bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: message.ReplyToken,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: message.Text,
			},
		},
	})

	return err
}
