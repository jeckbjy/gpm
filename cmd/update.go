package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Update struct {
}

func (self *Update) Cmd() cli.Command {
	return cli.Command{
		Name:        "update",
		ShortName:   "up",
		Usage:       "Update a project's dependencies",
		Description: "",
	}
}

// Run update all deps and update lock file
func (self *Update) Run(ctx *gpm.Ctx) {
	ctx.MustLoad()

	for _, dep := range ctx.Imports {
		if err := ctx.Get(dep.Remote(), dep.Version); err != nil {
			ctx.Die("%+v", err)
		}
	}
}
