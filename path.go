package osplus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func getenv(key string) (string, error) {
	v := os.Getenv(key)
	if v != "" {
		return v, nil
	}
	return "", fmt.Errorf("not found: %s", key)
}

func rejectBlank(ss []string) []string {
	results := make([]string, 0)
	for _, s := range ss {
		if s != "" {
			results = append(results, s)
		}
	}
	return results
}

func GetGoPath() (string, error) {
	vs := GetGoPathMulti()
	if len(vs) > 0 {
		return vs[0], nil
	}
	v, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(v, "go"), nil
}

func GetGoPathMulti() []string {
	v := os.Getenv("GOPATH")
	return rejectBlank(strings.Split(v, ":"))
}

func GetPath() []string {
	v := os.Getenv("PATH")
	return rejectBlank(strings.Split(v, ":"))
}

func LookPathWithPredicate(dirs []string, names []string, pred func(fpath string, fi os.FileInfo) (ok bool)) (string, error) {
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

func LookPath(dirs []string, names ...string) (string, error) {
	return LookPathWithPredicate(dirs, names, func(_ string, _ os.FileInfo) bool { return true })
}

func LookPathAll(dirs []string, names ...string) []string {
	results := make([]string, 0)
	LookPathWithPredicate(dirs, names, func(fpath string, _ os.FileInfo) bool {
		results = append(results, fpath)
		return false
	})
	return results
}
