package router

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/pkg"
	"github.com/Kdaito/accountant-line-bot/internal/service"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Router struct {
	Port string
}

func (r *Router) Set(channelSecret, channelToken, gptApiUrl, gptApiKey string) {
	// port setting
	if r.Port == "" {
		r.Port = "2001"
	}

	ctx := context.Background()

	clientBot, err := linebot.New(channelSecret, channelToken)

	// messaging api setting
	messagingBot, err := messaging_api.NewMessagingApiAPI(
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
		log.Fatalf("cannot init drive service with credentials json: %v", err)
	}

	// google sheet api setting
	sheetService, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("cannot init sheet service with credentials json: %v", err)
	}

	// DI
	drivePkg := pkg.NewDrive(driveService)
	sheetPkg := pkg.NewSheet(sheetService)
	chatAIPkg := pkg.NewChatAI(gptApiUrl, gptApiKey)
	messagePkg := &pkg.Message{ChannelSecret: channelSecret, ClientBot: clientBot, MessagingBot: messagingBot}

	callbackService := &service.CallbackService{
		Drive:   drivePkg,
		Sheet:   sheetPkg,
		Message: messagePkg,
		ChatAI:  chatAIPkg,
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
