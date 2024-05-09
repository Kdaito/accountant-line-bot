package pkg

import (
	"bytes"
	"fmt"

	"github.com/Kdaito/accountant-line-bot/internal/types"
	"google.golang.org/api/drive/v2"
)

type GDrive struct {
	Service *drive.Service
}

func (g *GDrive) Upload(parentId string, sheetForDrive *types.SheetForDrive) (string, error) {
	f := &drive.File{
		Title:    sheetForDrive.Title,
		Parents:  []*drive.ParentReference{{Id: parentId}},
		MimeType: "application/vnd.google-apps.spreadsheet",
	}

	r, err := g.Service.Files.Insert(f).Media(bytes.NewReader(sheetForDrive.ByteFile)).Do()
	if err != nil {
		return "", fmt.Errorf("Unable to upload sheet to Drive: %v", err)
	}

	return r.Id, nil
}
