package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Remove struct {
}

func (self *Remove) Cmd() cli.Command {
	return cli.Command{
		Name:      "remove",
		ShortName: "rm",
		Usage:     "",
	}
}

func (self *Remove) Run(ctx *gpm.Ctx) {
}

// Remove removes a dependncy from the configuration.
// func Remove(packages []string) {
// 	// EnsureConfig()
// }
