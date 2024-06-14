package service

import (
	"context"
	"errors"
	"log"
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

func NewCallbackService(
	drive interfaces.DriveInterface,
	message interfaces.MessageInterface,
	sheet interfaces.SheetInterface,
	chatAI interfaces.ChatAIInterface,
) *CallbackService {
	return &CallbackService{
		Drive:   drive,
		Message: message,
		Sheet:   sheet,
		ChatAI:  chatAI,
	}
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
		//
		// 受け取ったメッセージがテキストだった場合の処理
		//
		case types.MESSAGE_TYPE_TEXT:
			err = c.handleTextContent(parsedMessage)

			if err != nil {
				c.setErrorResponse(err, w)
				parsedMessage.Text = "なんか処理失敗したわ。もう一回送ってみてくれる？"
				c.Message.ReplyMessage(parsedMessage)
				return
			}

			c.setSuccessResponse(w)
			return
		//
		// 受け取ったメッセージが画像だった場合の処理
		//
		case types.MESSAGE_TYPE_IMAGE:
			err = c.handleImageContent(parsedMessage, ctx)

			if err != nil {
				c.setErrorResponse(err, w)
				parsedMessage.Text = "なんか処理失敗したわ。もう一回送ってみてくれる？"
				c.Message.ReplyMessage(parsedMessage)
				return
			}

			c.setSuccessResponse(w)
			return
		//
		// 受け取ったメッセージがテキストでも画像でもなかった場合の処理
		//
		case types.MESSAGE_TYPE_OTHERS:
			err = c.handleTextContent(parsedMessage)

			if err != nil {
				c.setErrorResponse(err, w)
				parsedMessage.Text = "なんか処理失敗したわ。もう一回送ってみてくれる？"
				c.Message.ReplyMessage(parsedMessage)
				return
			}

			c.setSuccessResponse(w)
			return
		default:
			return
		}
	}
}

func (c *CallbackService) handleTextContent(parsedMessage *types.ParsedMessage) error {
	parsedMessage.Text = "すまんな、私も忙しいんだ。さっさとレシートを送りたまえ。"
	if err := c.Message.ReplyMessage(parsedMessage); err != nil {
		return err
	}
	return nil
}

func (c *CallbackService) handleImageContent(parsedMessage *types.ParsedMessage, ctx context.Context) error {
	return c.Message.HandleImageContent(parsedMessage.ID, func(encodedImage string) error {

		receipt, err := c.ChatAI.ScanReceipt(encodedImage)

		if err != nil {
			return err
		}

		// 送られてきた画像がレシート画像でない場合はメッセージを送信して終了する
		if !receipt.IsReceipt {
			parsedMessage.Text = "この画像、レシートじゃないやろ。解析できんわ。"
			if err = c.Message.ReplyMessage(parsedMessage); err != nil {
				return err
			}
			return nil
		}

		sheetForDrive, err := c.Sheet.CreateSheet(ctx)

		if err != nil {
			return err
		}

		err = c.Sheet.WriteSheet(sheetForDrive.FileId, sheetForDrive.SheetId, receipt)

		if err != nil {
			return err
		}

		targetFolderId := os.Getenv("FOLDER_ID")

		_, err = c.Drive.Move(targetFolderId, sheetForDrive)

		if err != nil {
			return err
		}

		parsedMessage.Text = "解析終わったで。"
		if err = c.Message.ReplyMessage(parsedMessage); err != nil {
			return err
		}

		return nil
	})
}

func (c *CallbackService) setSuccessResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request successfully processed!"))
}

func (c *CallbackService) setErrorResponse(err error, w http.ResponseWriter) {
	var AppErrorType *app_error.AppError

	// エラーログをコンソールに出力するため
	log.Printf("エラーが発生しました。: %v", err)

	if errors.As(err, &AppErrorType) {
		w.WriteHeader(err.(*app_error.AppError).Code)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(err.Error()))

	return
}
