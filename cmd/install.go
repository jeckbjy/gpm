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

func (self *Install) Run(ctx *gpm.Ctx) {

}
