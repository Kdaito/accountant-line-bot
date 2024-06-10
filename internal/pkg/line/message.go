package line

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
	"github.com/Kdaito/accountant-line-bot/internal/types"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type Message struct {
	ChannelSecret string
	ClientBot     *linebot.Client
	MessagingBot  *messaging_api.MessagingApiAPI
}

func NewMessage(
	clientBot *linebot.Client,
	messagingBot *messaging_api.MessagingApiAPI,
	channelSecret string,
) *Message {
	return &Message{
		ChannelSecret: channelSecret,
		ClientBot:     clientBot,
		MessagingBot:  messagingBot,
	}
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
	_, err := m.MessagingBot.ReplyMessage(&messaging_api.ReplyMessageRequest{
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

func (m *Message) HandleImageContent(messageId string, callback func(string) error) error {
	// lineのメッセージから画像を取得する
	content, err := m.ClientBot.GetMessageContent(messageId).Do()

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot get image content of message.", err)
	}

	defer content.Content.Close()

	// 取得した画像をbase64エンコードする
	encodedImage, err := m.encodeImage(content.Content)

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot encode image.", err)
	}

	err = callback(encodedImage)

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Unexpected error occured.", err)
	}

	return nil
}

func (m *Message) encodeImage(content io.ReadCloser) (string, error) {
	// 一時保存用のファイルを作成する
	file, err := os.Create("tmp-image.jpg")

	if err != nil {
		return "", fmt.Errorf("Failed create tmp image file: %s", err)
	}

	// エンコードが終了したら、ファイルをクローズして削除する
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	buf := new(bytes.Buffer)

	size, err := io.Copy(buf, content)
	if err != nil {
		return "", fmt.Errorf("Error reading content: %s", err)
	}

	if size == 0 {
		return "", fmt.Errorf("No content retrieved, possible incorrect message ID or empty file.")
	}

	_, err = io.Copy(file, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return "", fmt.Errorf("Error saving content to file: %v", err)
	}

	// Base64エンコード
	encodeResult := base64.StdEncoding.EncodeToString(buf.Bytes())

	return encodeResult, nil
}

func checkMessage(text string) string {
	switch text {
	default:
		return "すまん、今寝とる。"
	}
}
