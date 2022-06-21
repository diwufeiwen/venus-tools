
package main

import (
	"strings"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	WithCategory("power", powerCmd),
	WithCategory("fil", filCmd),
	WithCategory("market", marketCmds),
}

func WithCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = strings.ToUpper(cat)
	return cmd
}

var msgConfig VenusConfig

func main() {
	app := &cli.App{
		Name:                 "venus-tool",
		Usage:                "Some tools related to venus",
		Version:              UserVersion(),
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "color",
			},
/*			&cli.StringFlag{
				Name:  "node-api",
				Usage: "node-api",
			},*/
			&cli.StringFlag{
				Name:  "msg-api",
				Usage: "msg-api",
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "token",
			},
		},

		Commands: append(Commands),

		Before: func(cctx *cli.Context) error {
			msgConfig = VenusConfig{
				NodeUrl:    cctx.String("node-api"),
				MessageUrl: cctx.String("msg-api"),
				Token:      cctx.String("token"),
			}

			return nil
		},
	}

	app.Setup()

	RunApp(app)
}
