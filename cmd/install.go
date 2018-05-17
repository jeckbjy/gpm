package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Install struct {
}

func (self *Install) Cmd() cli.Command {
	return cli.Command{
		Name:        "install",
		ShortName:   "i",
		Usage:       "Install a project's dependencies",
		Description: "",
	}
}

// Run install all deps
func (self *Install) Run(ctx *gpm.Ctx) {
	if len(ctx.Args()) != 0 {
		ctx.Die("install donot need args")
	}

	ctx.Load()

	// get all
	for _, dep := range ctx.Conf.Imports {
		if err := ctx.Get(dep); err != nil {
			ctx.Die("%+v", err)
		}
	}
}
