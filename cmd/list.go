package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type List struct {
}

func (self *List) Cmd() cli.Command {
	return cli.Command{
		Name:        "list",
		Usage:       "List prints all dependencies that the present code references.",
		Description: "",
	}
}

func (self *List) Run(ctx *gpm.Ctx) {

}
