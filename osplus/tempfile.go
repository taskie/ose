package osplus

import (
	"fmt"
	"io/ioutil"
	"os"
)

type TempFile struct {
	Destination string
	Dir         string
	Pattern     string
	MoveOptions *MoveOptions
	File        *os.File
}

func CreateTempFile(tempDir, tempPattern string) (*TempFile, error) {
	return CreateTempFileWithDestination("", tempDir, tempPattern)
}

func CreateTempFileWithDestination(dst, tempDir, tempPattern string) (*TempFile, error) {
	tmp := &TempFile{
		Destination: dst,
		Dir:         tempDir,
		Pattern:     tempPattern,
	}
	err := tmp.create()
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func (tmp *TempFile) create() error {
	f, err := ioutil.TempFile(tmp.Dir, tmp.Pattern)
	if err != nil {
		return err
	}
	tmp.File = f
	return nil
}

func (tmp *TempFile) Name() string {
	if tmp.File != nil {
		return tmp.File.Name()
	}
	return ""
}

func (tmp *TempFile) CloseFile() error {
	if tmp.File != nil {
		return tmp.File.Close()
	}
	return fmt.Errorf("tempfile is not opened yet")
}

func (tmp *TempFile) Close() error {
	if tmp.File != nil {
		err1 := tmp.CloseFile()
		if tmp.Destination == "" {
			// ignore err1
			err2 := tmp.remove()
			if err1 != nil {
				return err1
			}
			return err2
		}
		if err1 != nil {
			return err1
		}
		return tmp.move(tmp.Destination)
	}
	return fmt.Errorf("tempfile is not opened yet")
}

func (tmp *TempFile) Write(p []byte) (int, error) {
	if tmp.File != nil {
		return tmp.File.Write(p)
	}
	return 0, fmt.Errorf("tempfile is not opened yet")
}

func (tmp *TempFile) move(dst string) error {
	if tmp.File != nil {
		return MoveFile(tmp.Name(), dst, tmp.MoveOptions)
	}
	return fmt.Errorf("tempfile is not opened yet")
}

func (tmp *TempFile) remove() error {
	if tmp.File != nil {
		return os.Remove(tmp.Name())
	}
	return fmt.Errorf("tempfile is not opened yet")
}
