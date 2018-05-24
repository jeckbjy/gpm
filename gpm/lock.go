package gpm

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
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

// Load 加载lock配置
func (l *LockFile) Load() error {
	data, err := ioutil.ReadFile(LockName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, l)
}

// Save 保存lock文件
func (l *LockFile) Save() error {
	data, err := yaml.Marshal(l)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(LockName, data, 0666)
}
