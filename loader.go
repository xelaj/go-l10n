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
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ungerik/go-dry"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Load takes the JSON file name and returns string resources
// in a format compatible with loc.Resources map values.
// It will panic if there's a problem with loading/unmarshaling JSON.
func load(filename string) (Resource, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	switch ext := filepath.Ext(filename); ext {
	case ".yaml", ".yml":
		return loadYamlFile(data)
	case ".json":
		return loadJsonFile(data)
	default:
		return nil, errors.New("unknown format: " + ext)
	}
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

func loadYamlFile(data []byte) (Resource, error) {
	res := make(map[string]interface{})

	err := yaml.Unmarshal(data, &res)
	if err != nil {
		return nil, errors.Wrap(err, "parsing file")
	}

	tr := make(Resource)

	for name, data := range res {
		switch value := data.(type) {
		case map[interface{}]interface{}:
			found := false
			for ki, vi := range value {
				k, ok := ki.(string)
				if !ok {
					return nil, errors.New("keys only strings")
				}

				if dry.StringInSlice(k, []string{"msg", "message"}) {
					tr[name] = fmt.Sprint(vi)
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("resource " + name + "doesn't have message property")
			}
		default:
			tr[name] = fmt.Sprint(value)
		}
	}

	return tr, nil
}

func loadJsonFile(data []byte) (Resource, error) {
	res := make(map[string]interface{})

	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, errors.Wrap(err, "parsing file")
	}

	tr := make(Resource)

	for name, data := range res {
		switch value := data.(type) {
		case map[string]interface{}:
			found := false
			for key, vi := range value {
				if dry.StringInSlice(key, []string{"msg", "message"}) {
					tr[name] = fmt.Sprint(vi)
					found = true
					break
				}
			}
			if !found {
				return nil, errors.New("resource " + name + "doesn't have message property")
			}
		default:
			tr[name] = fmt.Sprint(value)
		}
	}

	return tr, nil
}
