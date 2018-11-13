package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

// Build cmd for build go code
type Build struct {
}

func (self *Build) Cmd() cli.Command {
	return cli.Command{
		Name:        "build",
		Usage:       "",
		Description: "",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name, n",
				Usage: "build target name",
			},
			cli.StringFlag{
				Name:  "gopath, g",
				Usage: "gopath for build",
			},
		},
	}
}

// Run build go code to bin
// build -n xxxxx
func (self *Build) Run(ctx *gpm.Ctx) {
	ctx.MustLoad()

	// get target name
	name := ctx.String("name")
	if name == "" {
		name = ctx.Name
		if name == "" || name == "." {
			dir, err := os.Getwd()
			if err == nil {
				name = filepath.Base(dir)
			}
		}

		if name == "" {
			name = "server"
		}
	}

	// get GOPATH,USE system gopath?
	gopath := ctx.String("gopath")
	if gopath == "" {
		gopath = self.findGoPath()
	}

	if gopath == "" {
		ctx.Die("cannot find gopath")
	}

	// build
	envGopath := fmt.Sprintf("GOPATH=%s", gopath)
	target := filepath.Join(gopath, "bin", name)
	ctx.Info("--> Build:%+v", name)
	cmd := exec.Command("go", "build", "-o", target)

	cmd.Env = append(cmd.Env, envGopath)
	output, err := cmd.CombinedOutput()
	if err == nil {
		ctx.Info("    build ok, output is %+v", filepath.Join("bin", name))
	} else {
		ctx.Info("    build fail, %s", output)
	}
}

// 查找gopath,三级目录内包含bin或者src
func (self *Build) findGoPath() string {
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
	// src/github.com/user/project
	dir := wd
	for i := 0; i < 3; i++ {
		node := filepath.Base(dir)
		if node == "src" || node == "bin" {
			// 上级目录则为gopath
			return filepath.Dir(dir)
		}

		dir = filepath.Dir(dir)
	}

	return ""
}
