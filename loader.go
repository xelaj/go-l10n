// Package locjson provides the `LoadTranslations()` function to
// load translations from JSON files of the following structure:
//
// {
//     "keyName": {
//         "message": "String to translate",
//         "description": "Description for translators"
//     },
//     ...
// }
package l10n

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
)

type translationFile map[string]struct {
	Message     string
	Description string
}

// Load takes the JSON file name and returns string resources
// in a format compatible with loc.Resources map values.
// It will panic if there's a problem with loading/unmarshaling JSON.
func load(filename string) (Resource, error) {
	t := make(Resource)
	var tf translationFile

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &tf)
	if err != nil {
		return nil, err
	}

	// convert JSON data structure into a destination map format
	for k, v := range tf {
		t[k] = v.Message
	}
	return t, nil
}

func walkDir(path string) (Resource, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	res := make(Resource)

	for _, f := range files {
		if f.IsDir() {
			recursiveres, err := walkDir(filepath.Join(path, f.Name()))
			if err != nil {
				return nil, err
			}
			for k, v := range recursiveres {
				res[k] = v
			}
		}

		r, err := load(filepath.Join(path, f.Name()))
		if err != nil {
			return nil, errors.New(filepath.Join(path, f.Name()) + ": " + err.Error())
		}
		for k, v := range r {
			res[k] = v
		}
	}

	return res, nil
}
