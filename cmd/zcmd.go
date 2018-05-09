package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

// Command 接口
type Command interface {
	Cmd() cli.Command
	Run(ctx *gpm.Ctx)
}

// New 创建所有命令
func New() []Command {
	cmds := []Command{
		&About{},
		&Build{},
		&Create{},
		&Get{},
		&Info{},
		&Install{},
		&List{},
		&Name{},
		&Remove{},
		&Update{},
	}

	return cmds
}
