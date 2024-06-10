package gcp

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/Kdaito/accountant-line-bot/internal/lib/app_error"
	"github.com/Kdaito/accountant-line-bot/internal/types"
	"google.golang.org/api/sheets/v4"
)

type Sheet struct {
	service *sheets.Service
}

func NewSheet(service *sheets.Service) *Sheet {
	return &Sheet{
		service: service,
	}
}

func (s *Sheet) CreateSheet(ctx context.Context) (*types.SheetForDrive, error) {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	timestamp := time.Now().In(jst).Format("2006-01-02-15-04-05")

	newSpreadsheet, err := s.service.Spreadsheets.Create(&sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title:    timestamp,
			Locale:   "ja_JP",
			TimeZone: "Asia/Tokyo",
		},
	}).Do()

	if err != nil || len(newSpreadsheet.Sheets) < 1 {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Unable to create sheet", err)
	}

	fileId := newSpreadsheet.SpreadsheetId

	sheetId := newSpreadsheet.Sheets[0].Properties.SheetId

	if err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Unable to marshal spreadsheet data.", err)
	}

	return &types.SheetForDrive{Title: newSpreadsheet.Properties.Title, FileId: fileId, SheetId: sheetId}, nil
}

func (s *Sheet) WriteSheet(fileId string, sheetId int64, receipt *types.Receipt) error {
	var vr sheets.ValueRange

	// 構造体にフィールドが設定されている場合は、シートに反映する
	receptReflect := reflect.TypeOf(*receipt)

	const UNDEFINED_VALUE = "undefined"

	// 日付
	date := UNDEFINED_VALUE
	if _, bol := receptReflect.FieldByName("Date"); bol {
		date = receipt.Date
	}
	vr.Values = append(vr.Values, []interface{}{"date", "", date})

	// 商品リスト
	itemRowItem := []interface{}{"name", "count", "amount"}
	vr.Values = append(vr.Values, itemRowItem)

	for _, item := range receipt.Items {
		var name, count, amount string

		itemReflect := reflect.TypeOf(item)
		if _, bol := itemReflect.FieldByName("Name"); bol {
			name = item.Name
		} else {
			name = UNDEFINED_VALUE
		}
		if _, bol := itemReflect.FieldByName("Amount"); bol {
			count = string(item.Count)
		} else {
			count = UNDEFINED_VALUE
		}
		if _, bol := itemReflect.FieldByName("Count"); bol {
			amount = string(item.Amount)
		} else {
			amount = UNDEFINED_VALUE
		}

		itemRow := []interface{}{name, count, amount}
		vr.Values = append(vr.Values, itemRow)
	}

	// 税抜価格
	totalAmount := UNDEFINED_VALUE
	if _, bol := receptReflect.FieldByName("TotalAmount"); bol {
		totalAmount = string(receipt.TotalAmount)
	}
	vr.Values = append(vr.Values, []interface{}{"total amount", "", totalAmount})

	// 税込価格
	totalAmountIncludingTax := UNDEFINED_VALUE
	if _, bol := receptReflect.FieldByName("TotalAmountIncludingTax"); bol {
		totalAmountIncludingTax = string(receipt.TotalAmountIncludingTax)
	}
	vr.Values = append(vr.Values, []interface{}{"total amount (including tax)", "", totalAmountIncludingTax})

	// 通貨
	currencySymbol := UNDEFINED_VALUE
	if _, bol := receptReflect.FieldByName("CurrencySymbol"); bol {
		currencySymbol = receipt.CurrencySymbol
	}
	vr.Values = append(vr.Values, []interface{}{"currency symbol", "", currencySymbol})

	writeRange := "B2"

	_, err := s.service.Spreadsheets.Values.Update(fileId, writeRange, &vr).ValueInputOption("RAW").Do()

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot wirte spread sheet values.", err)
	}

	const LEN_TO_ITEM_TABLE = 3
	itemLength := len(receipt.Items)

	// スプレッドシートのスタイリング
	// セルの結合リクエスト
	margeDateCellsRequest := &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartRowIndex:    1,
				EndRowIndex:      2,
				StartColumnIndex: 1,
				EndColumnIndex:   3,
			},
			MergeType: "MERGE_ALL",
		},
	}
	margeTotalAmountCellsRequest := &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartRowIndex:    int64(LEN_TO_ITEM_TABLE + itemLength),
				EndRowIndex:      int64(LEN_TO_ITEM_TABLE + itemLength + 1),
				StartColumnIndex: 1,
				EndColumnIndex:   3,
			},
			MergeType: "MERGE_ALL",
		},
	}
	margeTotalAmountIncludingTaxCellsRequest := &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartRowIndex:    int64(LEN_TO_ITEM_TABLE + itemLength + 1),
				EndRowIndex:      int64(LEN_TO_ITEM_TABLE + itemLength + 2),
				StartColumnIndex: 1,
				EndColumnIndex:   3,
			},
			MergeType: "MERGE_ALL",
		},
	}
	margeCurrencySymbolCellsRequest := &sheets.Request{
		MergeCells: &sheets.MergeCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartRowIndex:    int64(LEN_TO_ITEM_TABLE + itemLength + 2),
				EndRowIndex:      int64(LEN_TO_ITEM_TABLE + itemLength + 3),
				StartColumnIndex: 1,
				EndColumnIndex:   3,
			},
			MergeType: "MERGE_ALL",
		},
	}

	// 背景色の設定リクエスト
	updateCellsRequest := &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartRowIndex:    2,
				EndRowIndex:      3,
				StartColumnIndex: 1,
				EndColumnIndex:   4,
			},
			Fields: "userEnteredFormat.backgroundColor",
			Rows: []*sheets.RowData{
				{
					Values: []*sheets.CellData{
						{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.5,
									Green: 0.8,
									Blue:  1.0,
								},
							},
						},
						{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.5,
									Green: 0.8,
									Blue:  1.0,
								},
							},
						},
						{
							UserEnteredFormat: &sheets.CellFormat{
								BackgroundColor: &sheets.Color{
									Red:   0.5,
									Green: 0.8,
									Blue:  1.0,
								},
							},
						},
					},
				},
			},
		},
	}

	// 枠線の設定リクエスト
	updateBordersRequest := &sheets.Request{
		UpdateBorders: &sheets.UpdateBordersRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetId,
				StartRowIndex:    1,
				EndRowIndex:      int64(LEN_TO_ITEM_TABLE + itemLength + 3),
				StartColumnIndex: 1,
				EndColumnIndex:   4,
			},
			Top: &sheets.Border{
				Style: "SOLID",
				Color: &sheets.Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
				},
			},
			Bottom: &sheets.Border{
				Style: "SOLID",
				Color: &sheets.Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
				},
			},
			Left: &sheets.Border{
				Style: "SOLID",
				Color: &sheets.Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
				},
			},
			Right: &sheets.Border{
				Style: "SOLID",
				Color: &sheets.Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
				},
			},
			InnerHorizontal: &sheets.Border{
				Style: "SOLID",
				Color: &sheets.Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
				},
			},
			InnerVertical: &sheets.Border{
				Style: "SOLID",
				Color: &sheets.Color{
					Red:   0.0,
					Green: 0.0,
					Blue:  0.0,
				},
			},
		},
	}

	// リクエストのバッチ実行
	requests := []*sheets.Request{margeDateCellsRequest, margeTotalAmountCellsRequest, margeTotalAmountIncludingTaxCellsRequest, margeCurrencySymbolCellsRequest, updateCellsRequest, updateBordersRequest}
	batchUpdate := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err = s.service.Spreadsheets.BatchUpdate(fileId, batchUpdate).Do()
	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot batch update spread sheet for styling", err)
	}

	return nil
}
