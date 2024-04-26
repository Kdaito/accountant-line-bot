package interfaces

import (
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/types"
)

type MessageInterface interface {
	ParseRequest(w http.ResponseWriter, req *http.Request) ([]*types.ParsedMessage, error)
	ReplyMessage(*types.ParsedMessage) error
	HandleImageContent(messageId string, callback func(*os.File) error) error
}
