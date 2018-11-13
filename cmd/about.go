package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

const aboutMessage = `Go vendor管理工具`

type About struct {
}

func (self *About) Cmd() cli.Command {
	return cli.Command{
		Name:        "about",
		Usage:       "Learn about gpm",
		Description: "",
	}
}

func (self *About) Run(ctx *gpm.Ctx) {
	ctx.Debug("about")
	ctx.Puts(aboutMessage)
}
