package gpm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"

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
	gopath    string
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

// GoPath 查找可能的GOPATH
func (ctx *Ctx) GoPath() string {
	if ctx.gopath != "" {
		return ctx.gopath
	}

	ctx.gopath = ctx.findGoPath()
	return ctx.gopath
}

// FindGoPath 尝试查找gopath
func (ctx *Ctx) findGoPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// 当前目录含有src目录
	if fi, err := os.Stat("./src"); err == nil && fi.IsDir() {
		// ctx.Debug("find gopath:%+v", wd)
		return wd
	}

	// 从当前目录向上查找src,bin,pkg
	dir := wd
	for i := 0; i < 3; i++ {
		node := filepath.Base(dir)
		if node == "src" || node == "bin" || node == "pkg" {
			// 上级目录则为gopath
			return filepath.Dir(dir)
		}

		dir = filepath.Dir(dir)
	}

	return ""
}

// SetWorkDir 设置工作目录
func (ctx *Ctx) SetWorkDir(wd string) {
	os.Chdir(wd)
	ctx.WorkDir = wd
	ctx.gopath = ""
}

// AutoChangeDir 自动查找配置所在目录,并切换到此目录
func (ctx *Ctx) AutoChangeDir() bool {
	if ctx.Conf.Exist() {
		return true
	}

	gopath := ctx.GoPath()
	if gopath == "" {
		return false
	}

	root := filepath.Join(gopath, "src")
	if !Exists(root) {
		return false
	}

	wd := ""
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !fi.IsDir() && filepath.Base(path) == "gpm.yaml" {
			wd = filepath.Dir(path)
			return fmt.Errorf("find work dir")
		}

		return nil
	})

	if wd == "" {
		ctx.Info("cannot find gpm.yaml")
	}

	// ctx.Debug("change dir:%+v", wd)
	ctx.SetWorkDir(wd)

	return true
}

// Load 加载配置文件
func (ctx *Ctx) Load() {
	// 确保配置文件正常加载
	if !ctx.Conf.Exist() {
		ctx.Die("can not load cfg:%+v", ctx.Conf.ConfigPath())
	}

	ctx.Conf.Load()
}

// UpdateVersion 查找版本
func (ctx *Ctx) UpdateVersion(dep *Dependency, repo vcs.Repo) error {
	ver := dep.Version

	// References in Git can begin with a ^ which is similar to semver.
	// If there is a ^ prefix we assume it's a semver constraint rather than
	// part of the git/VCS commit id.
	if repo.IsReference(ver) && !strings.HasPrefix(ver, "^") {
		return repo.UpdateVersion(ver)
	}

	constraint, err := semver.NewConstraint(ver)
	if err != nil {
		return err
	}

	// Get the tags and branches (in that order)
	refs, err := ctx.GetAllRefs(repo)
	if err != nil {
		return err
	}

	// Convert and filter the list to semver.Version instances
	semvers := []*semver.Version{}
	for _, ref := range refs {
		v, err := semver.NewVersion(ref)
		if err == nil {
			semvers = append(semvers, v)
		}
	}

	// Sort semver list
	sort.Sort(sort.Reverse(semver.Collection(semvers)))

	found := ""
	for _, v := range semvers {
		if constraint.Check(v) {
			found = v.Original()
		}
	}

	//
	if found == "" {
		return nil
	}

	if err := repo.UpdateVersion(found); err != nil {
		return err
	}

	if found != dep.VersionLock {
		ctx.Debug("update lock version:%+v", found)
		dep.VersionLock = found
	}

	return nil
}

// UpdateSemver update Semantic Version 将最后一个版本设为默认版本
func (ctx *Ctx) UpdateSemver(dep *Dependency, repo vcs.Repo) error {
	if dep.Version != "" {
		return nil
	}

	// 最后一个版本号
	lastVersion := ""

	// Git endpoints allow for querying without fetching the codebase locally.
	// We try that first to avoid fetching right away. Is this premature
	// optimization?
	if repo.Vcs() == vcs.Git {
		pattern := regexp.MustCompile(`(?m-s)(?:tags)/(\S+)$`)
		out, err := exec.Command("git", "ls-remote", repo.Remote()).CombinedOutput()
		if err != nil {
			return err
		}

		lines := strings.Split(string(out), "\n")
		for _, i := range lines {
			ti := strings.TrimSpace(i)
			if found := pattern.FindString(ti); found != "" {
				lastVersion = strings.TrimPrefix(strings.TrimSuffix(found, "^{}"), "tags/")
			}
		}
	} else {
		if err := ctx.GetOrUpdateRepo(repo); err != nil {
			return err
		}

		tags, err := repo.Tags()
		if err != nil {
			return err
		}
		if len(tags) > 0 {
			lastVersion = tags[len(tags)-1]
		}
	}

	if lastVersion == "" {
		return fmt.Errorf("cannot find last version:%+v", repo.Remote())
	}

	if sv, err := semver.NewVersion(lastVersion); err == nil {
		// 使用语义版本号,默认使用~,只更新patch,^表示minor向下兼容更新
		dep.Version = "~" + sv.String()
	} else {
		// 非语义版本号,使用特定版本
		dep.Version = lastVersion
	}

	dep.VersionLock = lastVersion

	return nil
}

// GetAllRefs 从repo中获取所有refs
func (ctx *Ctx) GetAllRefs(repo vcs.Repo) ([]string, error) {
	tags, err := repo.Tags()
	if err != nil {
		return []string{}, err
	}

	branches, err := repo.Branches()
	if err != nil {
		return []string{}, err
	}

	refs := append(tags, branches...)
	return refs, nil
}

// GetOrUpdateRepo 获取或者更新repo
func (ctx *Ctx) GetOrUpdateRepo(repo vcs.Repo) error {
	if !Exists(repo.LocalPath()) {
		ctx.Debug("get repo:%+v", repo.Remote())
		return repo.Get()
	}

	ctx.Debug("udpate repo:%+v", repo.Remote())
	return repo.Update()
}

// Get download one dependency to vendor
func (ctx *Ctx) Get(dep *Dependency, update bool) error {
	// checkout repo to cache
	local := filepath.Join(ctx.CacheDir, dep.Name)
	repo, err := vcs.NewRepo(dep.Remote(), local)
	if err != nil {
		return fmt.Errorf("repo create fail:%+v, %+v", dep.Name, err)
	}

	ctx.Info("--> Fetching %s", dep.Name)
	if err := ctx.GetOrUpdateRepo(repo); err != nil {
		return fmt.Errorf("update repo fail:%+v", err)
	}

	if update {
		if dep.Version == "" {
			ctx.UpdateSemver(dep, repo)
		} else {
			ctx.UpdateVersion(dep, repo)
		}
	}

	// udpate lock version
	if dep.VersionLock != "" {
		needed := true
		if cur, err := repo.Current(); err == nil {
			needed = cur != dep.VersionLock
		}

		if needed {
			ctx.Info("update version:%+v", dep.VersionLock)
			if err := repo.UpdateVersion(dep.VersionLock); err != nil {
				ctx.Info("update version fail:%+v", err)
			} else {
				repo.Update()
			}
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
