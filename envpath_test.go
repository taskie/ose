package ose_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/taskie/ose"
)

func TestPath(t *testing.T) {
	p := ose.NewEnvPath(afero.NewMemMapFs(), ose.NewMapEnv())
	_, _ = p.GetGoPath()
	_ = p.GetPath()
}

func TestXdg(t *testing.T) {
	p := ose.NewEnvPath(afero.NewMemMapFs(), ose.NewMapEnv())
	_, _ = p.GetXdgConfigHome()
	_, _ = p.GetXdgConfigHome()
	_, _ = p.GetXdgCacheHome()
	_, _ = p.GetXdgDataHome()
	_, _ = p.GetXdgRuntimeDir()
	_ = p.GetXdgDataDirs()
	_ = p.GetXdgConfigDirs()
}
