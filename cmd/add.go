package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

// Add 下载代码到vendor中,并添加到gpm.yaml里
type Add struct {
}

func (self *Add) Cmd() cli.Command {
	return cli.Command{
		Name:        "add",
		Usage:       "add repo to gpm.yaml and install package to vendor/",
		Description: "",
	}
}

func (self *Add) Run(ctx *gpm.Ctx) {
	if len(ctx.Args()) != 1 {
		ctx.Die("gpm add need one repo!")
	}

	ctx.MustLoad()

	repo := ctx.Args()[0]
	if ctx.HasDependency(repo) {
		ctx.Die("has repo:%+v", repo)
	}

	dep, err := gpm.NewDependency(repo)
	if err != nil {
		ctx.Die("%+v", err)
	}

	ctx.Get(dep.Remote(), dep.Version)
}
