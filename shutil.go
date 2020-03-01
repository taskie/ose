package ose

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

func Exists(fs afero.Fs, filePath string) bool {
	if _, err := fs.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

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
			if Exists(fs, newname) {
				return fmt.Errorf("file exists: %s", newname)
			}
			err := fs.Rename(oldname, newname)
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

type CopyTreeOptions struct {
	NoOverwrite bool
}

func CopyTree(fs afero.Fs, oldname, newname string, opts *CopyTreeOptions) error {
	if opts == nil {
		opts = &CopyTreeOptions{}
	}
	oldFi, err := fs.Stat(oldname)
	if err != nil {
		return err
	}
	if !oldFi.IsDir() {
		return fmt.Errorf("not directory: %s", oldname)
	}
	if Exists(fs, newname) {
		if opts.NoOverwrite {
			return fmt.Errorf("already exists: %s", newname)
		}
	} else {
		err = fs.MkdirAll(newname, 0755)
		if err != nil {
			return err
		}
	}
	return copyTreeContent(fs, oldname, newname, opts, 1)
}

// https://stackoverflow.com/questions/51779243/copy-a-folder-in-go

func copyTreeContent(fs afero.Fs, oldname, newname string, opts *CopyTreeOptions, depth int) error {
	if depth > 127 {
		return fmt.Errorf("max depth exceeded")
	}
	entries, err := afero.ReadDir(fs, oldname)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		src := filepath.Join(oldname, entry.Name())
		dst := filepath.Join(newname, entry.Name())

		srcFI, err := fs.Stat(src)
		if err != nil {
			return err
		}

		_, ok := srcFI.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", src)
		}

		switch srcFI.Mode() & os.ModeType {
		case os.ModeDir:
			if !Exists(fs, dst) {
				if err := fs.MkdirAll(dst, 0755); err != nil {
					return err
				}
			}
			if err := copyTreeContent(fs, src, dst, opts, depth+1); err != nil {
				return err
			}
		case os.ModeSymlink:
			return fmt.Errorf("unimplemented: copy symlink")
		default:
			if err := Copy(fs, src, dst); err != nil {
				return err
			}
		}

		// if err := os.Lchown(dst, int(stat.Uid), int(stat.Gid)); err != nil {
		//   return err
		// }

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := fs.Chmod(dst, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

type MoveTreeOptions struct {
	NoOverwrite bool
	NoRename    bool
}

func MoveTree(fs afero.Fs, oldname, newname string, opts *MoveTreeOptions) error {
	if opts == nil {
		opts = &MoveTreeOptions{}
	}
	if !opts.NoRename {
		if opts.NoOverwrite {
			if Exists(fs, newname) {
				return fmt.Errorf("already exists: %s", newname)
			}
			err := fs.Rename(oldname, newname)
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
	copyOpts := &CopyTreeOptions{
		NoOverwrite: opts.NoOverwrite,
	}
	err := CopyTree(fs, oldname, newname, copyOpts)
	if err != nil {
		return err
	}
	return fs.RemoveAll(oldname)
}
