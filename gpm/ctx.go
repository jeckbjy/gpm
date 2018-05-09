package gpm

import (
	"os"

	"github.com/codegangsta/cli"
)

// Ctx defines the supporting context
type Ctx struct {
	*cli.Context
	*Logger
	Conf     *Config
	WorkDir  string
	RootDir  string
	CacheDir string
	GOPATH   string
	Args     []string
}

// NewCtx create context
func NewCtx() *Ctx {
	ctx := &Ctx{}
	ctx.Logger = NewLogger()
	ctx.Conf = &Config{}
	ctx.init()
	return ctx
}

func (ctx *Ctx) init() {
	wd, err := os.Getwd()
	if err != nil {
		ctx.Die("cannot get work dir:%s", err)
	}

	ctx.WorkDir = wd
}
