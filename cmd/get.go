package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/jeckbjy/gpm/gpm"
)

type Get struct {
}

func (self *Get) Cmd() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Install one or more packages into `vendor/` and add dependency to gpm.yaml.",
		Description: `Gets one or more package (like 'go get') and then adds that file
		to the glide.yaml file. Multiple package names can be specified on one line.
	 
			$ glide get github.com/Masterminds/cookoo/web
	 
		The above will install the project github.com/Masterminds/cookoo and add
		the subpackage 'web'.
	 
		If a fetched dependency has a glide.yaml file, configuration from Godep,
		GPM, GOM, or GB Glide that configuration will be used to find the dependencies
		and versions to fetch. If those are not available the dependent packages will
		be fetched as either a version specified elsewhere or the latest version.
	 
		When adding a new dependency Glide will perform an update to work out
		the versions for the dependencies of this dependency (transitive ones). This
		will generate an updated glide.lock file with specific locked versions to use.
	 
		The '--strip-vendor' flag will remove any nested 'vendor' folders and
		'Godeps/_workspace' folders after an update (along with undoing any Godep
		import rewriting). Note, The Godeps specific functionality is deprecated and
		will be removed when most Godeps users have migrated to using the vendor
		folder.`,
	}
}

func (self *Get) Run(ctx *gpm.Ctx) {

}

// Get 获取代码
// func Get(names []string) {

// 	// remote := "https://github.com/Masterminds/vcs"
// 	// // local, _ := ioutil.TempDir("", "go-vcs")
// 	// local := "aa"
// 	// repo, err := vcs.NewRepo(remote, local)
// 	// msg.Info("get data:%+v", local)

// 	// if err != nil {
// 	// 	msg.Die("bad repo:%+v", err)
// 	// 	return
// 	// }

// 	// repo.Get()

// 	// repo.ExportDir("github.com/Masterminds/vcs")

// 	// EnsureVendor()

// 	// base := conf.BasePath()
// 	// EnsureVendorDir()
// 	// conf := EnsureConfig()
// }
