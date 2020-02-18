package ose

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

type EnvPath struct {
	fs  afero.Fs
	env Env
}

func NewEnvPath(fs afero.Fs, env Env) *EnvPath {
	return &EnvPath{fs: fs, env: env}
}

// GOPATH, PATH

func (p *EnvPath) GetGoPath() (string, error) {
	vs := p.GetGoPathMulti()
	if len(vs) > 0 {
		return vs[0], nil
	}
	v, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, "go"), nil
}

func (p *EnvPath) GetGoPathMulti() []string {
	v := p.env.Get("GOPATH")
	return rejectEmpty(strings.Split(v, ":"))
}

func (p *EnvPath) GetPath() []string {
	v := p.env.Get("PATH")
	return rejectEmpty(strings.Split(v, ":"))
}

// XDG_*

const (
	XdgConfigHomeKey string = "XDG_CONFIG_HOME"
	XdgCacheHomeKey         = "XDG_CACHE_HOME"
	XdgDataHomeKey          = "XDG_DATA_HOME"
	XdgRuntimeDirKey        = "XDG_RUNTIME_DIR"
	XdgDataDirsKey          = "XDG_DATA_DIRS"
	XdgConfigDirsKey        = "XDG_CONFIG_DIRS"
)

func (p *EnvPath) GetXdgConfigHome() (string, error) {
	if v, ok := p.env.Lookup(XdgConfigHomeKey); ok {
		return v, nil
	}
	v, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, ".config"), nil
}

func (p *EnvPath) GetXdgCacheHome() (string, error) {
	if v, ok := p.env.Lookup(XdgCacheHomeKey); ok {
		return v, nil
	}
	v, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, ".cache"), nil
}

func (p *EnvPath) GetXdgDataHome() (string, error) {
	if v, ok := p.env.Lookup(XdgCacheHomeKey); ok {
		return v, nil
	}
	v, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, ".local", "share"), nil
}

func (p *EnvPath) GetXdgRuntimeDir() (string, error) {
	if v, ok := p.env.Lookup(XdgRuntimeDirKey); ok {
		return v, nil
	}
	return "", fmt.Errorf("not found: %s", XdgRuntimeDirKey)
}

func (p *EnvPath) GetXdgDataDirs() []string {
	v := p.env.Get(XdgDataDirsKey)
	return rejectEmpty(strings.Split(v, ":"))
}

func (p *EnvPath) GetXdgConfigDirs() []string {
	v := p.env.Get(XdgConfigDirsKey)
	return rejectEmpty(strings.Split(v, ":"))
}

func (p *EnvPath) LookPathWithPredicate(dirs []string, names []string, pred func(fpath string, fi os.FileInfo) (ok bool)) (string, error) {
	for _, dir := range dirs {
		for _, name := range names {
			fpath := filepath.Join(dir, name)
			fi, err := os.Stat(fpath)
			if err != nil {
				continue
			}
			if pred(fpath, fi) {
				return fpath, nil
			}
		}
	}
	return "", fmt.Errorf("not found: %s in %s", strings.Join(names, ", "), strings.Join(dirs, ", "))
}

func (p *EnvPath) LookPath(dirs []string, names ...string) (string, error) {
	return p.LookPathWithPredicate(dirs, names, func(_ string, _ os.FileInfo) bool { return true })
}

func (p *EnvPath) LookPathAll(dirs []string, names ...string) []string {
	results := make([]string, 0)
	p.LookPathWithPredicate(dirs, names, func(fpath string, _ os.FileInfo) bool {
		results = append(results, fpath)
		return false
	})
	return results
}
