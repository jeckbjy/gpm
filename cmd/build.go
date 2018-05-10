package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Build struct {
}

func (self *Build) Cmd() cli.Command {
	return cli.Command{
		Name:        "build",
		Usage:       "",
		Description: "",
	}
}

func (self *Build) Run(ctx *gpm.Ctx) {
	ctx.Debug("Build")
}
