package pkg

import (
	"fmt"

	"github.com/Kdaito/accountant-line-bot/internal/types"
	"google.golang.org/api/drive/v2"
)

type Drive struct {
	Service *drive.Service
}

func NewDrive(service *drive.Service) *Drive {
	return &Drive{
		Service: service,
	}
}

func (d *Drive) Move(folderId string, sheetForDrive *types.SheetForDrive) (string, error) {
	f := &drive.File{
		Title:    sheetForDrive.Title,
		Parents:  []*drive.ParentReference{{Id: folderId}},
		MimeType: "application/vnd.google-apps.spreadsheet",
	}

	r, err := d.Service.Files.Update(sheetForDrive.FileId, f).Do()

	if err != nil {
		return "", fmt.Errorf("Unable to upload sheet to Drive: %v", err)
	}

	return r.Id, nil
}

