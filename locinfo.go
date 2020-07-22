package l10n

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/zieckey/goini"
)

type Locale struct {
	Code      string
	Name      string
	Translate string
	Extends   string
}

func langInfo(path, name string) (*Locale, error) {
	cfg := goini.New()
	err := cfg.ParseFile(filepath.Join(path, name, "info.cfg"))
	if err != nil {
		return nil, errors.Wrap(err, "parsing file")
	}

	l := &Locale{}
	l.Code = name
	l.Name, _ = cfg.Get("name")
	l.Translate, _ = cfg.Get("translate")
	l.Extends, _ = cfg.Get("extends")

	return l, nil
}
