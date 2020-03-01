package ose

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/afero"
)

type CopyOptions struct {
	NoOverwrite bool
}

func CopyFile(fs afero.Fs, oldname string, newname string, opts *CopyOptions) error {
	if opts == nil {
		opts = &CopyOptions{}
	}
	oldFile, err := fs.Open(oldname)
	if err != nil {
		return err
	}
	defer oldFile.Close()
	oldInfo, err := oldFile.Stat()
	if err != nil {
		return err
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if opts.NoOverwrite {
		flag |= os.O_EXCL
	}
	newFile, err := fs.OpenFile(newname, flag, oldInfo.Mode())
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		return err
	}
	err = newFile.Sync()
	return err
}

func Copy(fs afero.Fs, oldname string, newname string) error {
	return CopyFile(fs, oldname, newname, nil)
}

type MoveOptions struct {
	NoOverwrite bool
	NoRename    bool
}

func MoveFile(fs afero.Fs, oldname string, newname string, opts *MoveOptions) error {
	if opts == nil {
		opts = &MoveOptions{}
	}
	if !opts.NoRename {
		if opts.NoOverwrite {
			_, err := fs.Stat(newname)
			if err == nil {
				return fmt.Errorf("file exists")
			}
			err = fs.Rename(oldname, newname)
			if err == nil {
				return nil
			}
			// fall through
		} else {
			err := fs.Rename(oldname, newname)
			if err == nil {
				return nil
			}
			// fall through
		}
	}
	err := CopyFile(fs, oldname, newname, &CopyOptions{
		NoOverwrite: opts.NoOverwrite,
	})
	if err != nil {
		return err
	}
	err = fs.Remove(oldname)
	return err
}

func Move(fs afero.Fs, oldname string, newname string) error {
	return MoveFile(fs, oldname, newname, nil)
}

func Touch(fs afero.Fs, path string) error {
	file, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Close()
	now := time.Now()
	err = fs.Chtimes(path, now, now)
	return err
}
