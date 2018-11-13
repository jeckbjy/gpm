package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

// 仅下载代码放到vendor中,但并不放入gpm.yaml中
type Get struct {
}

func (self *Get) Cmd() cli.Command {
	return cli.Command{
		Name:        "get",
		Usage:       "like go get,but install package to vendor/",
		Description: ``,
	}
}

// Run 添加并下载某个repo
func (self *Get) Run(ctx *gpm.Ctx) {
	if len(ctx.Args()) != 1 {
		ctx.Die("get need just one repo!")
	}

	ctx.MustLoad()

	url := ctx.Args()[0]

	repo, version := gpm.ParseRepo(url)

	ctx.Info("get repo:%+v", url)
	if err := ctx.Get(repo, version); err != nil {
		ctx.Die("get repo fail:%+v", err)
	}
	ctx.Info("save repo to vendor, but not insert to gpm.yaml")
}
