package setup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func NewPkgServices(ctx context.Context) *PkgServices {
	driveService, sheetService := getGcpService(ctx)

	lineBotService := getLineBotService()

	gptService := getGptService()

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

func (p *PkgServices) Gpt() (string, string) {
	return p.gptService.apiUrl, p.gptService.apiKey
}

func getGcpService(ctx context.Context) (*drive.Service, 
	*sheets.Service) {
	credentialsJSON := os.Getenv("CREDENTIALS_JSON")
	cwd, err := os.Getwd()
	if err != nil {
			log.Fatalf("Unable to get current working directory: %v", err)
	}
	filePath := filepath.Join(cwd, "credentials.json")
	err = os.WriteFile(filePath, []byte(credentialsJSON), 0644);
	if err != nil {
		log.Fatalf("Unable to read credentials.json %v", err)
	}
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Unable to read credentials.json %v", err)
	}

	config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope, drive.DriveScope)

	client := getGcpClient(config, ctx)

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable new google drive api service %v", err)
	}

	sheetService, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable new google sheet api service %v", err)
	}

	return driveService, sheetService
}

type GptService struct {
	apiUrl string
	apiKey string
}

func getGptService() *GptService {
	gptApiUrl := os.Getenv("GPT_API_URL")
	gptApiKey := os.Getenv("GPT_API_KEY")
	return &GptService{
		apiKey: gptApiKey,
		apiUrl: gptApiUrl,
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

// Retrieve a token, saves the token, then returns the generated client.
func getGcpClient(config *oauth2.Config, ctx context.Context) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	const TOKEN_FILE = "token.json"

	// 環境変数からtoken-jsonを取得する
	tokenJSON := os.Getenv("TOKEN_JSON")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current working directory: %v", err)
	}
	filePath := filepath.Join(cwd, TOKEN_FILE)

	// 環境変数にtoken-jsonが設定されている場合は、作業ディレクトリにtoken.jsonとして書き出す
	if (tokenJSON != "") {
		err = os.WriteFile(filePath, []byte(tokenJSON), 0644);
		if err != nil {
			log.Fatalf("Unable to write file %v", err)
		}
	}

	// token.jsonから認証情報を取得する
	tok, err := tokenFromFile(filePath)
	if err != nil {
		// token.jsonが存在しない場合は、認証後token.jsonを作成する
		tok = getTokenFromWeb(config)
		saveToken(filePath, tok)
	}

	return config.Client(ctx, tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// token.jsonから認証用トークンを取得する
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// 指定されたパス/ファイル名にトークンを保存する
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
