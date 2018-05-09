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

func (self *Create) Run(ctx *gpm.Ctx) {
	if ctx.Conf.Exist() {
		ctx.Die("Cowardly refusing to overwrite existing YAML.")
	}

	ctx.Info("Writing configuration")
	ctx.Conf.Create()
	ctx.Conf.Save()
}

// Create 创建配置文件
// func Create() {
// 	// if conf.Exist() {
// 	// 	msg.Die("Cowardly refusing to overwrite existing YAML.")
// 	// }

// 	// msg.Info("Writing configuration")
// 	// cfg := conf.New()
// 	// cfg.Save()
// }

// func buildConfig(base string) *conf.Config {
// 	builder, err := util.GetBuildContext()
// 	if err != nil {
// 		msg.Die("Failed to build an import context: %s", err)
// 	}

// 	name := builder.PackageName(base)

// 	config := new(conf.Config)
// 	config.Name = name

// 	return config
// }
