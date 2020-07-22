package l10n

import (
	"github.com/pkg/errors"
	"github.com/xelaj/errs"

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
func (lc *Context) Lang() string {
	return lc.lang
}

// Tr returns a translated version of the string
func (lc *Context) Tr(key string) string {
	res, err := lc.tr(key)
	if err != nil {
		return key
	}
	return res
}

func (lc *Context) TrWithError(key string) (string, error) {
	return lc.tr(key)
}

// Trf returns a formatted version of the string
func (lc *Context) Trf(key string, params plurr.Params) (string, error) {
	translated, err := lc.tr(key)
	if err != nil {
		return "", errors.Wrap(err, "translating")
	}

	s, err := lc.plurr.Format(translated, params)
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

// tr исключительно достает из пула только сами переводы, дополнением занимаются публичные функции
// tr будет стараться искать альтернативные ключи
func (lc *Context) tr(key string) (string, error) {
	err := lc.pool.LoadResource(lc.lang)
	if err != nil {
		return key, errors.Wrap(err, "loading resource")
	}

	resource := lc.pool.Resources[lc.lang]

	s, ok := resource[key]
	if ok {
		return s, nil
	}

	alternativeLocale := ""
	for _, locale := range lc.pool.locales {
		if locale.Code == lc.lang {
			alternativeLocale = locale.Extends
		}
	}
	if alternativeLocale == "" {
		return key, errs.NotFound("key", key)
	}

	ctx, err := lc.pool.GetContext(alternativeLocale)
	if err != nil {
		return key, errors.Wrap(err, "getting context")
	}

	return ctx.tr(key)
}
