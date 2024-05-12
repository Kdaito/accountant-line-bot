package pkg

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
	"github.com/Kdaito/accountant-line-bot/internal/types"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Message struct {
	ChannelSecret string
	Bot           *messaging_api.MessagingApiAPI
	Blob          *messaging_api.MessagingApiBlobAPI
}

func (m *Message) ParseRequest(w http.ResponseWriter, req *http.Request) ([]*types.ParsedMessage, error) {
	res := []*types.ParsedMessage{}

	cb, err := webhook.ParseRequest(m.ChannelSecret, req)
	if err != nil {
		if errors.Is(err, webhook.ErrInvalidSignature) {
			return nil, app_error.NewAppError(http.StatusBadRequest, "Cannnot parse request. Invalid signature.", err)
		}
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Cannnot parse request.", err)
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				res = append(res, &types.ParsedMessage{
					MessageType: types.MESSAGE_TYPE_TEXT,
					Text:        checkMessage(message.Text),
					ID:          message.Id,
					ReplyToken:  e.ReplyToken,
				})
			case webhook.ImageMessageContent:
				res = append(res, &types.ParsedMessage{
					MessageType: types.MESSAGE_TYPE_IMAGE,
					Text:        "",
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

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannnot reply message.", err)
	}

	return nil
}

func (m *Message) HandleImageContent(messageId string, callback func(*os.File) error) error {
	content, err := m.Blob.GetMessageContent(messageId)

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot get image content of message.", err)
	}

	defer content.Body.Close()

	file, err := m.SaveTmpImage(content.Body)

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot save temporary image.", err)
	}

	defer os.Remove(file.Name())

	err = callback(file)

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Unexpected error occured.", err)
	}

	return nil
}

func (m *Message) SaveTmpImage(content io.ReadCloser) (*os.File, error) {
	file, err := os.Create("tmp-image")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func checkMessage(text string) string {
	switch text {
	default:
		return "すまん、今寝とる。"
	}
}
