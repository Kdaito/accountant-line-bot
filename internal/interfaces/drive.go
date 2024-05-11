package interfaces

import "github.com/Kdaito/accountant-line-bot/internal/types"

type DriveInterface interface {
	Move(parentId string, sheetForDrive *types.SheetForDrive) (string, error)
}
