package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/interfaces"
	"github.com/Kdaito/accountant-line-bot/internal/types"
)

type CallbackService struct {
	Drive   interfaces.DriveInterface
	Message interfaces.MessageInterface
	Sheet   interfaces.SheetInterface
}

func (c *CallbackService) Callback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	sheetForDrive, err := c.Sheet.CreateSheet(ctx)

	if err != nil {
		fmt.Printf("%v", err)
	}

	c.Drive.Upload("10KFLeu8bbmUfVUhDv0UjqL6XtEgWFG1N", sheetForDrive)

	if err != nil {
		fmt.Printf("%v", err)
	}

	// parsedMessages, err := c.Message.ParseRequest(w, req)

	// if err != nil {
	// 	return
	// }

	// for _, parsedMessage := range parsedMessages {
	// 	switch parsedMessage.MessageType {
	// 	case types.MESSAGE_TYPE_TEXT:
	// 		c.handleImageContent(parsedMessage)
	// 	case types.MESSAGE_TYPE_OTHERS:
	// 		c.handleTextContent(parsedMessage)
	// 	case types.MESSAGE_TYPE_IMAGE:
	// 		c.handleImageContent(parsedMessage)
	// 	default:
	// 		return
	// 	}
	// }
}

func (c *CallbackService) handleTextContent(parsedMessage *types.ParsedMessage) {
	c.Message.ReplyMessage(parsedMessage)
}

func (c *CallbackService) handleImageContent(parsedMessage *types.ParsedMessage) {
	c.Message.HandleImageContent(parsedMessage.ID, func(file *os.File) error {
		fmt.Println(file.Name())
		return nil
	})
}
