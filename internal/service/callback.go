package service

import (
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/interfaces"
	"github.com/Kdaito/accountant-line-bot/internal/types"
)

type CallbackService struct {
	Drive   interfaces.DriveInterface
	Message interfaces.MessageInterface
}

func (c *CallbackService) Callback(w http.ResponseWriter, req *http.Request) {
	parsedMessages, err := c.Message.ParseRequest(w, req)

	if err != nil {
		return
	}

	for _, parsedMessage := range parsedMessages {
		switch parsedMessage.MessageType {
		case types.MESSAGE_TYPE_TEXT:
			c.Message.ReplyMessage(parsedMessage)
		case types.MESSAGE_TYPE_OTHERS:
			c.Message.ReplyMessage(parsedMessage)
		default:
			return
		}
	}
}
