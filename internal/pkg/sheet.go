package pkg

import (
	"context"
	"net/http"
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
	}).Context(ctx).Do()

	if err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Unable to create sheet", err)
	}

	spreadSheetFileId := newSpreadsheet.SpreadsheetId

	if err != nil {
		return nil, app_error.NewAppError(http.StatusInternalServerError, "Unable to marshal spreadsheet data.", err)
	}

	return &types.SheetForDrive{Title: newSpreadsheet.Properties.Title, FileId: spreadSheetFileId}, nil
}

func (s *Sheet) WriteSheet(fileId string) error {
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{"Name", "Age", "Department"})
	for _, emp := range getSampleEmployee() {
		row := []interface{}{emp.Name, emp.Age, emp.Department}
		vr.Values = append(vr.Values, row)
	}

	writeRange := "A1"

	_, err := s.service.Spreadsheets.Values.Update(fileId, writeRange, &vr).ValueInputOption("RAW").Do()

	if err != nil {
		return app_error.NewAppError(http.StatusInternalServerError, "Cannot wirte spread sheet values.", err)
	}

	return nil
}

// ============ for sample ============

type Employee struct {
	Name       string `json:"Name"`
	Age        int    `json:"Age"`
	Department string `json:"Department"`
}

func getSampleEmployee() []Employee {
	return []Employee{{Name: "Hiroto", Age: 22, Department: "Development"}, {Name: "Hayato", Age: 21, Department: "Manager"}}
}
