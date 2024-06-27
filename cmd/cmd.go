package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
	"github.com/vanillaiice/quran-cli/version"
)

const (
	perm    = 0644         // file mode
	dataDir = ".quran-cli" // data directory
)

// Exec executes the app.
func Exec() {
	app := cli.App{
		Name:    "quran-cli",
		Usage:   "Read the Holy Quran from your terminal",
		Authors: []*cli.Author{{Name: "vanillaiice", Email: "vanillaiice1@proton.me"}},
		Version: version.Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Aliases: []string{"g"},
				Usage:   "set log level",
				Value:   "info",
			},
		},
		Commands: []*cli.Command{
			initCmd,
			readCmd,
		},
	}

	app.Before = func(ctx *cli.Context) error {
		var logLevel log.Level

		switch ctx.String("log-level") {
		case "debug":
			logLevel = log.DebugLevel
		case "warn":
			logLevel = log.WarnLevel
		case "error":
			logLevel = log.ErrorLevel
		case "fatal":
			logLevel = log.FatalLevel
		case "info":
			logLevel = log.InfoLevel
		default:
			return fmt.Errorf("invalid log level: %q", ctx.String("level"))
		}

		log.SetLevel(logLevel)

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
