package gpm

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	// ConfName config file name
	ConfName = "gpm.yaml"
	// LockName lock file name
	LockName = "gpm.lock"
)

// Lock represents an individual locked dependency.
type Lock struct {
	Name      string
	Reversion string
}

// LockFile represents a gpm.lock file.
type LockFile struct {
	Imports []*Lock
}

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
	Name       string `yaml:"package"`
	Version    string `yaml:"version,omitempty"` // semantic version
	Repository string `yaml:"repo,omitempty"`    //
	Reversion  string `yaml:"-"`                 // Version for lock
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

// NewDependency create dependency
func NewDependency(repo string) (*Dependency, error) {
	if strings.HasPrefix(repo, "git@") {
		return nil, fmt.Errorf("not support git@")
	}

	name, version := "", ""
	if strings.Contains(repo, "@") {
		tokens := strings.Split(repo, "@")
		name = tokens[0]
		version = tokens[1]
	} else {
		name = repo
	}

	// remove https://
	if strings.HasPrefix(name, "https://") {
		name = name[len("https://"):]
	}

	dep := &Dependency{Name: name, Version: version}
	return dep, nil
}

// NewConfig create config
func NewConfig() *Config {
	cfg := &Config{}
	cfg.init()
	return cfg
}

// Init 初始化
func (cfg *Config) init() {
	cfg.Version = "0.0.0"
}

// Exist 判断配置文件是否存在
func (cfg *Config) Exist() bool {
	_, err := os.Stat(ConfName)
	return err == nil
}

// Load 加载配置文件
func (cfg *Config) Load() error {
	data, err := ioutil.ReadFile(ConfName)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}

	// try fix name and version
	for _, dep := range cfg.Imports {
		if dep.Version == "" && strings.Contains(dep.Name, "@") {
			tokens := strings.Split(dep.Name, "@")
			dep.Name = tokens[0]
			dep.Version = tokens[1]
		}
	}

	return nil
}

// Save 保存配置文件
func (cfg *Config) Save() error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(ConfName, data, 0666)
}

// func (cfg *Config) LoadLock() error {
// 	l := &LockFile{}
// 	data, err := ioutil.ReadFile(LockName)
// 	if err != nil {
// 		return err
// 	}

// 	return yaml.Unmarshal(data, l)
// }

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
// func (cfg *Config) NewDependency(name string) (*Dependency, error) {
// 	//
// 	name, version := cfg.ParseRepo(name)
// 	// try remove https://, git
// 	repo := ""
// 	if index := strings.Index(name, "://"); index != -1 {
// 		repo = name
// 		begIndex := index + 3
// 		endIndex := strings.LastIndex(name, ".")
// 		if endIndex == -1 {
// 			name = name[begIndex:]
// 		} else {
// 			name = name[begIndex:endIndex]
// 		}
// 	}

// 	if cfg.HasDependency(name) {
// 		return nil, fmt.Errorf("dependency has exist: %+v", name)
// 	}

// 	dep := &Dependency{Name: name, Version: version, Repository: repo}
// 	return dep, nil
// }

// func (cfg *Config) findDependency(deps []*Dependency, name string) *Dependency {
// 	for _, dep := range deps {
// 		if dep.Name == name {
// 			return dep
// 		}
// 	}

// 	return nil
// }

// // ImportLock 从lock文件导入
// func (cfg *Config) ImportLock(l *LockFile) {
// 	for _, ldep := range l.Imports {
// 		cdep := cfg.findDependency(cfg.Imports, ldep.Name)
// 		if cdep != nil {
// 			cdep.Reversion = ldep.Reversion
// 		}
// 	}
// }

// // ExportLock 导出lock文件
// func (cfg *Config) ExportLock(l *LockFile) {
// 	l.Imports = l.Imports[:0]
// 	for _, dep := range cfg.Imports {
// 		ldep := &Lock{Name: dep.Name, Reversion: dep.Reversion}
// 		l.Imports = append(l.Imports, ldep)
// 	}
// }
