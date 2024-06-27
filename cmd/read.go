package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/quran-cli/db"
	"github.com/vanillaiice/quran-cli/tui"
	"github.com/vanillaiice/quran-cli/tui/list"
	"github.com/vanillaiice/quran-cli/tui/tview"
)

// maxSurahId is the maximum surah id in the Quran.
const maxSurahId = 114

// readCmd is the read command.
// It prints a surah by providing its name or number.
var readCmd = &cli.Command{
	Name:    "read",
	Aliases: []string{"r"},
	Usage:   "read a surah",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "data-path",
			Aliases: []string{"p"},
			Usage:   "data path `PATH`",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "style",
			Aliases: []string{"t"},
			Usage:   "use terminal ui style `STYLE` (tview, list)",
			Value:   "list",
		},
		&cli.StringFlag{
			Name:    "language",
			Aliases: []string{"l"},
			Usage:   "read in `LANGUAGE`",
			Value:   "en",
		},
		&cli.StringFlag{
			Name:    "mode",
			Aliases: []string{"m"},
			Usage:   "reading mode `MODE` (arabic, translation, both)",
			Value:   "both",
		},
		&cli.StringFlag{
			Name:    "surah",
			Aliases: []string{"s"},
			Usage:   "read surah `SURAH`",
		},
		&cli.BoolFlag{
			Name:    "exact",
			Aliases: []string{"e"},
			Usage:   "read surah `SURAH` by exact name (case insensitive)",
		},
		&cli.IntFlag{
			Name:    "number",
			Aliases: []string{"n"},
			Usage:   "read surah number `NUMBER`",
			Value:   1,
		},
		&cli.BoolFlag{
			Name:    "random",
			Aliases: []string{"r"},
			Usage:   "read a random surah",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		lang := langCode(ctx.String("language"))
		switch lang {
		case Arabic, Bengali, Chinese, English, Spanish, French, Indonesian, Russian, Swedish, Turkish, Urdu, Transliteration:
		default:
			return fmt.Errorf("unsupported language: %q", lang)
		}

		var mode tui.Lang
		switch ctx.String("mode") {
		case "arabic", "ar":
			mode = tui.Arabic
		case "translation", "tr":
			mode = tui.Translation
		case "both", "bo":
			mode = tui.Both
		default:
			return fmt.Errorf("unsupported mode: %q", ctx.String("mode"))
		}

		dataPath := ctx.String("data-path")

		if dataPath == "" {
			dataPath, err = os.UserHomeDir()
			if err != nil {
				return
			}
			dataPath = path.Join(dataPath, dataDir)
		}

		dbPath := path.Join(dataPath, fmt.Sprintf("quran_%s.db", lang))

		if _, err = os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("database %q not found, create it ? (y/N)\n -> ", dbPath)

			var ans string
			_, err = fmt.Scan(&ans)
			if err != nil {
				return
			}
			ans = strings.ToLower(ans)

			if ans == "y" || ans == "yes" {
				err = initFunc(lang, dataPath, true)
				if err != nil {
					return
				}
			} else {
				log.Warn("not creating database")
				return
			}
		}

		d, err := db.New(dbPath)
		if err != nil {
			return
		}
		defer d.Close()

		var surah *db.Surah

		if ctx.Bool("random") {
			surah, err = d.GetSurahById(rand.Intn(maxSurahId) + 1)
			if err != nil {
				return
			}
		} else {
			if ctx.String("surah") != "" {
				if !ctx.Bool("exact") {
					surah, err = d.GetSurahByNameLike(ctx.String("surah"))
					if err != nil {
						return
					}
				} else {
					surah, err = d.GetSurahByName(ctx.String("surah"))
					if err != nil {
						return
					}
				}

				if len(surah.Verses) == 0 {
					return fmt.Errorf("surah %q not found", ctx.String("surah"))
				}
			} else if ctx.Int("number") != 0 {
				surah, err = d.GetSurahById(ctx.Int("number"))
				if err != nil {
					return
				}

				if len(surah.Verses) == 0 {
					return fmt.Errorf("surah #%d not found", ctx.Int("number"))
				}
			} else {
				return fmt.Errorf("please specify surah name or number")
			}
		}

		switch ctx.String("style") {
		case "tview", "tv":
			err = tview.Run(surah, mode)
		case "list", "li":
			err = list.Run(surah, mode)
		default:
			err = fmt.Errorf("invalid style: %q", ctx.String("style"))
		}

		return
	},
}
