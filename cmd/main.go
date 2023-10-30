package main

import (
	"flag"
	"github.com/Chalk-Technologies/beta-translator/internal/notion"
	"github.com/Chalk-Technologies/beta-translator/internal/repoUpload"
	"github.com/Chalk-Technologies/beta-translator/internal/translation"
	"log"
)

func main() {
	// our main server!
	// we are going to use the notion api to get the values from our translations table
	// then we will create json files to handle the translations in our react app
	// get flags
	var export = flag.Bool("export", false, "export all translations from the notion DB")
	var importFile = flag.String("importFile", "", "import all translations file")
	var importField = flag.String("importField", "Text English", "import all translations into field")
	var notionSecret = flag.String("notionSecret", "", "notion secret")
	var notionDB = flag.String("db", "", "notion database id")
	var uploadRepo = flag.String("uploadRepo", "Chalk-Technologies/beta-web", "The repo to upload to, including repo owner and without the .git extension, not case-sensitive")
	var uploadPath = flag.String("uploadPath", "public/translations/", "The repo to upload to, including repo owner and without the .git extension, not case-sensitive")
	// todo add another flag for adding a new translation, need category, title, value
	//var nFlag = flag.Int("n", 1234, "help message for flag n")

	flag.Parse()

	//init the notion client
	// todo check for nil value in notionsecret

	//log.Println(*notionSecret)
	notion.Init(*notionSecret)

	tEN, tES, tFR, tDE, tPT, tKM, err := notion.GetTranslations(*notionDB)
	if err != nil {
		log.Fatalf("got error in translation %v", err)
		return
	}

	if importFile != nil && *importFile != "" {
		if importField == nil {
			log.Fatalf("importField must be set")
			return
		}
		t, err := translation.ImportFromFile(*importFile)
		if err != nil {
			log.Fatalf("got error importing translation %v", err)
			return
		}
		err = notion.UpdateOrCreateAllTranslationObjects(*t, *importField, *notionDB)
		if err != nil {
			log.Fatalf("got error updating or creating translation %v", err)
			return
		}
	}

	if uploadPath != nil && uploadRepo != nil && *uploadRepo != "" && *uploadPath != "" {
		repoUpload.Init()
		if err = repoUpload.UploadFile("en.json", *uploadRepo, *uploadPath, tEN); err != nil {
			log.Fatalf("got error on en.json upload: %v", err)
		}
		if err = repoUpload.UploadFile("es.json", *uploadRepo, *uploadPath, tES); err != nil {
			log.Fatalf("got error on es.json upload: %v", err)
		}
		if err = repoUpload.UploadFile("de.json", *uploadRepo, *uploadPath, tDE); err != nil {
			log.Fatalf("got error on de.json upload: %v", err)
		}
		if err = repoUpload.UploadFile("pt.json", *uploadRepo, *uploadPath, tPT); err != nil {
			log.Fatalf("got error on pt.json upload: %v", err)
		}
		if err = repoUpload.UploadFile("fr.json", *uploadRepo, *uploadPath, tFR); err != nil {
			log.Fatalf("got error on fr.json upload: %v", err)
		}
		//if err = repoUpload.UploadFile("km.json", *uploadRepo, *uploadPath, tKM); err != nil {
		//	log.Fatalf("got error on km.json upload: %v", err)
		//}
	}
	if export != nil && *export {
		if err = tEN.Export("en.json"); err != nil {
			log.Fatalf("got error on en.json export: %v", err)
		}
		if err = tES.Export("es.json"); err != nil {
			log.Fatalf("got error on es.json export: %v", err)
		}
		if err = tFR.Export("fr.json"); err != nil {
			log.Fatalf("got error on fr.json export: %v", err)
		}
		if err = tDE.Export("de.json"); err != nil {
			log.Fatalf("got error on de.json export: %v", err)
		}
		if err = tPT.Export("pt.json"); err != nil {
			log.Fatalf("got error on pt.json export: %v", err)
		}
		if err = tKM.Export("km.json"); err != nil {
			log.Fatalf("got error on km.json export: %v", err)
		}
	} else {
		log.Printf("%v\n\n", tEN)
		log.Printf("%v\n\n", tES)
		log.Printf("%v\n\n", tFR)
		log.Printf("%v\n\n", tDE)
		log.Printf("%v\n\n", tPT)
		log.Printf("%v\n\n", tKM)
	}

}
