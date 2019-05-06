package main

import (
	"fmt"

	"github.com/iafan/Plurr/go/plurr"
	"github.com/xelaj/go-l10n"
)

func main() {
	// create localization pool with "en" as a default language
	// and load string resources
	lp, err := l10n.NewPool("/path/to/your/system/locale", "YourAppName", "en_GB")
	if err != nil {
		fmt.Println("you don't have en_GB localization!")
		fmt.Println("'lp' will work, but return message keys only!")
	}

	// locale packs are preload when you call Tr function (if pool have no data about lang, it
	// use PreloadResource function)
	// if you want preload it obviously, use lp.PreloadResource("yourlang")
	err = lp.PreloadResource("ru_RU")
	if err != nil {
		panic(err)
	}

	// you can simple translate messages!
	fmt.Println(lp.Tr("Hello"))

	// and format it. Trf method will return error, if plurr.Params are incorrect.
	fmt.Println(lp.Trf("YouHaveNMessages", plurr.Params{"N": 1}))
	fmt.Println(lp.Trf("YouHaveNMessages", plurr.Params{"N": 4}))

	// if you don't want get errors, use Strf. it will return key, if error is happened.
	fmt.Println(lp.Strf("YouHaveNMessages", plurr.Params{"N": 2}))
	fmt.Println(lp.Strf("YouHaveNMessages", plurr.Params{"some": "bad params"}))

	// you can get locale context from pool
	an, err := lp.GetContext("andorian")
	ru, err := lp.GetContext("ru_RU")
	if err != nil {
		fmt.Println("you don't have these locales! but context will partially work if you need")
	}

	// and use it in one time
	fmt.Println(lp.Tr("Hello"))
	fmt.Println(an.Tr("Hello"))
	fmt.Println(ru.Tr("Hello"))

	// you can change lang of pool
	lp.SetLanguage("ru_RU")

	// you can even format messages with context.
	fmt.Println(ru.Strf("YouHaveNMessages", plurr.Params{"N": 4}))
}
