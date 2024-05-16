package interfaces

import (
	"context"

	"github.com/Kdaito/accountant-line-bot/internal/types"
)

type SheetInterface interface {
	CreateSheet(ctx context.Context) (*types.SheetForDrive, error)
	WriteSheet(fileId string) error
}
