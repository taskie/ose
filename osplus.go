package osplus

import (
	"io"
	"os"
	"strings"
	"time"
)

var (
	Version = "0.1.0-beta"
	Revision = ""
)

func RenameUsingLink(oldpath string, newpath string) error {
	err := os.Link(oldpath, newpath)
	if err != nil {
		return err
	}
	err = os.Remove(oldpath)
	return err
}

type CopyOptions struct {
	NoOverwrite bool
}

func CopyWithOptions(oldpath string, newpath string, opts *CopyOptions) error {
	oldFile, err := os.Open(oldpath)
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
	newFile, err := os.OpenFile(newpath, flag, oldInfo.Mode())
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, oldFile)
	return err
}

func Copy(oldpath string, newpath string) error {
	return CopyWithOptions(oldpath, newpath, &CopyOptions{})
}

type MoveOptions struct {
	NoOverwrite bool
	NoRename    bool
}

func MoveWithOptions(oldpath string, newpath string, opts *MoveOptions) error {
	if !opts.NoRename {
		if opts.NoOverwrite {
			err := RenameUsingLink(oldpath, newpath)
			if err == nil {
				return nil
			} else if v, ok := err.(*os.LinkError); ok {
				if strings.Contains(v.Error(), "file exists") {
					return err
				}
			}
		} else {
			err := os.Rename(oldpath, newpath)
			if err == nil {
				return nil
			}
		}
	}
	err := CopyWithOptions(oldpath, newpath, &CopyOptions{
		NoOverwrite: opts.NoOverwrite,
	})
	if err != nil {
		return err
	}
	err = os.Remove(oldpath)
	return err
}

func Move(oldpath string, newpath string) error {
	return MoveWithOptions(oldpath, newpath, &MoveOptions{})
}

func Touch(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Close()
	now := time.Now()
	err = os.Chtimes(path, now, now)
	return err
}
