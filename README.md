# quran-cli [![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/github.com/vanillaiice/quran-cli) [![Go Report Card](https://goreportcard.com/badge/github.com/vanillaiice/quran-cli)](https://goreportcard.com/report/github.com/vanillaiice/quran-cli)

Read the Holy Quran from your terminal.

# Installation

```sh
$ go install github.com/vanillaiice/quran-cli@latest
```
# Demo

https://github.com/vanillaiice/quran-cli/assets/120596571/b2a9ab2c-cd67-44d7-b2aa-29aa8f4c2202

# Example Usage

```sh
# read the first surah of the Quran in english with arabic text
$ quran-cli read

# read a random surah in english
$ quran-cli read --random

# read surah #55 in french translation only
$ quran-cli read --language fr --number 55 --mode tr

# initialize data for chinese
$ quran-cli init -l zh

# read surah Al-Mulk in chinese
$ quran-cli read -l zh -s mulk
```

> if data for a language is not initialized, it can be initialized
> automatically by answering to the shown prompt.

> By default, the data is stored in the $HOME/.quran-cli directory.

# Help

```sh
NAME:
   quran-cli - Read the Holy Quran from your terminal

USAGE:
   quran-cli [global options] command [command options]

VERSION:
   0.1.0

AUTHOR:
   vanillaiice <vanillaiice1@proton.me>

COMMANDS:
   init, i  initialize data for a language
   read, r  read a surah
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value, -g value  set log level (default: "info")
   --help, -h                   show help
   --version, -v                print the version
```

# Misc

- Please make sure to have the proper arabic [fonts](https://wiki.archlinux.org/title/Localization/Arabic#Fonts) installed on your system (if you plan to read in arabic).
- If the terminal screen is too small, some of the text might not be displayed correctly.

# Todo

- Add more styling to the current display.

# Acknowledgements

- The Quran JSON files are sourced from [risan/quran-json](https://github.com/risan/quran-json), which are in turn taken from [tanzil.net](https://tanzil.net/trans/en.transliteration)
- The different terminal UI styles use [tview](https://github.com/rivo/tview), and [termenv](https://github.com/muesli/termenv).

# License

This project is licensed under the GPLv3 License.

# Author

[vanillaiice](https://github.com/vanillaiice)
