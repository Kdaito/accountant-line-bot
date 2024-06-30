package setup

import (
	"context"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type PkgServices struct {
	driveService   *drive.Service
	sheetService   *sheets.Service
	lineBotService *LineBotService
	gptService     *GptService
}

func NewPkgServices(ctx context.Context, isSkipGpt bool) *PkgServices {
	driveService, sheetService := getGcpService(ctx)

	lineBotService := getLineBotService()

	gptService := getGptService(isSkipGpt)

	return &PkgServices{
		driveService:   driveService,
		sheetService:   sheetService,
		lineBotService: lineBotService,
		gptService:     gptService,
	}
}

func (p *PkgServices) Drive() *drive.Service {
	return p.driveService
}

func (p *PkgServices) Sheet() *sheets.Service {
	return p.sheetService
}

func (p *PkgServices) LineBot() (*linebot.Client, *messaging_api.MessagingApiAPI, string) {
	return p.lineBotService.lineBotClient, p.lineBotService.messagingApi, p.lineBotService.channelSecret
}

func (p *PkgServices) Gpt() (string, string, bool) {
	return p.gptService.apiUrl, p.gptService.apiKey, p.gptService.isSkipGpt
}

func getGcpService(ctx context.Context) (*drive.Service,
	*sheets.Service) {
	// 環境変数からシークレットを取得
	serviceAccountJSON := os.Getenv("SERVICE_ACCOUNT_JSON")

	// シークレットをファイルに書き込む
	err := os.WriteFile("service-account.json", []byte(serviceAccountJSON), 0644)

	if err != nil {
		log.Fatalf("Unable to write file %v", err)
	}

	// credentials.jsonファイルを読み込む
	b, err := os.ReadFile("service-account.json")

	if err != nil {
		log.Fatalf("Unable to write read service-account.json %v", err)
	}

	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("Unable new google drive api service %v", err)
	}

	sheetService, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("Unable new google sheet api service %v", err)
	}

	return driveService, sheetService
}

type GptService struct {
	apiUrl    string
	apiKey    string
	isSkipGpt bool
}

func getGptService(isSkipGpt bool) *GptService {
	gptApiUrl := os.Getenv("GPT_API_URL")
	gptApiKey := os.Getenv("GPT_API_KEY")
	return &GptService{
		apiKey:    gptApiKey,
		apiUrl:    gptApiUrl,
		isSkipGpt: isSkipGpt,
	}
}

type LineBotService struct {
	lineBotClient *linebot.Client
	messagingApi  *messaging_api.MessagingApiAPI
	channelSecret string
}

func getLineBotService() *LineBotService {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")

	lineBotClient, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		log.Fatalf("Unable to init line bot client: %v", err)
	}

	messagingApi, err := messaging_api.NewMessagingApiAPI(
		channelToken,
	)
	if err != nil {
		log.Fatalf("Unable to init line messaging api client: %v", err)
	}

	return &LineBotService{
		lineBotClient: lineBotClient,
		messagingApi:  messagingApi,
		channelSecret: channelSecret,
	}
}
