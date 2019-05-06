package l10n

import (
	"errors"

	"github.com/iafan/Plurr/go/plurr"
)

// Context holds a localization context
// (a combination of a language, a corresponding Plurr formatter,
// and a link back to parent localization Pool object)
type Context struct {
	pool  *Pool
	lang  string
	plurr *plurr.Plurr
}

// GetLanguage returns the current language of the context.
func (lc *Context) GetLang() string {
	return lc.lang
}

// Tr returns a translated version of the string
func (lc *Context) Tr(key string) string {
	return lc.tr(key, lc.lang)
}

// Trf returns a formatted version of the string
func (lc *Context) Trf(key string, params plurr.Params) (string, error) {
	s, err := lc.plurr.Format(lc.Tr(key), params)
	if err != nil {
		return "", errors.New(key + ": " + err.Error())
	}

	return s, nil
}

func (lc *Context) Strf(key string, params plurr.Params) string {
	r, err := lc.Trf(key, params)
	if err != nil {
		return lc.Tr(key)
	}
	return r
}

func (lc *Context) tr(key, lang string) string {
	resource, ok := lc.pool.Resources[lang]
	if !ok {
		err := lc.pool.PreloadResource(lang)
		if err != nil {
			return key
		}

		resource = lc.pool.Resources[lang]
	}

	s, ok := resource[key]
	if ok {
		return s
	}
	return lc.trAlternate(key, lang)
}

func (lc *Context) trAlternate(key, lang string) string {
	info, err := langInfo(lc.pool.resourcePath, lang)
	if err != nil {
		return key
	}

	if info.Extends != "" {
		return lc.tr(key, info.Extends)
	}

	return key
}
