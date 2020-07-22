package l10n

import (
	"path/filepath"

	"github.com/xelaj/errs"

	"github.com/gobuffalo/envy"
	"github.com/iafan/Plurr/go/plurr"
	"github.com/pkg/errors"
	"github.com/ungerik/go-dry"
	"github.com/xelaj/v"
)

var (
	AppName     = ""
	LocalesPath = ""
	builtinPool *Pool
	mainContext Translator
)

func Init() {
	if builtinPool != nil {
		return
	}
	if AppName == "" {
		AppName = v.AppName
	}
	if LocalesPath == "" {
		LocalesPath = "/usr/share/locale"
	}

	var err error
	builtinPool, err = NewPool(LocalesPath, AppName, false)
	dry.PanicIfErr(err)
	mainContext, err = GetContext(envy.Get("LANGUAGE", "en_US"))
	dry.PanicIfErr(err)
}

func SetLanguage(lang string) error {
	Init()
	ctx, err := builtinPool.GetContext(lang)
	if err != nil {
		return errors.Wrap(err, "getting language context")
	}
	mainContext = ctx
	return nil
}

func Tr(key string) string {
	Init()
	return mainContext.Tr(key)
}

func Trf(key string, params plurr.Params) (string, error) {
	Init()
	return mainContext.Trf(key, params)
}

func Strf(key string, params plurr.Params) string {
	Init()
	return mainContext.Strf(key, params)
}

func GetContext(lang string) (Translator, error) {
	Init()

	err := LoadResource(lang)
	if err != nil {
		return nil, errors.Wrap(err, "loading resource")
	}

	return &Context{
		pool:  builtinPool,
		lang:  lang,
		plurr: plurr.New().SetLocale(lang),
	}, nil
}

func LoadResource(lang string) error {
	Init()

	res, err := walkDir(filepath.Join(builtinPool.resourcePath, lang, builtinPool.appName))
	if err == nil {
		builtinPool.Resources[lang] = res
	}
	return errors.Wrap(err, "can't load '"+lang+"'")
}

func MustAll(items []string) error {
	Init()

	errorsMultiple := &errs.MultipleErrors{}
	for _, item := range items {
		_, err := mainContext.TrWithError(item)
		errorsMultiple.Add(err)
	}
	return errorsMultiple.Normalize()
}

type MsgCodeReturner interface {
	Code() string
}

type codeGetter uint8

type MessageCode func(privateParamss ...interface{}) string

func Message(key string) MessageCode {
	return func(privateParams ...interface{}) string {
		if len(privateParams) > 0 {
			if _, ok := privateParams[0].(codeGetter); !ok {
				panic("do not use additional parameters")
			}
			return key
		}
		return Tr(key)
	}
}

func (m MessageCode) Code() string {
	return m(codeGetter(0))
}

type MessageWithParamsCode func(params map[string]interface{}, privateParams ...interface{}) string

func MessageWithParams(key string) MessageWithParamsCode {
	return func(params map[string]interface{}, privateParams ...interface{}) string {
		if len(privateParams) > 0 {
			if _, ok := privateParams[0].(codeGetter); !ok {
				panic("do not use additional parameters")
			}
			return key
		}
		return Strf(key, params)
	}
}

func (m MessageWithParamsCode) Code() string {
	return m(nil, codeGetter(0))
}
