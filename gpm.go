package main

import (
	"os"

	"github.com/jeckbjy/gpm/gpm"

	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/cmd"
)

const version = "0.0.1"
const usage = `Vendor Package Management for your Go projects.

   Each project should have a 'glide.yaml' file in the project directory. Files
   look something like this:

       package: github.com/Masterminds/glide
       imports:
       - package: github.com/Masterminds/cookoo
         version: 1.1.0
       - package: github.com/kylelemons/go-gypsy
         subpackages:
         - yaml

   For more details on the 'glide.yaml' files see the documentation at
   https://glide.sh/docs/glide.yaml
`

func main() {
	app := cli.NewApp()
	app.Name = "gpm"
	app.Usage = usage
	app.Version = version

	// setup commands
	cmds := cmd.New()
	for _, c := range cmds {
		action := c
		cliCmd := c.Cmd()
		cliCmd.Action = func(cliCtx *cli.Context) error {
			ctx := gpm.NewCtx()
			ctx.Context = cliCtx

			ctx.Debug("run cmd:%+s", cliCtx.Command.Name)
			action.Run(ctx)
			return nil
		}

		app.Commands = append(app.Commands, cliCmd)
	}

	// run
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
