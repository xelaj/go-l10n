package l10n

import (
	"path/filepath"

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

	return mainContext.MustAll(items)
}

func Message(key string) MessageCode {
	return ImplicitCtxMessage(&mainContext, key)
}

func MessageWithParams(key string) MessageWithParamsCode {
	return ImplicitCtxMessageWithParams(&mainContext, key)
}

func DefaultTranslator() Translator {
	Init()

	return mainContext
}
