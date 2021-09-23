package main

import (
	"strings"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	WithCategory("power", powerCmd),
}

func WithCategory(cat string, cmd *cli.Command) *cli.Command {
	cmd.Category = strings.ToUpper(cat)
	return cmd
}

var msgConfig MessagerConfig

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
			&cli.StringFlag{
				Name:  "msg-api",
				Usage: "msg-api",
			},
			&cli.StringFlag{
				Name:  "msg-token",
				Usage: "msg-token",
			},
		},

		Commands: append(Commands),

		Before: func(cctx *cli.Context) error {
			msgConfig = MessagerConfig{
				Url:   cctx.String("msg-api"),
				Token: cctx.String("msg-token"),
			}

			return nil
		},
	}

	app.Setup()

	RunApp(app)
}
