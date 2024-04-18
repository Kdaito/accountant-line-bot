package module

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	service *sheets.Service
}

// サービスアカウントのクライアントを作成する
func NewSheetsService(ctx context.Context, b []byte) *SheetsService {
	// サービスアカウントのクライアントを作成する
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("cannot create client from service account: %v", err)
	}
	return &SheetsService{
		service: srv,
	}
}
