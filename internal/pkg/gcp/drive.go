package gcp

import (
	"net/http"

	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
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
		return "", app_error.NewAppError(http.StatusInternalServerError, "Unable to upload sheet to Drive.", err)
	}

	return r.Id, nil
}

// func (d *Drive) transferOwnerShip(fileId, newOwnerEmail string) error {
// 	permission := &drive.Permission{
// 		Type:  "user",
// 		Role:  "writer",
// 		Value: newOwnerEmail,
// 	}

// 	res, err := d.Service.Permissions.Insert(fileId, permission).Do()

// 	// 新しいオーナーに編集権限を与える
// 	if err != nil {
// 		return fmt.Errorf("Unable to grant write permission: %v", err)
// 	}

// 	// 新しいオーナーに所有権を移譲する
// 	permission.Role = "owner"

// 	_, err = d.Service.Permissions.Patch(fileId, res.Id, permission).TransferOwnership(true).Do()

// 	fmt.Print("owner権限を移行しました")

// 	if err != nil {
// 		return fmt.Errorf("Unable to transfer ownership: %v", err)
// 	}

// 	return nil
// }
