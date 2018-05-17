package gpm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/vcs"
	"github.com/codegangsta/cli"
)

// Ctx defines the supporting context
type Ctx struct {
	*cli.Context
	*Logger
	Conf      *Config
	HomeDir   string
	CacheDir  string
	WorkDir   string
	VendorDir string
}

// NewCtx create context
func NewCtx() *Ctx {
	ctx := &Ctx{}
	ctx.Logger = NewLogger()
	ctx.Conf = NewConfig()
	ctx.init()
	return ctx
}

func (ctx *Ctx) init() {
	wd, err := os.Getwd()
	if err != nil {
		ctx.Die("cannot get work dir:%s", err)
	}

	home, _ := Home()
	ctx.HomeDir = home
	ctx.CacheDir = filepath.Join(".gpm", "cache")
	ctx.WorkDir = wd
	ctx.VendorDir = filepath.Join("vendor")
}

// Load 加载配置文件
func (ctx *Ctx) Load() {
	// 确保配置文件正常加载
	if !ctx.Conf.Exist() {
		ctx.Die("can not load cfg:%+v", ctx.Conf.ConfigPath())
	}

	ctx.Conf.Load()
}

// Get download one dependency to vendor
func (ctx *Ctx) Get(dep *Dependency) error {
	// checkout repo to cache
	local := filepath.Join(ctx.CacheDir, dep.Name)
	repo, err := vcs.NewRepo(dep.Remote(), local)
	if err != nil {
		return fmt.Errorf("repo create fail:%+v, %+v", dep.Name, err)
	}

	// TODO:check version
	ctx.Info("--> Fetching %s", dep.Name)
	if !Exists(local) {
		ctx.Debug("repo get:%+v", dep.Name)
		if err := repo.Get(); err != nil {
			return fmt.Errorf("repo get fail:%+v, %+v", dep.Name, err)
		}
	} else {
		ctx.Debug("repo up:%+v", dep.Name)
		if err := repo.Update(); err != nil {
			return fmt.Errorf("repo update fail:%+v, %+v", dep.Name, err)
		}
	}

	// export if not exist or not same version(TODO)
	relDir := filepath.Join(ctx.VendorDir, dep.Name)
	absDir, _ := filepath.Abs(relDir)
	if !Exists(absDir) {
		ctx.Info("--> Export %s, %s", dep.Name, relDir)
		if err := repo.ExportDir(absDir); err != nil {
			return fmt.Errorf("repo export fail:%+v", err)
		}
	}

	return nil
}
