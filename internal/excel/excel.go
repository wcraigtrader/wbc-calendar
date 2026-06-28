package excel

import (
	"fmt"
	"log"
	"time"

	"wbc-calendar/internal/event"

	"github.com/xuri/excelize/v2"
)

type ExcelReader struct {
	File     *excelize.File
	Headings *event.Headings
}

func NewExcelReader(filePath string) (*ExcelReader, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}

	return &ExcelReader{File: f}, nil
}

func (er *ExcelReader) Close() error {
	if er.File != nil {
		return er.File.Close()
	}
	return nil
}

func (er *ExcelReader) GetSheetNames() []string {
	return er.File.GetSheetList()
}

func (er *ExcelReader) sheetExists(sheetName string) bool {
	sheets := er.File.GetSheetList()
	for _, sheet := range sheets {
		if sheet == sheetName {
			return true
		}
	}
	return false
}

func (er *ExcelReader) ReadSheet(sheetName string, zone *time.Location, year int) ([]*event.Event, error) {
	if !er.sheetExists(sheetName) {
		availableSheets := er.File.GetSheetList()
		return nil, fmt.Errorf("sheet '%s' not found. Available sheets: %v", sheetName, availableSheets)
	}

	rows, err := er.File.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet '%s': %w", sheetName, err)
	}

	var result []*event.Event

	for i, row := range rows {
		if er.Headings == nil {
			if len(row) > 0 && row[0] == "Date" {
				if er.Headings, err = event.CreateHeadings(row); err != nil {
					return nil, fmt.Errorf("invalid header row in sheet '%s': %w", sheetName, err)
				}
			}
			continue
		}

		if event, err := event.NewEvent(er.Headings, i+1, row, zone, year); err != nil {
			log.Printf("Error parsing row %d in sheet '%s': %v", i+1, sheetName, err)
		} else {
			result = append(result, event)
			if event.HasErrors() {
				log.Printf("Row %4d %s: %s\n", event.Line, event, event.ErrorString())
			}
		}
	}

	return result, nil
}
