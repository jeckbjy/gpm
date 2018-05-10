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

func (self *Update) Run(ctx *gpm.Ctx) {
	ctx.Debug("Update")
}
