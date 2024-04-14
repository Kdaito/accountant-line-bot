package module

import (
	"context"
	"log"

	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

type DriveService struct {
	service *drive.Service
}

// サービスアカウントのクライアントを作成する
func NewDriveService(ctx context.Context, b []byte) *DriveService {
	// サービスアカウントのクライアントを作成する
	srv, err := drive.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		log.Fatalf("cannot create client from service account: %v", err)
	}
	return &DriveService{
		service: srv,
	}
}