package googleSheets

import (
	"context"
	"github.com/Chalk-Technologies/beta-translator/internal/translation"
	"google.golang.org/api/sheets/v4"
)

var sheetsSvc *sheets.Service

func Init() error {
	ctx := context.Background()
	var err error
	sheetsSvc, err = sheets.NewService(ctx) //option.WithCredentialsFile("beta-29103-74da29af10d8.json"))
	return err
}

func getTranslationsDoc(id string) (*sheets.Spreadsheet, error) {
	r, err := sheetsSvc.Spreadsheets.Get(id).Context(context.Background()).Do()
	return r, err
}

var langs = []string{"en", "es", "fr", "de", "pt", "km", "ko"}

// get the translations from the notion table
func GetTranslations(id string) (translation.Translation, translation.Translation, translation.Translation, translation.Translation, translation.Translation, translation.Translation, translation.Translation, error) {
	//doc, err := getTranslationsDoc(id)
	//if err != nil {
	//	return nil, nil, nil, nil, nil, nil, nil, err
	//}

	translations := make([]translation.Translation, 7)
	for i := range langs {
		// get the sheet
		valuesService := sheets.NewSpreadsheetsValuesService(sheetsSvc)
		//var values *sheets.ValueRange
		values, err := valuesService.Get(id, langs[i]).MajorDimension("ROWS").Do()
		if err != nil {
			return nil, nil, nil, nil, nil, nil, nil, err
		}
		// get the translations
		trans := make(translation.Translation)
		// first col is category
		// second col is key
		// third col is english
		// fourth col is auto-translated
		// fifth col is override

		for _, row := range values.Values[1:] {
			var category, key, t string
			var skip bool
			for ii, val := range row {
				var ok bool
				switch ii {
				case 0:
					// category
					category, ok = val.(string)
				case 1:
					// key
					key, ok = val.(string)
				default:
					// start with english, override with auto-translate, override with manual translate
					t, ok = val.(string)
				}
				if !ok {
					skip = true
					break
				}
			}
			if skip || category == "" || key == "" || t == "" || t == "#VALUE!" {
				continue
			}
			trans.AddValue(category, key, t)
		}

		// then add it to the list
		translations[i] = trans
	}

	return translations[0], translations[1], translations[2], translations[3], translations[4], translations[5], translations[6], nil
}
