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

func (lc *Context) Message(key string) MessageCode {
	t := Translator(lc)
	return ImplicitCtxMessage(&t, key)
}
func (lc *Context) MessageWithParams(key string) MessageWithParamsCode {
	t := Translator(lc)
	return ImplicitCtxMessageWithParams(&t, key)
}

func (lc *Context) MustAll(items []string) error {
	errorsMultiple := &errs.MultipleErrors{}
	for _, item := range items {
		_, err := lc.TrWithError(item)
		errorsMultiple.Add(err)
	}
	return errorsMultiple.Normalize()
}

type MsgCodeReturner interface {
	Code() string
}

type codeGetter uint8

type MessageCode func(privateParamss ...interface{}) string

func ImplicitCtxMessage(lc *Translator, key string) MessageCode {
	return func(privateParams ...interface{}) string {
		if len(privateParams) > 0 {
			if _, ok := privateParams[0].(codeGetter); !ok {
				panic("do not use additional parameters")
			}
			return key
		}
		return (*lc).Tr(key)
	}
}

func (m MessageCode) Code() string {
	return m(codeGetter(0))
}

type MessageWithParamsCode func(params map[string]interface{}, privateParams ...interface{}) string

func ImplicitCtxMessageWithParams(lc *Translator, key string) MessageWithParamsCode {
	return func(params map[string]interface{}, privateParams ...interface{}) string {
		if len(privateParams) > 0 {
			if _, ok := privateParams[0].(codeGetter); !ok {
				panic("do not use additional parameters")
			}
			return key
		}
		return (*lc).Strf(key, params)
	}
}

func (m MessageWithParamsCode) Code() string {
	return m(nil, codeGetter(0))
}
