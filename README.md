# BETA TRANSLATOR


## Export Translations
`go run cmd/main.go -notionSecret <YOUR_NOTION_SECRET> -export -db <YOUR_NOTION_DB_ID>`

go run cmd/main.go -notionSecret secret_rzwHEGnB6sQ3gBTAqhLVcryjxMBtuWASw4yszZjn5On -export -db 4bf9b6b8e2b74d528dd5419f5406caa7
go run cmd/main.go -sheetId 1cNwPjNEHhUdYrCx31G88trjBv7ZZDsoEXBMnfevSNJo -export

## Upload Translations to GCloud
go run cmd/main.go -sheetId 1cNwPjNEHhUdYrCx31G88trjBv7ZZDsoEXBMnfevSNJo -upload


## Import Translations
`go run cmd/main.go -notionSecret <YOUR_NOTION_SECRET> -importFile es_current.json -importField "Text Spanish" -db <YOUR_NOTION_DB_ID>`


