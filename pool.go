package l10n

import (
	"path/filepath"

	"github.com/iafan/Plurr/go/plurr"
)

// Resources is an array that holds string resources for a particular language
type Resource map[string]string

// Pool object holds all string resources and language-specific contexts.
// Pool can be created either for an entire project, or for each logical
// part of the application that has its own set of strings.
type Pool struct {
	DefaultLanguage string
	Resources       map[string]Resource
	contexts        map[string]*Context
	resourcePath    string
	appName         string
}

// NewPool creates a new localization pool object
func NewPool(resourcePath, appName, lang string) (*Pool, error) {
	lp := &Pool{
		DefaultLanguage: lang,
		resourcePath:    resourcePath,
		appName:         appName,
		Resources:       make(map[string]Resource),
		contexts:        make(map[string]*Context),
	}
	lp.GetContext(lang)
	return lp, lp.PreloadResource(lang)
}

// GetContext returns a translation context based on the provided language.
// Contexts are thread-safe.
func (lp *Pool) GetContext(lang string) (*Context, error) {
	if lp.contexts[lang] == nil {
		lp.contexts[lang] = &Context{
			pool:  lp,
			lang:  lang,
			plurr: plurr.New().SetLocale(lang),
		}
	}
	err := lp.PreloadResource(lang)
	return lp.contexts[lang], err
}

func (lp *Pool) SetLanguage(lang string) error {
	lp.DefaultLanguage = lang
	lp.GetContext(lang)
	return lp.PreloadResource(lang)
}

func (lp *Pool) Tr(key string) string {
	return lp.contexts[lp.DefaultLanguage].Tr(key)
}

func (lp *Pool) Trf(key string, params plurr.Params) (string, error) {
	return lp.contexts[lp.DefaultLanguage].Trf(key, params)
}

func (lp *Pool) Strf(key string, params plurr.Params) string {
	return lp.contexts[lp.DefaultLanguage].Strf(key, params)
}

func (lp *Pool) PreloadResource(lang string) error {
	res, err := walkDir(filepath.Join(lp.resourcePath, lang, lp.appName))
	if err == nil {
		lp.Resources[lang] = res
	}
	return err
}
