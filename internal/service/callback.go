package service

import (
	"context"
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
	ChatAI  interfaces.ChatAIInterface
}

func (c *CallbackService) Callback(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	parsedMessages, err := c.Message.ParseRequest(w, req)

	if err != nil {
		c.setErrorResponse(err, w)
		return
	}

	for _, parsedMessage := range parsedMessages {
		switch parsedMessage.MessageType {
		case types.MESSAGE_TYPE_TEXT:
			c.handleTextContent(parsedMessage)
			return
		case types.MESSAGE_TYPE_OTHERS:
			c.handleTextContent(parsedMessage)
			return
		case types.MESSAGE_TYPE_IMAGE:
			err = c.handleImageContent(parsedMessage, ctx)
			if err != nil {
				c.setErrorResponse(err, w)
			}
			return
		default:
			return
		}
	}
}

func (c *CallbackService) handleTextContent(parsedMessage *types.ParsedMessage) {
	// res, _ := c.ChatAI.ScanReceipt("")
	// fmt.Printf("result struct: %#v\n", res)
	// c.Message.ReplyMessage(parsedMessage)
	return
}

func (c *CallbackService) handleImageContent(parsedMessage *types.ParsedMessage, ctx context.Context) error {
	return c.Message.HandleImageContent(parsedMessage.ID, func(encodedImage string) error {

		receipt, err := c.ChatAI.ScanReceipt(encodedImage)

		if err != nil {
			return err
		}

		sheetForDrive, err := c.Sheet.CreateSheet(ctx)

		if err != nil {
			return err
		}

		err = c.Sheet.WriteSheet(sheetForDrive.FileId, sheetForDrive.SheetId, receipt)

		if err != nil {
			return err
		}

		targetFolderId := os.Getenv("DRIVE_FOLDER_ID")

		_, err = c.Drive.Move(targetFolderId, sheetForDrive)

		if err != nil {
			return err
		}

		return nil
	})
}

func (c *CallbackService) setErrorResponse(err error, w http.ResponseWriter) {
	var AppErrorType *app_error.AppError

	// エラーログをコンソールに出力するため
	fmt.Print(err)

	if errors.As(err, &AppErrorType) {
		w.WriteHeader(err.(*app_error.AppError).Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(err.Error()))

	return
}
