package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Name struct {
}

func (self *Name) Cmd() cli.Command {
	return cli.Command{
		Name:        "name",
		Usage:       "Print the name of this project.",
		Description: "Read the glide.yaml file and print the name given on the 'package' line.",
	}
}

// Run return name from config
func (self *Name) Run(ctx *gpm.Ctx) {
	ctx.Load()
	ctx.Puts(ctx.Conf.Name)
}
