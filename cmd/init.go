package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/quran-cli/db"
)

// initCmd is the init command.
var initCmd = &cli.Command{
	Name:    "init",
	Aliases: []string{"i"},
	Usage:   "initialize data for a language",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "data-path",
			Aliases: []string{"p"},
			Usage:   "save data in `PATH`",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "language",
			Aliases: []string{"l"},
			Usage:   "init for language `LANGUAGE`",
			Value:   "en",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "delete existing data if exists",
			Value:   false,
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		return initFunc(langCode(ctx.String("language")), ctx.String("data-path"), ctx.Bool("force"))
	},
}

// initFunc downloads the needed data and initializes
// the database for a specific language.
var initFunc = func(lang langCode, dataPath string, force bool) (err error) {
	switch lang {
	case Arabic, Bengali, Chinese, English, Spanish, French, Indonesian, Russian, Swedish, Turkish, Urdu, Transliteration:
	default:
		return fmt.Errorf("unsupported language: %q", lang)
	}

	if dataPath == "" {
		dataPath, err = os.UserHomeDir()
		if err != nil {
			return
		}
		dataPath = path.Join(dataPath, dataDir)
	}

	dbPath := path.Join(dataPath, fmt.Sprintf("quran_%s.db", lang))

	if _, err = os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dataPath, os.ModePerm)
		if err != nil {
			return
		}
		log.Debugf("created data directory %q", dataPath)
	} else if err != nil {
		return
	} else {
		if _, err = os.Stat(dbPath); err == nil {
			if force {
				if err = os.Remove(dbPath); err != nil {
					return
				}
				log.Warnf("deleted existing database for language %s", lang)
			} else {
				return fmt.Errorf("database already exists for language %s", lang)
			}
		}
	}

	log.Debugf("downloading file quran_%s.json...", lang)

	resp, err := http.Get(sources[lang])
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Debugf("downloaded file quran_%s.json", lang)

	d, err := db.New(dbPath)
	if err != nil {
		return
	}
	defer d.Close()

	log.Debugf("intializing quran database for language %s...", lang)

	if err = d.InitFromReader(resp.Body); err != nil {
		return
	}

	log.Infof("initialized quran database for language %s", lang)

	return
}
