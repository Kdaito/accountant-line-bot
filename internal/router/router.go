package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/gdrive"
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
	_, err := messaging_api.NewMessagingApiAPI(
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
	drive := &gdrive.GDrive{Service: driveService}

	callbackService := &service.CallbackService{
		Drive: drive,
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
