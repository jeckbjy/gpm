package gpm

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/Masterminds/vcs"
	"github.com/codegangsta/cli"
)

const (
	PREFIX_GIT   = "git@"
	PREFIX_HTTPS = "https://"
)

type Ctx struct {
	*cli.Context
	*Logger
	*Config
	CacheDir string
}

// NewCtx create context
func NewCtx() *Ctx {
	ctx := &Ctx{}
	ctx.init()
	return ctx
}

func (ctx *Ctx) init() {
	ctx.Logger = NewLogger()
	ctx.Config = NewConfig()
	if home, err := Home(); err != nil {
		ctx.Die("cannot get home")
		ctx.CacheDir = filepath.Join(home, ".gpm")
	}
}

// MustLoad load config and die if not exists
func (ctx *Ctx) MustLoad() {
	err := ctx.Load()
	if err != nil {
		ctx.Die("not find config,use gpm init to create")
	}
}

// Get 类似go get,获取代码放入vendor中
func (ctx *Ctx) Get(url string, version string) error {
	if !strings.HasPrefix(url, PREFIX_GIT) && !strings.Contains(url, "://") {
		url += "https://"
	}

	local, err := CacheLocal(url)
	if err != nil {
		return err
	}

	// step1: get repo
	repo, err := vcs.NewRepo(url, local)
	if err != nil {
		return err
	}

	if !Exists(repo.LocalPath()) {
		err = repo.Get()
	} else {
		err = repo.Update()
	}

	if err != nil {
		return err
	}

	// step2: update version
	oldVersion, _ := repo.Current()
	ctx.UpdateVersion(repo, version)
	newVersion, _ := repo.Current()

	// step3: update repo
	if err := repo.Update(); err != nil {
		return err
	}

	// step4: export repo to vendor
	// git@github.com:user/repo.git -> github.com/user/repo
	// https://github.com/user/repo
	path := url
	if strings.HasPrefix(url, PREFIX_GIT) {
		path = path[len(PREFIX_GIT):]
		path = strings.TrimRight(path, ".git")
		path = strings.Replace(path, ":", "/", 1)
	} else if strings.HasPrefix(url, PREFIX_HTTPS) {
		path = path[len(PREFIX_HTTPS):]
	}

	dir := filepath.Join("vendor", path)
	exportDir, _ := filepath.Abs(dir)
	if !Exists(exportDir) || oldVersion != newVersion {
		ctx.Info("--> Export %s, %s", path, exportDir)
		if err := repo.ExportDir(exportDir); err != nil {
			return fmt.Errorf("repo export fail:%+v", err)
		}
	}

	return nil
}

