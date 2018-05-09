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

func (self *Name) Run(ctx *gpm.Ctx) {
}

// Name prints the name of the package, according to the config file.
// func Name() {
// 	// cfg := EnsureConfig()
// 	// msg.Puts(cfg.Name)
// }
