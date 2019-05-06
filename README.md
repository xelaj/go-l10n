# go-l10n - Localization for humans

**english** [русский](https://github.com/xelaj/go-l10n/blob/master/doc/ru_RU/README.md)

Quick and easy localization of your Go projects!

This package is based on the [go-l10n package](https://github.com/iafan/go-l10n) created by
[iafan](https://github.com/iafan).

## Getting Started

as simple as another packages:

`go get github.com/xelaj/go-l10n`

## how to use

**code examples are [here](https://github.com/xelaj/go-l10n/examples/main.go)**

The main struct you will use is `l10n.Pool`. on init func it takes filepath to
locale dirs, the name of your application and the main language into which you will translate
messages.

messages are taken from the directory that you specify during initialization. however, it must have
certain structure:

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

`info.cfg` is a summary of the locale. it is written to the l10n.Locale object. she keeps in
The following parameters:

```ini
# name of language on ascii
name = Russian

# localized name of language
translate = Russian

# which language it extends. if you have unsupported locale name, package returns key name
extends =
```

all three parameters are optional. however, the file itself in each folder with the same name is required.

`YourApp` essentially implements the namespaces of various applications.

inside the directory of your application, you can do absolutely any structure of your files,
the only important thing is that the names of the messages in all these files do not overlap. if that happens then
The translation closest to the root directory will be selected.

each locale file is a json file with the following parameters:

``` json
{
    "about": {
        "message": "this is locale file!",
        "description": "use this description to message message"
    },
    "hello text": {
        "message": "hello!"
    }
}
```

examples of file struvture can be seen [here](https://github.com/xelaj/go-l10n/examples/locale)

## Contributing

Please read [contributing guide](https://github.com/xelaj/go-l10n/doc/ru_RU/CONTRIBUTING.md) if you want to help. And the help is very necessary!

## Authors

* **Igor Afanasyev** — *base library* - [iafan](https://github.com/iafan)
* **Richard Cooper** — *other stuff* - [ololosha228](https://github.com/ololosha228)


## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/xelaj/go-l10n/doc/en_US/LICENSE) file for details
