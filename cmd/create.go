package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Create struct {
}

func (self *Create) Cmd() cli.Command {
	return cli.Command{
		Name:      "create",
		ShortName: "init",
		Usage:     "Initialize a new project, creating a gpm.yaml file",
		Description: `This command starts from a project without Glide and
		sets it up. It generates a glide.yaml file, parsing your codebase to guess
		the dependencies to include. Once this step is done you may edit the
		glide.yaml file to update imported dependency properties such as the version
		or version range to include.
	
		To fetch the dependencies you may run 'glide install'.`,
	}
}

// 在当前目录创建gpm.yaml
func (self *Create) Run(ctx *gpm.Ctx) {
	if len(ctx.Args()) > 1 {
		ctx.Die("don't support")
		return
	}

	// if len(ctx.Args()) == 1 {
	// 	// 在子目录创建 src/xxx/xxx
	// 	name := ctx.Args()[0]
	// 	if strings.HasPrefix(name, "src/") {
	// 		name = name[4:]
	// 	}

	// 	dir := name
	// 	if filepath.Base(ctx.WorkDir) != "src" {
	// 		dir = filepath.Join("src", name)
	// 	}

	// 	if gpm.Exists(dir) {
	// 		ctx.Die("folder exists:%+v", dir)
	// 	}

	// 	if err := os.MkdirAll(dir, 0755); err != nil {
	// 		ctx.Die("can not create path")
	// 	}

	// 	// 修改工作目录
	// 	if wd, err := filepath.Abs(dir); err == nil {
	// 		ctx.SetWorkDir(wd)
	// 	} else {
	// 		ctx.Die("change work dir fail:%+v", dir)
	// 	}

	// 	ctx.Conf.Name = name
	// } else {
	// 	// 在当前目录创建
	// 	if ctx.Conf.Exist() {
	// 		ctx.Die("Cowardly refusing to overwrite existing YAML.")
	// 	}

	// 	ctx.Conf.Name = filepath.Base(ctx.WorkDir)
	// }

	// dir, err := filepath.Rel(ctx.GoPath(), ctx.WorkDir)
	// if err != nil {
	// 	ctx.Die("cannot parse dir:%+v", err)
	// }

	if ctx.Exist() {
		ctx.Die("Cowardly refusing to overwrite existing YAML.")
	}

	ctx.Info("Writing configuration gpm.yaml")
	ctx.Save()
}
