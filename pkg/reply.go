package pkg

import (
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Reply struct {
	bot *messaging_api.MessagingApiAPI
}

func newReply(bot *messaging_api.MessagingApiAPI) *Reply {
	return &Reply{
		bot: bot,
	}
}

func (r *Reply) ReplyMessage(message string, event webhook.MessageEvent) error {
	_, err := r.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: event.ReplyToken,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: message,
			},
		},
	})

	return err
}