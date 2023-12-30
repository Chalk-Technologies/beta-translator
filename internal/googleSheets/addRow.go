package googleSheets

import (
	"context"
	"google.golang.org/api/sheets/v4"
)

// add a row
func AddRow(id, category, key, english string) error {
	// we add it to the "en" sheet
	doc, err := sheetsSvc.Spreadsheets.Get(id).Context(context.Background()).Do()
	if err != nil {
		return err
	}
	enSheetId := int64(0)
	for _, sh := range doc.Sheets {
		if sh.Properties.Title == "en" {
			enSheetId = sh.Properties.SheetId
			break
		}
	}
	_, err = sheetsSvc.Spreadsheets.BatchUpdate(id,
		&sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{{AppendCells: &sheets.AppendCellsRequest{
				Fields: "*",
				Rows: []*sheets.RowData{{
					Values: []*sheets.CellData{{
						UserEnteredValue: &sheets.ExtendedValue{
							StringValue: &category,
						},
					}, {
						UserEnteredValue: &sheets.ExtendedValue{
							StringValue: &key,
						},
					}, {
						UserEnteredValue: &sheets.ExtendedValue{
							StringValue: &english,
						},
					}},
				}},
				SheetId: enSheetId,
			}}},
		}).Context(context.Background()).Do()
	return err
}
