package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

// Info prints information about a project based on a passed in format.
type Info struct {
}

func (self *Info) Cmd() cli.Command {
	return cli.Command{
		Name:  "info",
		Usage: "Info prints information about this project",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "format, f",
				Usage: `Format of the information wanted (required).`,
			},
		},
		Description: `A format containing the text with replacement variables
		has to be passed in. Those variables are:
	 
			%n - name
			%d - description
			%h - homepage
			%l - license
	 
		For example, given a project with the following glide.yaml:
	 
			package: foo
			homepage: https://example.com
			license: MIT
			description: Some example description
	 
		Then running the following commands:
	 
			glide info -f %n
			   prints 'foo'
	 
			glide info -f "License: %l"
			   prints 'License: MIT'
	 
			glide info -f "%n - %d - %h - %l"
			   prints 'foo - Some example description - https://example.com - MIT'`,
	}
}

func (self *Info) Run(ctx *gpm.Ctx) {
	if !ctx.IsSet("format") {
		cli.ShowCommandHelp(ctx.Context, ctx.Command.Name)
		return
	}

	//
}
