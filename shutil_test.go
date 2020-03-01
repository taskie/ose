package ose_test

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/taskie/ose"
)

func testCopy(t *testing.T) {
	tmp := "."
	fs := afero.NewMemMapFs()
	fooPath := filepath.Join(tmp, "foo")
	barPath := filepath.Join(tmp, "bar")
	bazPath := filepath.Join(tmp, "baz")

	err := afero.WriteFile(fs, fooPath, []byte("hello, world!"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ose.Move(fs, fooPath, barPath)
	if err != nil {
		t.Fatal(err)
	}
	err = ose.Copy(fs, barPath, bazPath)
	if err != nil {
		t.Fatal(err)
	}
	err = ose.CopyFile(fs, barPath, bazPath, &ose.CopyOptions{
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Copy: overwrite must be inhibited")
	}
	err = ose.MoveFile(fs, barPath, bazPath, &ose.MoveOptions{
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Move: overwrite must be inhibited")
	}
	err = ose.MoveFile(fs, barPath, bazPath, &ose.MoveOptions{
		NoRename:    true,
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Move: overwrite must be inhibited")
	}
	err = ose.MoveFile(fs, barPath, bazPath, &ose.MoveOptions{
		NoRename: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}
