package cloudStorage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Chalk-Technologies/beta-translator/internal/translation"
)

var cloudStorageClient *storage.Client

const betaAssetsBucketName = "beta-assets"

// UploadTranslations adds a new version of each language file to the beta-assets cloud storage bucket at translations/en.json for example
func UploadTranslations(lang string, t translation.Translation) error {
	ctx := context.Background()
	bucket := cloudStorageClient.Bucket(betaAssetsBucketName)
	wc := bucket.Object(fmt.Sprintf("translations/%v.json", lang)).NewWriter(ctx)
	wc.ContentType = "application/json"
	wc.CacheControl = "public"

	defer wc.Close()
	jerr := json.NewEncoder(wc).Encode(t)
	return jerr
}

func Init() error {
	ctx := context.Background()
	client, gerr := storage.NewClient(ctx)
	if gerr != nil {
		return gerr
	}
	cloudStorageClient = client
	return nil
}
