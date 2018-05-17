package cmd

import (
	"io"
	"os"
	"path/filepath"

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

// Run 删除某个repo
func (self *Remove) Run(ctx *gpm.Ctx) {
	if len(ctx.Args()) != 1 {
		ctx.Die("remove need just one repo!")
	}

	ctx.Load()

	name := ctx.Args()[0]
	if !ctx.Conf.DelDependency(name) {
		ctx.Die("cannot find dependency:%+v", name)
	}

	ctx.Info("remove deps: %+v", name)

	// remove from cache
	removeAll(ctx.CacheDir, name)
	// cache := filepath.Join(ctx.CacheDir, name)
	// os.RemoveAll(cache)

	// remove from vendor
	removeAll(ctx.VendorDir, name)
	// vendor := filepath.Join(ctx.VendorDir, name)
	// os.RemoveAll(vendor)

	// TODO: remove sub dependency tree
	// TODO: remove deps from lock file

	// TODO: remove empty dir

	ctx.Conf.Save()
}

// 递归向上删除所有空文件夹
func removeAll(prefix, name string) {
	dir := filepath.Join(prefix, name)
	os.RemoveAll(dir)
	for {
		dir = filepath.Dir(dir)
		if dir == prefix || !isEmptyDir(dir) {
			break
		}

		os.RemoveAll(dir)
	}
}

// 判断是否是空文件夹
func isEmptyDir(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}
