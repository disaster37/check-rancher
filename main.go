package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/disaster37/check-rancher/v2/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var version = "develop"
var commit = ""

func run(args []string) error {

	// Logger setting
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Check Rancher2"
	app.Version = fmt.Sprintf("%s-%s", version, commit)
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "url",
			Usage:   "The rancher URL",
			EnvVars: []string{"RANCHER_URL"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "access-key",
			Usage:   "The rancher access key",
			EnvVars: []string{"RANCHER_ACCESS_KEY"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:    "secret-key",
			Usage:   "The rancher secret key",
			EnvVars: []string{"RANCHER_SECRET_KEY"},
		}),
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Display debug output",
		},
		altsrc.NewInt64Flag(&cli.Int64Flag{
			Name:  "timeout",
			Usage: "The timeout in second",
			Value: 0,
		}),
		&cli.BoolFlag{
			Name:  "no-color",
			Usage: "No print color",
		},
		&cli.BoolFlag{
			Name:  "self-signed-certificate",
			Usage: "Don't check certificate validity. It usefull when use self signed certificate",
		},
		&cli.StringSliceFlag{
			Name:  "ca-path",
			Usage: "List of CA path",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:     "check-snapshot",
			Usage:    "Check the snapshot state",
			Category: "Cluster",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "cluster-name",
					Usage:    "The cluster name",
					Required: true,
				},
				&cli.DurationFlag{
					Name:  "max-older-than",
					Usage: "How many time old the backup",
					Value: 24 * time.Hour,
				},
			},
			Action: cmd.CheckSnapshot,
		},
	}

	app.Before = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if !c.Bool("no-color") {
			formatter := new(prefixed.TextFormatter)
			formatter.FullTimestamp = true
			formatter.ForceFormatting = true
			log.SetFormatter(formatter)
		}

		if c.String("config") != "" {
			before := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))
			return before(c)
		}
		return nil
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(args)
	return err
}

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
