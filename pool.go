package l10n

import (
	"path/filepath"

	"github.com/xelaj/errs"

	"github.com/ungerik/go-dry"

	"github.com/iafan/Plurr/go/plurr"
	"github.com/pkg/errors"
)

type Translator interface {
	Tr(key string) string
	TrWithError(key string) (string, error)
	Trf(key string, params plurr.Params) (string, error)
	Strf(key string, params plurr.Params) string
	Lang() string
}

// Resources is an array that holds string resources for a particular language
type Resource map[string]string

// Pool object holds all string resources and language-specific contexts.
// Pool can be created either for an entire project, or for each logical
// part of the application that has its own set of strings.
type Pool struct {
	Resources    map[string]Resource
	resourcePath string
	appName      string
	locales      []*Locale
}

// NewPool creates a new localization pool object
func NewPool(resourcePath, appName string, preloadAll bool) (*Pool, error) {
	lp := &Pool{
		resourcePath: resourcePath,
		appName:      appName,
		Resources:    make(map[string]Resource),
	}
	locales, err := lp.availableLocales()
	if err != nil {
		return nil, errors.Wrap(err, "getting all locales")
	}
	lp.locales = locales

	if preloadAll {
		err := lp.PreloadAll()
		if err != nil {
			return nil, errors.Wrap(err, "preloading all locales")
		}
	}

	return lp, nil
}

// GetContext returns a translation context based on the provided language.
// Contexts are thread-safe.
func (lp *Pool) GetContext(lang string) (*Context, error) {
	err := lp.LoadResource(lang)
	if err != nil {
		return nil, errors.Wrap(err, "loading resource")
	}

	return &Context{
		pool:  builtinPool,
		lang:  lang,
		plurr: plurr.New().SetLocale(lang),
	}, nil
}

func (lp *Pool) ForceLoadResource(lang string) error {
	res, err := walkDir(filepath.Join(lp.resourcePath, lang, lp.appName))
	if err == nil {
		lp.Resources[lang] = res
	}
	return err
}

func (lp *Pool) LoadResource(lang string) error {
	if _, ok := lp.Resources[lang]; ok {
		return nil
	}

	return lp.ForceLoadResource(lang)
}

func (lp *Pool) PreloadAll() error {
	resources := make(map[string]Resource)

	errMultiple := &errs.MultipleErrors{}
	for _, locale := range lp.locales {
		path := filepath.Join(lp.resourcePath, locale.Code, lp.appName)
		if !dry.FileExists(path) {
			continue
		}

		res, err := walkDir(path)
		if err != nil {
			errMultiple.Add(err)
			continue
		}

		resources[locale.Code] = res
	}

	if err := errMultiple.Normalize(); err != nil {
		return errors.Wrap(err, "loading all languages")
	}

	for code, resource := range resources {
		lp.Resources[code] = resource
	}

	return nil
}

func (lp *Pool) availableLocales() ([]*Locale, error) {
	dirs, err := dry.ListDirDirectories(lp.resourcePath)
	dry.PanicIfErr(err)

	locales := make([]*Locale, 0)
	errMultiple := &errs.MultipleErrors{}
	for _, dir := range dirs {
		locale, err := langInfo(lp.resourcePath, dir)
		if err != nil {
			errMultiple.Add(err)
			continue
		}
		locales = append(locales, locale)
	}
	if err := errMultiple.Normalize(); err != nil {
		return nil, err
	}
	return locales, nil
}