// UpdateVersion update repo version
func (ctx *Ctx) UpdateVersion(repo vcs.Repo, ver string) error {
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
	refs := []string{}
	if tags, err := repo.Tags(); err == nil {
		refs = tags
	}

	if branches, err := repo.Branches(); err == nil {
		refs = append(refs, branches...)
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

	return nil
}

// const (
// 	GetModeInit = iota
// 	GetModeUpdate
// 	GetModeInstall
// )

// // Ctx 工作上下文，所有的操作都是相对于配置文件所在的目录设置的
// type Ctx struct {
// 	*cli.Context
// 	*Logger
// 	Conf     *Config
// 	HomeDir  string
// 	CacheDir string
// 	WorkDir  string
// 	gopath   string
// }

// // NewCtx create context
// func NewCtx() *Ctx {
// 	ctx := &Ctx{}
// 	ctx.Logger = NewLogger()
// 	ctx.Conf = NewConfig()
// 	ctx.init()
// 	return ctx
// }

// func (ctx *Ctx) init() {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		ctx.Die("cannot get work dir:%s", err)
// 	}

// 	home, _ := Home()
// 	ctx.HomeDir = home
// 	ctx.CacheDir = filepath.Join(home, ".gpm")
// 	ctx.WorkDir = wd
// }

// // GoPath 查找可能的GOPATH
// func (ctx *Ctx) GoPath() string {
// 	if ctx.gopath != "" {
// 		return ctx.gopath
// 	}

// 	ctx.gopath = ctx.findGoPath()
// 	return ctx.gopath
// }

// // VendorDir 默认和配置同级目录,可以在配置中修改
// func (ctx *Ctx) VendorDir() string {
// 	if ctx.Conf.Vendor != "" {
// 		if strings.HasSuffix(ctx.Conf.Vendor, "/vendor") {
// 			return ctx.Conf.Vendor
// 		} else {
// 			return filepath.Join(ctx.Conf.Vendor, "vendor")
// 		}
// 	}

// 	return "vendor"
// }

// // FindGoPath 尝试查找gopath
// func (ctx *Ctx) findGoPath() string {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return ""
// 	}

// 	// 当前目录含有src目录
// 	if fi, err := os.Stat("./src"); err == nil && fi.IsDir() {
// 		// ctx.Debug("find gopath:%+v", wd)
// 		return wd
// 	}

// 	// 从当前目录向上查找src,bin,pkg
// 	dir := wd
// 	for i := 0; i < 3; i++ {
// 		node := filepath.Base(dir)
// 		if node == "src" || node == "bin" || node == "pkg" {
// 			// 上级目录则为gopath
// 			return filepath.Dir(dir)
// 		}

// 		dir = filepath.Dir(dir)
// 	}

// 	return ""
// }

// // SetWorkDir 设置工作目录
// func (ctx *Ctx) SetWorkDir(wd string) {
// 	os.Chdir(wd)
// 	ctx.WorkDir = wd
// 	ctx.gopath = ""
// }

// // AutoChangeDir 自动查找配置所在目录,并切换到此目录
// func (ctx *Ctx) AutoChangeDir() bool {
// 	if ctx.Conf.Exist() {
// 		return true
// 	}

// 	gopath := ctx.GoPath()
// 	if gopath == "" {
// 		return false
// 	}

// 	root := filepath.Join(gopath, "src")
// 	if !Exists(root) {
// 		return false
// 	}

// 	wd := ""
// 	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
// 		if err != nil {
// 			return nil
// 		}

// 		if !fi.IsDir() && filepath.Base(path) == "gpm.yaml" {
// 			wd = filepath.Dir(path)
// 			return fmt.Errorf("find work dir")
// 		}

// 		return nil
// 	})

// 	if wd == "" {
// 		ctx.Info("cannot find gpm.yaml")
// 	}

// 	// ctx.Debug("change dir:%+v", wd)
// 	ctx.SetWorkDir(wd)

// 	return true
// }

// // Load 加载配置文件
// func (ctx *Ctx) Load() {
// 	// 确保配置文件正常加载
// 	if !ctx.Conf.Exist() {
// 		ctx.Die("can not load cfg:%+v", ctx.Conf.ConfigPath())
// 	}

// 	if err := ctx.Conf.Load(); err != nil {
// 		ctx.Die("load config fail:%+v", err)
// 	}

// 	if !Exists(LockName) {
// 		return
// 	}

// 	lockfile := &LockFile{}
// 	if err := lockfile.Load(); err != nil {
// 		ctx.Die("load lock fail:%+v", err)
// 	}

// 	ctx.Conf.ImportLock(lockfile)
// }

// // Save 保存配置文件和lock文件
// func (ctx *Ctx) Save() {
// 	ctx.Conf.Save()

// 	if len(ctx.Conf.Imports) > 0 || Exists(LockName) {
// 		l := &LockFile{}
// 		ctx.Conf.ExportLock(l)
// 		l.Save()
// 	}
// }

// // UpdateVersion 查找版本
// func (ctx *Ctx) UpdateVersion(dep *Dependency, repo vcs.Repo) error {
// 	ver := dep.Version

// 	// References in Git can begin with a ^ which is similar to semver.
// 	// If there is a ^ prefix we assume it's a semver constraint rather than
// 	// part of the git/VCS commit id.
// 	if repo.IsReference(ver) && !strings.HasPrefix(ver, "^") {
// 		return repo.UpdateVersion(ver)
// 	}

// 	constraint, err := semver.NewConstraint(ver)
// 	if err != nil {
// 		return err
// 	}

// 	// Get the tags and branches (in that order)
// 	refs, err := ctx.GetAllRefs(repo)
// 	if err != nil {
// 		return err
// 	}

// 	// Convert and filter the list to semver.Version instances
// 	semvers := []*semver.Version{}
// 	for _, ref := range refs {
// 		v, err := semver.NewVersion(ref)
// 		if err == nil {
// 			semvers = append(semvers, v)
// 		}
// 	}

// 	// Sort semver list
// 	sort.Sort(sort.Reverse(semver.Collection(semvers)))

// 	found := ""
// 	for _, v := range semvers {
// 		if constraint.Check(v) {
// 			found = v.Original()
// 		}
// 	}

// 	//
// 	if found == "" {
// 		return nil
// 	}

// 	if err := repo.UpdateVersion(found); err != nil {
// 		return err
// 	}

// 	if found != dep.Reversion {
// 		dep.Reversion = found
// 	}

// 	return nil
// }

// // UpdateSemver update Semantic Version 将最后一个版本设为默认版本
// func (ctx *Ctx) UpdateSemver(dep *Dependency, repo vcs.Repo) error {
// 	if dep.Version != "" {
// 		return nil
// 	}

// 	// 最后一个版本号
// 	lastVersion := ""

// 	// Git endpoints allow for querying without fetching the codebase locally.
// 	// We try that first to avoid fetching right away. Is this premature
// 	// optimization?
// 	if repo.Vcs() == vcs.Git {
// 		pattern := regexp.MustCompile(`(?m-s)(?:tags)/(\S+)$`)
// 		out, err := exec.Command("git", "ls-remote", repo.Remote()).CombinedOutput()
// 		if err != nil {
// 			return err
// 		}

// 		lines := strings.Split(string(out), "\n")
// 		for _, i := range lines {
// 			ti := strings.TrimSpace(i)
// 			if found := pattern.FindString(ti); found != "" {
// 				lastVersion = strings.TrimPrefix(strings.TrimSuffix(found, "^{}"), "tags/")
// 			}
// 		}
// 	} else {
// 		if err := ctx.GetOrUpdateRepo(repo); err != nil {
// 			return err
// 		}

// 		tags, err := repo.Tags()
// 		if err != nil {
// 			return err
// 		}
// 		if len(tags) > 0 {
// 			lastVersion = tags[len(tags)-1]
// 		}
// 	}

// 	if lastVersion == "" {
// 		return fmt.Errorf("cannot find last version:%+v", repo.Remote())
// 	}

// 	if sv, err := semver.NewVersion(lastVersion); err == nil {
// 		// 使用语义版本号,默认使用~,只更新patch,^表示minor向下兼容更新
// 		dep.Version = "~" + sv.String()
// 	} else {
// 		// 非语义版本号,使用特定版本
// 		dep.Version = lastVersion
// 	}

// 	dep.Reversion = lastVersion

// 	return nil
// }

// // GetAllRefs 从repo中获取所有refs
// func (ctx *Ctx) GetAllRefs(repo vcs.Repo) ([]string, error) {
// 	tags, err := repo.Tags()
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	branches, err := repo.Branches()
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	refs := append(tags, branches...)
// 	return refs, nil
// }

// // GetOrUpdateRepo 获取或者更新repo
// func (ctx *Ctx) GetOrUpdateRepo(repo vcs.Repo) error {
// 	if !Exists(repo.LocalPath()) {
// 		return repo.Get()
// 	}

// 	return repo.Update()
// }

// // scpSyntaxRe matches the SCP-like addresses used to access repos over SSH.
// var scpSyntaxRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

// // GetCacheLocal 返回一个cache地址
// func (ctx *Ctx) GetCacheLocal(repo string) (string, error) {
// 	var u *url.URL
// 	var err error
// 	var strip bool
// 	if m := scpSyntaxRe.FindStringSubmatch(repo); m != nil {
// 		// Match SCP-like syntax and convert it to a URL.
// 		// Eg, "git@github.com:user/repo" becomes
// 		// "ssh://git@github.com/user/repo".
// 		u = &url.URL{
// 			Scheme: "ssh",
// 			User:   url.User(m[1]),
// 			Host:   m[2],
// 			Path:   "/" + m[3],
// 		}
// 		strip = true
// 	} else {
// 		u, err = url.Parse(repo)
// 		if err != nil {
// 			return "", err
// 		}
// 	}

// 	if strip {
// 		u.Scheme = ""
// 	}

// 	var key string
// 	if u.Scheme != "" {
// 		key = u.Scheme + "-"
// 	}
// 	if u.User != nil && u.User.Username() != "" {
// 		key = key + u.User.Username() + "-"
// 	}
// 	key = key + u.Host
// 	if u.Path != "" {
// 		key = key + strings.Replace(u.Path, "/", "-", -1)
// 	}

// 	key = strings.Replace(key, ":", "-", -1)

// 	return filepath.Join(ctx.CacheDir, key), nil
// }

// // Get download one dependency to vendor
// func (ctx *Ctx) Get(dep *Dependency, mode int) error {
// 	// checkout repo to cache
// 	local, err := ctx.GetCacheLocal(dep.Name)
// 	if err != nil {
// 		return err
// 	}

// 	repo, err := vcs.NewRepo(dep.Remote(), local)
// 	if err != nil {
// 		return fmt.Errorf("repo create fail:%+v, %+v", dep.Name, err)
// 	}

// 	oldVersion, _ := repo.Current()

// 	ctx.Info("--> Fetching %s", dep.Name)
// 	if err := ctx.GetOrUpdateRepo(repo); err != nil {
// 		return fmt.Errorf("update repo fail:%+v", err)
// 	}

// 	// 自动获取Semantic Version
// 	if mode == GetModeInit {
// 		if dep.Version == "" {
// 			ctx.UpdateSemver(dep, repo)
// 		}
// 	}

// 	// 查找合适的版本
// 	if mode == GetModeInit || mode == GetModeUpdate {
// 		ctx.UpdateVersion(dep, repo)
// 	}

// 	// 更新到合适的版本
// 	if dep.Reversion != "" {
// 		needed := true
// 		if cur, err := repo.Current(); err == nil {
// 			needed = cur != dep.Reversion
// 		}

// 		if needed {
// 			ctx.Info("update version:%+v", dep.Reversion)
// 			if err := repo.UpdateVersion(dep.Reversion); err != nil {
// 				ctx.Info("update version fail:%+v", err)
// 			} else {
// 				repo.Update()
// 			}
// 		}
// 	}

// 	// TODO:记录版本信息
// 	// export if not exist or not same version(TODO)
// 	relDir := filepath.Join(ctx.VendorDir(), dep.Name)
// 	absDir, _ := filepath.Abs(relDir)
// 	curVersion, _ := repo.Current()
// 	if oldVersion != curVersion || !Exists(absDir) {
// 		ctx.Info("--> Export %s, %s", dep.Name, relDir)
// 		if err := repo.ExportDir(absDir); err != nil {
// 			return fmt.Errorf("repo export fail:%+v", err)
// 		}
// 	}

// 	return nil
// }
