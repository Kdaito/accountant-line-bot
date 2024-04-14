package processer

import (
	"errors"
	// "context"
	// "io/ioutil"
	"log"
	"net/http"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	// "github.com/Kdaito/accountant-line-bot/module"
)

type Processer struct {
	channelSecret string
	bot           *messaging_api.MessagingApiAPI
}

func NewProcesser(channelSecret string, bot *messaging_api.MessagingApiAPI) *Processer {
	return &Processer{
		channelSecret: channelSecret,
		bot:           bot,
	}
}

func (p *Processer) ParseRequest(w http.ResponseWriter, req *http.Request) *webhook.CallbackRequest {
	cb, err := webhook.ParseRequest(p.channelSecret, req)
	if err != nil {
		log.Printf("Cannot parse request: %+v\n", err)
		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return nil
	}
	return cb
}

func (p *Processer) AnalyzeTextMessage(message webhook.TextMessageContent) string {
	text := message.Text

	if text == "おはよう" {
		return "おはよう"
	} else {
		return "は？"
	}
}

func (p *Processer) ReplyMessage(message string, replyToken string) error {
	_, err := p.bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
		ReplyToken: replyToken,
		Messages: []messaging_api.MessageInterface{
			messaging_api.TextMessage{
				Text: message,
			},
		},
	})

	return err
}

func (p *Processer) ExportSheet() {
	// ctx := context.Background()

	// // サービスアカウントの秘密鍵を読み込む
	// b, err := ioutil.ReadFile("service-account.json")
	// if err != nil {
	// 	log.Fatalf("cannot read service account json file: %v", err)
	// }

	// サービスアカウントのクライアントを作成する
	// driveSrv := module.NewDriveService(ctx, b)
	// sheetsSrv := module.NewSheetsService(ctx, b)
}
