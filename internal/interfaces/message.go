package interfaces

import (
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/types"
)

type MessageInterface interface {
	ParseRequest(w http.ResponseWriter, req *http.Request) ([]*types.ParsedMessage, error)
	ReplyMessage(*types.ParsedMessage) error
	HandleImageContent(messageId string, callback func(string) error) error
}
