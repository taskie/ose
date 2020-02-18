package osplus_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/taskie/ose/osplus"
)

func testRenameUsingLink(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ose-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}()

	fooPath := filepath.Join(tmp, "foo")
	barPath := filepath.Join(tmp, "bar")
	bazPath := filepath.Join(tmp, "baz")
	err = osplus.Touch(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	err = osplus.Touch(bazPath)
	if err != nil {
		t.Fatal(err)
	}
	err = osplus.RenameUsingLink(fooPath, barPath)
	if err != nil {
		t.Fatal(err)
	}
	err = osplus.RenameUsingLink(barPath, bazPath)
	if err == nil || !strings.Contains(err.Error(), "file exists") {
		t.Fatalf("RenameUsingLink: error 'exists' must occur, but %v", err)
	}
}

func testCopy(t *testing.T) {
	tmp, err := ioutil.TempDir("", "ose-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Fatal(err)
		}
	}()

	fooPath := filepath.Join(tmp, "foo")
	barPath := filepath.Join(tmp, "bar")
	bazPath := filepath.Join(tmp, "baz")

	err = ioutil.WriteFile(fooPath, []byte("hello, world!"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = osplus.Move(fooPath, barPath)
	if err != nil {
		t.Fatal(err)
	}
	err = osplus.Copy(barPath, bazPath)
	if err != nil {
		t.Fatal(err)
	}
	err = osplus.CopyFile(barPath, bazPath, &osplus.CopyOptions{
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Copy: overwrite must be inhibited")
	}
	err = osplus.MoveFile(barPath, bazPath, &osplus.MoveOptions{
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Move: overwrite must be inhibited")
	}
	err = osplus.MoveFile(barPath, bazPath, &osplus.MoveOptions{
		NoRename:    true,
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Move: overwrite must be inhibited")
	}
	err = osplus.MoveFile(barPath, bazPath, &osplus.MoveOptions{
		NoRename: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}
