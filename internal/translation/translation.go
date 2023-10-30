package translation

import (
	"encoding/json"
	"io"
	"os"
)

type Translation map[string]Values
type Values map[string]interface{}

// writes a Translations object to a file specified by fileName
func (t Translation) Export(fileName string) error {
	// write it
	jsonString, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write(jsonString); err != nil {
		return err
	}

	return nil
}

// writes a Translations object to a file specified by fileName
func ImportFromFile(fileName string) (*Translation, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	r, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	t := new(Translation)
	err = json.Unmarshal(r, t)
	return t, err
}

func (t *Translation) AddValue(cat, label, val string) {
	if _, ok := (*t)[cat]; ok {
		// add it
		(*t)[cat][label] = val
	} else {
		(*t)[cat] = Values{}
		(*t)[cat][label] = val
	}
}
