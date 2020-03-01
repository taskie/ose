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
	newname, err := s.TempFileScopeLazy(dir, prefix, func(f afero.File) (string, error) {
		ok, err := handler(f)
		if !ok {
			return "", err
		}
		return newname, err
	})
	return newname != "", err
}

func (s *TempScope) TempFileScopeLazy(dir, prefix string, handler func(f afero.File) (string, error)) (string, error) {
	f, err := afero.TempFile(s.fs, dir, prefix)
	if err != nil {
		return "", err
	}
	oldname := f.Name()
	newname, err := handler(f)
	if newname == "" || err != nil {
		_ = f.Close()
		_ = s.fs.Remove(oldname)
		return newname, err
	}
	err = f.Close()
	if err != nil {
		_ = s.fs.Remove(oldname)
		return newname, err
	}
	err = Move(s.fs, oldname, newname)
	return newname, err
}

func (s *TempScope) TempDirScope(dir, prefix, newname string, handler func(tempname string) (bool, error)) (bool, error) {
	newname, err := s.TempDirScopeLazy(dir, prefix, func(tempname string) (string, error) {
		ok, err := handler(tempname)
		if !ok {
			return "", err
		}
		return newname, err
	})
	return newname != "", err
}

func (s *TempScope) TempDirScopeLazy(dir, prefix string, handler func(tempname string) (string, error)) (string, error) {
	oldname, err := afero.TempDir(s.fs, dir, prefix)
	if err != nil {
		return "", err
	}
	newname, err := handler(oldname)
	if newname == "" || err != nil {
		_ = s.fs.RemoveAll(oldname)
		return newname, err
	}
	// TODO: implement cp -r
	err = s.fs.Rename(oldname, newname)
	return newname, err
}
