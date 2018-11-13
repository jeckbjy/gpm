package gpm

import (
	"bytes"
	"errors"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Exists check dir of file exists
func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		home := drive + path
		if drive == "" || path == "" {
			home = os.Getenv("USERPROFILE")
		}
		if home == "" {
			return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
		}

		return home, nil
	} else {
		// Unix-like system, so just assume Unix
		// First prefer the HOME environmental variable
		if home := os.Getenv("HOME"); home != "" {
			return home, nil
		}

		// If that fails, try the shell
		var stdout bytes.Buffer
		cmd := exec.Command("sh", "-c", "eval echo ~$USER")
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			return "", err
		}

		result := strings.TrimSpace(stdout.String())
		if result == "" {
			return "", errors.New("blank output when reading home directory")
		}

		return result, nil
	}
}

// scpSyntaxRe matches the SCP-like addresses used to access repos over SSH.
var scpSyntaxRe = regexp.MustCompile(`^([a-zA-Z0-9_]+)@([a-zA-Z0-9._-]+):(.*)$`)

// CacheLocal return local path for cache
func CacheLocal(repo string) (string, error) {
	var u *url.URL
	var err error
	var strip bool
	if m := scpSyntaxRe.FindStringSubmatch(repo); m != nil {
		// Match SCP-like syntax and convert it to a URL.
		// Eg, "git@github.com:user/repo" becomes
		// "ssh://git@github.com/user/repo".
		u = &url.URL{
			Scheme: "ssh",
			User:   url.User(m[1]),
			Host:   m[2],
			Path:   "/" + m[3],
		}
		strip = true
	} else {
		u, err = url.Parse(repo)
		if err != nil {
			return "", err
		}
	}

	if strip {
		u.Scheme = ""
	}

	var key string
	if u.Scheme != "" {
		key = u.Scheme + "-"
	}
	if u.User != nil && u.User.Username() != "" {
		key = key + u.User.Username() + "-"
	}
	key = key + u.Host
	if u.Path != "" {
		key = key + strings.Replace(u.Path, "/", "-", -1)
	}

	key = strings.Replace(key, ":", "-", -1)

	home, err := Home()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".gpm", key), nil
}

// ParseRepo repo = url@version
func ParseRepo(repo string) (string, string) {
	url := repo
	// ignore "git@"
	if strings.HasPrefix(repo, "git@") {
		url = repo[4:]
	}

	if strings.Contains(url, "@") {
		tokens := strings.SplitAfterN(repo, "@", 1)
		return tokens[0], tokens[1]
	}

	return repo, ""
}
