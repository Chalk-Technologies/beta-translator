package notion

import (
	"context"
	"github.com/Chalk-Technologies/beta-translator/internal/translation"
	"github.com/jomei/notionapi"
	"log"
)

// case "Category":
// pr := p.(*notionapi.SelectProperty)
// cat = pr.Select.Name
// case "Name":
// pr := p.(*notionapi.TitleProperty)
// //log.Printf("\n\n GOT TITLE %#v \n\n", pr)
// l = pr.Title[0].Text.Content

func AddTranslationRow(category, label, en, es, de, fr, pt, km, dbID string) error {
	log.Printf("adding row %v.%v : %v, %v, %v, %v, %v, %v", category, label, en, es, de, fr, pt, km)
	prop := notionapi.Properties{}
	prop["Name"] = notionapi.TitleProperty{
		Title: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: label,
			},
		}},
	}
	prop["Category"] = notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: category,
		},
	}
	prop["Text English"] = notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: en,
			},
		}},
	}
	prop["Text Spanish"] = notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: es,
			},
		}},
	}
	prop["Text German"] = notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: de,
			},
		}},
	}
	prop["Text French"] = notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: fr,
			},
		}},
	}
	prop["Text Portuguese"] = notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: pt,
				Link:    nil,
			},
		}},
	}
	prop["Text Khmer"] = notionapi.RichTextProperty{
		RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: km,
			},
		}},
	}

	prq := notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(dbID),
		},
		Properties: prop,
	}
	_, err := client.Page.Create(context.Background(), &prq)
	return err
}

func UpdateOrCreateTranslationRow(label, category, property, value, databaseID string, doNotOverwrite bool) error {
	// first try to get the row
	log.Printf("updating or creating translation row %v.%v %v %v\n", category, label, property, value)
	dbr := notionapi.DatabaseQueryRequest{
		Filter: notionapi.AndCompoundFilter{notionapi.PropertyFilter{Property: "Name", RichText: &notionapi.TextFilterCondition{
			Equals: label,
		}}},
	}
	query, err := client.Database.Query(context.Background(), notionapi.DatabaseID(databaseID), &dbr)
	if err != nil {
		log.Println("got error while querying")
		return err
	}
	var pageInQuestion *notionapi.Page

	for i, r := range query.Results {
		var cat string
		var l string
		for lb, p := range r.Properties {
			switch lb {
			case "Category":
				//log.Printf("got category %v\n", p)
				pr := p.(*notionapi.SelectProperty)
				cat = pr.Select.Name
			case "Name":
				//log.Printf("got name %v\n", p)
				pr := p.(*notionapi.TitleProperty)
				l = pr.Title[0].Text.Content
			}
		}
		//log.Printf("comparing %v to %v, %v to %v\n\n", cat, category, l, label)
		if cat == category && l == label {
			pageInQuestion = &query.Results[i]
			break
		}
	}

	if pageInQuestion == nil {
		log.Printf("no matching page found")

		// make a new one
		var en, es, de, fr, pt, km string
		switch property {
		case "Text English":
			en = value
		case "Text Spanish":
			es = value
		case "Text German":
			de = value
		case "Text French":
			fr = value
		case "Text Portuguese":
			pt = value
		case "Text Khmer":
			km = value
		}
		return AddTranslationRow(category, label, en, es, de, fr, pt, km, databaseID)
	}
	//log.Printf("got matching page %v\n", *pageInQuestion)

	// update the property in question
	if _, ok := pageInQuestion.Properties[property]; !ok {
		// add the prop
		pageInQuestion.Properties[property] = notionapi.RichTextProperty{RichText: []notionapi.RichText{{
			Text: &notionapi.Text{
				Content: value,
			},
		},
		}}
	} else {
		//update the content
		// check if there's a value already andif so just return, we don't want to overwrite
		if doNotOverwrite && pageInQuestion.Properties[property].(*notionapi.RichTextProperty).RichText[0].Text.Content != "" {
			return nil
		}
		pageInQuestion.Properties[property].(*notionapi.RichTextProperty).RichText[0].Text.Content = value
	}
	return UpdateTranslationRow(string(pageInQuestion.ID), pageInQuestion.Properties)
}

func UpdateTranslationRow(pageID string, properties notionapi.Properties) error {
	log.Printf("updating translation %v\n", properties)
	puq := notionapi.PageUpdateRequest{
		Properties: properties,
	}
	ctx := context.Background()
	_, err := client.Page.Update(ctx, notionapi.PageID(pageID), &puq)
	return err
}

func UpdateOrCreateAllTranslationObjects(t translation.Translation, property string, dbID string) error {
	//ctx := context.Background()
	//query, err := client.Database.Query(ctx, notionapi.DatabaseID(dbID), nil)
	//if err != nil {
	//	return err
	//}
	for category, values := range t {
		for label, val := range values {
			if err := UpdateOrCreateTranslationRow(label, category, property, val.(string), dbID, true); err != nil {
				return err
			}
		}
	}
	return nil
}
