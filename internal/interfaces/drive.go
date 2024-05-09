package interfaces

import "github.com/Kdaito/accountant-line-bot/internal/types"

type DriveInterface interface {
	Upload(parentId string, sheetForDrive *types.SheetForDrive) (string, error)
}
