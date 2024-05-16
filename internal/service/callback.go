package service

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/interfaces"
	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
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

	err = c.Sheet.WriteSheet(sheetForDrive.FileId)

	if err != nil {
		fmt.Printf("%v", err)
	}


	targetFolderId := os.Getenv("DRIVE_FOLDER_ID")
	c.Drive.Move(targetFolderId, sheetForDrive)

	if err != nil {
		fmt.Printf("%v", err)
	}

	// parsedMessages, err := c.Message.ParseRequest(w, req)

	// if err != nil {
	// 	c.setErrorResponse(err, w)
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

func (c *CallbackService) setErrorResponse(err error, w http.ResponseWriter) {
	var AppErrorType *app_error.AppError
	if errors.As(err, &AppErrorType) {
		w.WriteHeader(err.(*app_error.AppError).Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write([]byte(err.Error()))
	return
}
