package osplus

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenameUsingLink(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-")
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
	err = Touch(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	err = Touch(bazPath)
	if err != nil {
		t.Fatal(err)
	}
	err = RenameUsingLink(fooPath, barPath)
	if err != nil {
		t.Fatal(err)
	}
	err = RenameUsingLink(barPath, bazPath)
	if err == nil || !strings.Contains(err.Error(), "file exists") {
		t.Fatalf("RenameUsingLink: error 'exists' must occur, but %v", err)
	}
}

func TestCopy(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-")
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
	err = Move(fooPath, barPath)
	if err != nil {
		t.Fatal(err)
	}
	err = Copy(barPath, bazPath)
	if err != nil {
		t.Fatal(err)
	}
	err = CopyWithOptions(barPath, bazPath, &CopyOptions{
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Copy: overwrite must be inhibited")
	}
	err = MoveWithOptions(barPath, bazPath, &MoveOptions{
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Move: overwrite must be inhibited")
	}
	err = MoveWithOptions(barPath, bazPath, &MoveOptions{
		NoRename:    true,
		NoOverwrite: true,
	})
	if err == nil {
		t.Fatal("Move: overwrite must be inhibited")
	}
	err = MoveWithOptions(barPath, bazPath, &MoveOptions{
		NoRename: true,
	})
	if err != nil {
		t.Fatal(err)
	}
}
