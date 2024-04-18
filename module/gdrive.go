package module

import (
	"context"
	"log"
	"os"

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

func (d *DriveService) Upload(parentId string, title string, file *os.File) (string, error) {
	f := &drive.File{
		Title: title,
		Parents: []*drive.ParentReference{{Id: parentId}},
	}

	r, err := d.service.Files.Insert(f).Media(file).Do()
	if err != nil {
		return "", err
	}

	return r.Id, nil
}
