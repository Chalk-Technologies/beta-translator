package repoUpload

import (
	"context"
	"encoding/json"
	"github.com/Chalk-Technologies/beta-translator/internal/translation"
	"github.com/google/go-github/v56/github"
	"log"
	"os"
	"strings"
)

// add export files to repo

var client *github.Client

func Init() {
	githubToken := os.Getenv("GITHUB_TOKEN")
	log.Printf("authenticating to github with token %v\n", githubToken)
	client = github.NewClient(nil).WithAuthToken(githubToken)
	return
}

var uploadMsg = "Auto-update translations - beta-translator"

func UploadFile(fileName string, repo string, path string, content translation.Translation) error {
	owner := strings.Split(repo, "/")[0]
	r := strings.Split(repo, "/")[1]
	p := path + fileName
	// get the sha of the existing file
	f, _, _, err := client.Repositories.GetContents(context.Background(), owner, r, p, nil)
	if err != nil {
		return err
	}

	jsonString, err := json.Marshal(content)
	if err != nil {
		return err
	}
	//var jsonEncodedString []byte
	//jsonEncodedString := base64.StdEncoding.EncodeToString(jsonString)

	opts := &github.RepositoryContentFileOptions{
		SHA:     f.SHA,
		Message: &uploadMsg,
		Content: jsonString,
	}
	if _, _, err = client.Repositories.UpdateFile(context.Background(), owner, r, p, opts); err != nil {
		return err
	}
	return nil
}
