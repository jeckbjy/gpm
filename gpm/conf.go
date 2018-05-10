package gpm

import (
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	VendorName = "vendor"
	ConfName   = "gpm.yaml"
	LockName   = "gpm.lock"
)

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
	Name       string   `yaml:"package"`
	Version    string   `yaml:"version,omitempty"` // Version
	Pin        string   `yaml:"-"`                 // Version fot lock
	Repository string   `yaml:"repo,omitempty"`
	VCS        string   `yaml:"vcs,omitempty"`
	Arch       []string `yaml:"arch,omitempty"`
	Os         []string `yaml:"os,omitempty"`
}

// Parse name and version
func (d *Dependency) Parse() {
	parts := strings.Split(d.Name, "#")
	if len(parts) > 1 {
		d.Name = parts[0]
		d.Version = parts[1]
	}
}

// Remote returns the remote location to fetch source from. This location is
// the central place where mirrors can alter the location.
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
	Path    string        `yaml:"-"` // 配置文件所在目录
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.Init()
	return cfg
}

// Init 初始化
func (cfg *Config) Init() {
	cfg.Name = "."
}

func (cfg *Config) SetPath(dir string) {
	if !strings.HasSuffix(dir, "/") {
		cfg.Path = dir + "/"
	} else {
		cfg.Path = dir
	}
}

func (cfg *Config) ConfigPath() string {
	return cfg.Path + ConfName
}

// Exist 判断配置文件是否存在
func (cfg *Config) Exist() bool {
	filename := cfg.ConfigPath()
	_, err := os.Stat(filename)
	return err == nil
}

// Load 加载配置文件
func (cfg *Config) Load() error {
	filename := cfg.ConfigPath()
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}

	// parse name and version
	for _, dep := range cfg.Imports {
		dep.Parse()
	}

	return nil
}

// Save 保存配置文件
func (cfg *Config) Save() error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	filename := cfg.Path + ConfName
	return ioutil.WriteFile(filename, data, 0666)
}

// HasDependency returns true if the given name is listed as an import or dev import.
func (cfg *Config) HasDependency(name string) bool {
	for _, d := range cfg.Imports {
		if d.Name == name {
			return true
		}
	}

	return false
}

// AddDependency returns true if the given name is listed as an import or dev import.
func (cfg *Config) AddDependency(name string) bool {
	if cfg.HasDependency(name) {
		return false
	}

	dep := &Dependency{Name: name}
	dep.Parse()
	cfg.Imports = append(cfg.Imports, dep)

	return true
}
