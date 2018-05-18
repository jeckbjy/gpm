package cmd

import (
	"fmt"
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
		},
	}
}

// Run build go code to bin
func (self *Build) Run(ctx *gpm.Ctx) {
	ctx.AutoChangeDir()
	ctx.Load()
	count := len(ctx.Args())
	name := ctx.String("name")
	if name == "" {
		name = ctx.Conf.Name
		if name == "" || name == "." {
			name = ctx.WorkDir
		}
		name = filepath.Base(name)
	}

	if count == 0 {
		gopath := ctx.GoPath()
		if gopath == "" {
			ctx.Die("cannot detect gopath")
		}

		// ctx.Debug("gopath:%+v", gopath)

		envGopath := fmt.Sprintf("GOPATH=%s", gopath)
		target := filepath.Join(gopath, "bin", name)
		ctx.Info("--> Build:%+v", name)
		cmd := exec.Command("go", "build", "-o", target)

		cmd.Env = append(cmd.Env, envGopath)
		output, err := cmd.CombinedOutput()
		if err == nil {
			ctx.Info("    build succeed:output is %+v", filepath.Join("bin", name))
		} else {
			ctx.Info("    build fail, %s", output)
		}
	}
}
