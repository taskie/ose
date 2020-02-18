package ose

import (
	"github.com/spf13/afero"
)

type TempScope struct {
	fs afero.Fs
}

func NewTempScope(fs afero.Fs) *TempScope {
	return &TempScope{fs: fs}
}

func (s *TempScope) TempFileScope(dir, prefix, newname string, handler func(f afero.File) (bool, error)) (bool, error) {
	f, err := afero.TempFile(s.fs, dir, prefix)
	if err != nil {
		return false, err
	}
	oldname := f.Name()
	ok, err := handler(f)
	if !ok || err != nil {
		_ = f.Close()
		_ = s.fs.Remove(oldname)
		return ok, err
	}
	err = f.Close()
	if err != nil {
		_ = s.fs.Remove(oldname)
		return ok, err
	}
	err = s.fs.Rename(oldname, newname)
	return ok, err
}

func (s *TempScope) TempDirScope(dir, prefix, newname string, handler func(tempname string) (bool, error)) (bool, error) {
	oldname, err := afero.TempDir(s.fs, dir, prefix)
	if err != nil {
		return false, err
	}
	ok, err := handler(oldname)
	if !ok || err != nil {
		_ = s.fs.RemoveAll(oldname)
		return ok, err
	}
	err = s.fs.Rename(oldname, newname)
	return ok, err
}
