package osplus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const (
	XdgConfigHomeKey string = "XDG_CONFIG_HOME"
	XdgCacheHomeKey         = "XDG_CACHE_HOME"
	XdgDataHomeKey          = "XDG_DATA_HOME"
	XdgRuntimeDirKey        = "XDG_RUNTIME_DIR"
	XdgDataDirsKey          = "XDG_DATA_DIRS"
	XdgConfigDirsKey        = "XDG_CONFIG_DIRS"
)

func Getenv(key string) (string, error) {
	v := os.Getenv(key)
	if v != "" {
		return v, nil
	}
	return "", fmt.Errorf("not found: %s", key)
}

func GetXdgConfigHome() (string, error) {
	v, err := Getenv(XdgConfigHomeKey)
	if err == nil {
		return v, nil
	}
	v, err = homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, ".config"), nil
}

func GetXdgCacheHome() (string, error) {
	v, err := Getenv(XdgCacheHomeKey)
	if err == nil {
		return v, nil
	}
	v, err = homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, ".cache"), nil
}

func GetXdgDataHome() (string, error) {
	v, err := Getenv(XdgCacheHomeKey)
	if err == nil {
		return v, nil
	}
	v, err = homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, ".local", "share"), nil
}

func GetXdgRuntimeDir() (string, error) {
	return Getenv(XdgRuntimeDirKey)
}

func GetXdgDataDirs() []string {
	v := os.Getenv(XdgDataDirsKey)
	return rejectBlank(strings.Split(v, ":"))
}

func GetXdgConfigDirs() []string {
	v := os.Getenv(XdgConfigDirsKey)
	return rejectBlank(strings.Split(v, ":"))
}
