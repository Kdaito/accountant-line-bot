package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/pkg"
	"github.com/Kdaito/accountant-line-bot/internal/service"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

type Router struct {
	Port string
}

func (r *Router) Set(channelSecret string, channelToken string) {
	// port setting
	if r.Port == "" {
		r.Port = "2001"
	}

	ctx := context.Background()

	// messaging api setting
	bot, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)
	blob, err := messaging_api.NewMessagingApiBlobAPI(
		channelToken,
	)
	if err != nil {
		log.Fatal(err)
	}

	// google drive api setting
	b, err := ioutil.ReadFile("service-account.json")
	if err != nil {
		log.Fatalf("cannot read service account json file: %v", err)
	}

	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON(b))

	if err != nil {
		log.Fatalf("cannot create client from service account: %v", err)
	}

	// google sheet api setting

	// DI
	drive := &pkg.GDrive{Service: driveService}
	message := &pkg.Message{ChannelSecret: channelSecret, Bot: bot, Blob: blob}

	callbackService := &service.CallbackService{
		Drive:   drive,
		Message: message,
	}

	// set routing
	http.HandleFunc("/callback", callbackService.Callback)

}

func (r *Router) Run() {
	fmt.Println("http://localhost:" + r.Port + "/")
	if err := http.ListenAndServe(":"+r.Port, nil); err != nil {
		log.Fatal(err)
	}
}
