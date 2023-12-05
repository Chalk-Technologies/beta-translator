package notion

import (
	"context"
	"github.com/Chalk-Technologies/beta-translator/internal/translation"
	"github.com/jomei/notionapi"
)

var client *notionapi.Client

// todo abstract and allow public
func Init(secret string) {
	client = notionapi.NewClient(notionapi.Token(secret))
	return
}

// get the translations from the notion table
func GetTranslations(databaseID string) (translation.Translation, translation.Translation, translation.Translation, translation.Translation, translation.Translation, translation.Translation, translation.Translation, error) {
	ctx := context.Background()
	query, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), nil)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	tEN := make(translation.Translation)
	tES := make(translation.Translation)
	tDE := make(translation.Translation)
	tFR := make(translation.Translation)
	tPT := make(translation.Translation)
	tKM := make(translation.Translation)
	tKO := make(translation.Translation)

	next := true

	for next {
		for _, r := range query.Results {
			// do we really have to query the api again to get the property vals  how annoying!
			var cat string
			var l string
			var en string
			var es string
			var fr string
			var de string
			var pt string
			var km string
			var ko string

			//r.Properties.UnmarshalJSON()
			for label, p := range r.Properties {
				switch label {
				case "Category":
					pr := p.(*notionapi.SelectProperty)
					cat = pr.Select.Name
				case "Name":
					pr := p.(*notionapi.TitleProperty)
					//log.Printf("\n\n GOT TITLE %#v \n\n", pr)
					l = pr.Title[0].Text.Content
				case "Text English":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {
						en = pr.RichText[0].PlainText
					}
				case "Text Spanish":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {

						es = pr.RichText[0].PlainText
					}
				case "Text German":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {

						de = pr.RichText[0].PlainText
					}
				case "Text French":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {

						fr = pr.RichText[0].PlainText
					}
				case "Text Portuguese":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {

						pt = pr.RichText[0].PlainText
					}
				case "Text Khmer":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {

						km = pr.RichText[0].PlainText
					}
				case "Text Korean":
					pr := p.(*notionapi.RichTextProperty)
					if len(pr.RichText) > 0 {

						ko = pr.RichText[0].PlainText
					}
				}
			}
			tEN.AddValue(cat, l, en)
			tES.AddValue(cat, l, es)
			tDE.AddValue(cat, l, de)
			tFR.AddValue(cat, l, fr)
			tPT.AddValue(cat, l, pt)
			tKM.AddValue(cat, l, km)
			tKO.AddValue(cat, l, ko)

		}

		if query.HasMore {
			dbqr := &notionapi.DatabaseQueryRequest{
				StartCursor: query.NextCursor,
			}
			query, err = client.Database.Query(ctx, notionapi.DatabaseID(databaseID), dbqr)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, err
			}
		} else {
			next = false
		}
	}

	return tEN, tES, tFR, tDE, tPT, tKM, tKO, nil
}
