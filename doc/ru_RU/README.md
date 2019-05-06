# go-l10n — Localization for humans

[english](https://github.com/xelaj/go-l10n/README.md) **русский**

быстрая и простая локализация ваших проектов на Go!

этот пакет базируется на пакете [go-l10n](https://github.com/iafan/go-l10n) от
[iafan](https://github.com/iafan).

## Getting Started

as simple as another packages:

`go get github.com/xelaj/go-l10n`

## how to use

**примеры кода находятся [здесь](https://github.com/xelaj/go-l10n/examples/main.go)**

главная структура, которую вы будете использовать, это `l10n.Pool`. На вход она принимает путь до
директории с локализацией, имя вашего приложения и основной язык, на который вы будете переводить
сообщения.

сообщения берутся из директории, которую вы укажете при инициализации. однако, она должна иметь
определенную структуру:

```
locale
    en_US
        YourApp
            locale.json
            another folder
                locale.json
                third.json
        info.cfg
    ru_RU
        YourApp
            locale.json
            another folder
                locale.json
                third.json
        info.cfg
    ua_UA
        ...
    es_ES
        ...
```

`info.cfg` это краткая информация о локали. она записывается в объект l10n.Locale. она хранит в
себе следующие параметры:

``` ini
# name of language on english with ascii symbols only
name=Russian

# localized name of language
translate=Русский

# which language it extends. if you have unsupported locale name, package returns key name
extends=
```

все три параметра являются опциональными. однако, сам файл в каждой папке с данным названием обязателен.

`YourApp` по сути реализует неймспейсы различных приложений.

внутри директории вашего приложения вы можете делать абсолютно любую структуру ваших файлов,
важно лишь, что бы названия сообщений во всех этих файлах не пересекались. если это произойдет, то
будет выбран перевод ближайший к корневой директории

каждый файл локали является json файлом следующими параметрами:

``` json
{
    "about":{
        "message":"this is locale file!",
        "description":"use this description to describe message"
    },
    "hello text":{
        "message":"hello!"
    }
}
```

примеры структуры файлов можете посмотреть [здесь](https://github.com/xelaj/go-l10n/examples/locale)

## Contributing

пожалуйста, прочитайте [информацию о помощи](https://github.com/xelaj/go-l10n/doc/ru_RU/CONTRIBUTING.md), если хотите помочь. А помощь очень нужна!

## Authors

* **Igor Afanasyev** - *базовая библиотека* - [iafan](https://github.com/iafan)
* **Richard Cooper** — *other stuff* - [ololosha228](https://github.com/ololosha228)

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/xelaj/go-l10n/doc/ru_RU/LICENSE) file for details
