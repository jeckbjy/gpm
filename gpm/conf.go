package gpm

import (
	"fmt"
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
	Version    string   `yaml:"version,omitempty"` // semantic version Semver
	Pin        string   `yaml:"-"`                 // Version for lock
	Repository string   `yaml:"repo,omitempty"`
	Vcs        string   `yaml:"vcs,omitempty"`
	Arch       []string `yaml:"arch,omitempty"`
	Os         []string `yaml:"os,omitempty"`
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
	Version string        `yaml:"version"`
	Home    string        `yaml:"home,omitempty"`
	Desc    string        `yaml:"description,omitempty"`
	License string        `yaml:"license,omitempty"`
	Owners  []*Owner      `yaml:"owners,omitempty"`
	Imports []*Dependency `yaml:"import"`
}

func NewConfig() *Config {
	cfg := &Config{}
	cfg.Init()
	return cfg
}

// Init 初始化
func (cfg *Config) Init() {
	cfg.Version = "0.0.0"
}

func (cfg *Config) ConfigPath() string {
	return ConfName
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

	// try fix name and version
	for _, dep := range cfg.Imports {
		dep.Name, dep.Version = cfg.ParseRepo(dep.Name)
	}

	return nil
}

// Save 保存配置文件
func (cfg *Config) Save() error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	filename := ConfName
	return ioutil.WriteFile(filename, data, 0666)
}

// ParseRepo 解析版本号,例如:github.com/jeckbjy/fairy@^2.0.1
func (cfg *Config) ParseRepo(repo string) (string, string) {
	tokens := strings.Split(repo, "@")
	if len(tokens) > 1 {
		return tokens[0], tokens[1]
	}

	return repo, ""
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

// DelDependency delete
func (cfg *Config) DelDependency(name string) bool {
	for i, d := range cfg.Imports {
		if d.Name == name {
			cfg.Imports = append(cfg.Imports[:i], cfg.Imports[i+1:]...)
			return true
		}
	}

	return false
}

// AddDependency add dependency to imports or dev
func (cfg *Config) AddDependency(dep *Dependency) {
	if !cfg.HasDependency(dep.Name) {
		cfg.Imports = append(cfg.Imports, dep)
	}
}

// NewDependency create Dependency if not exist
func (cfg *Config) NewDependency(name string) (*Dependency, error) {
	//
	name, version := cfg.ParseRepo(name)
	// try remove https://, git
	repo := ""
	if index := strings.Index(name, "://"); index != -1 {
		repo = name
		begIndex := index + 3
		endIndex := strings.LastIndex(name, ".")
		if endIndex == -1 {
			name = name[begIndex:]
		} else {
			name = name[begIndex:endIndex]
		}
	}

	if cfg.HasDependency(name) {
		return nil, fmt.Errorf("dependency has exist: %+v", name)
	}

	dep := &Dependency{Name: name, Version: version, Repository: repo}
	return dep, nil
}
