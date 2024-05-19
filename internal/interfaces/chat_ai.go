package interfaces

import "github.com/Kdaito/accountant-line-bot/internal/types"

type ChatAIInterface interface {
	ScanReceipt(string) (*types.Receipt, error)
}
