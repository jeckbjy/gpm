package gpm

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const VendorDir = "vendor"
const ConfFile = "gpm.yaml"
const LockFile = "gpm.lock"

// Owner describes an owner of a package. This can be a person, company, or
// other organization. This is useful if someone needs to contact the
// owner of a package to address things like a security issue.
type Owner struct {
	Name  string `yaml:"name,omitempty"`
	Email string `yaml:"email,omitempty"`
	Home  string `yaml:"home,omitempty"`
}

// Dependency describes a package that the present package depends upon.
type Dependency struct {
	Name       string   `yaml:package`
	Version    string   `yaml:version,omitempty` // Version
	Pin        string   `yaml:"-"`               // Version fot lock
	Repository string   `yaml:repo,omitempty`
	VCS        string   `yaml:vcs,omitempty`
	Arch       []string `yaml:arch,omitempty`
	Os         []string `yaml:os,omitempty`
}

func (d *Dependency) Remote() string {
	var r string

	if d.Repository != "" {
		r = d.Repository
	} else {
		r = "https://" + d.Name
	}

	return r
}

// Config is the top-level configuration object.
type Config struct {
	Name    string        `yaml:"package"`
	Desc    string        `yaml:"description,omitempty"`
	Home    string        `yaml:"homepage,omitempty"`
	License string        `yaml:"license,omitempty"`
	Owners  []*Owner      `yaml:"owners,omitempty"`
	Imports []*Dependency `yaml:"import"`
}

// Exist 判断配置文件是否存在
func (cfg *Config) Exist() bool {
	_, err := os.Stat(ConfFile)
	return err == nil
}

// Create 初始化信息
func (cfg *Config) Create() {
}

// Load 加载配置文件
func (cfg *Config) Load() error {
	data, err := ioutil.ReadFile(ConfFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

// Save 保存配置文件
func (cfg *Config) Save() error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ConfFile, data, 0666)
}
