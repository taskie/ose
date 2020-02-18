package ose_test

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"
	"github.com/taskie/ose"
)

func TestTempFileScope(t *testing.T) {
	fs := afero.NewMemMapFs()
	s := ose.NewTempScope(fs)
	ok, err := s.TempFileScope("", "bar", "baz", func(f afero.File) (bool, error) {
		t.Log(f.Name())
		_, err := f.WriteString("hello")
		return true, err
	})
	if err != nil {
		t.Fatalf("returns error: %v", err)
	}
	if !ok {
		t.Fatal("returns not ok")
	}
	f, err := fs.Open("baz")
	if err != nil {
		t.Fatalf("can't open file: %v", err)
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	t.Log(f)
	if err != nil {
		t.Fatalf("can't read file: %v", err)
	}
	actual := string(bs)
	if actual != "hello" {
		t.Fatalf("invalid content (actual): %v", actual)
	}
}

func TestTempDirScope(t *testing.T) {
	fs := afero.NewMemMapFs()
	s := ose.NewTempScope(fs)
	ok, err := s.TempDirScope("", "bar", "baz", func(tempname string) (bool, error) {
		return true, nil
	})
	if err != nil {
		t.Fatalf("returns error: %v", err)
	}
	if !ok {
		t.Fatal("returns not ok")
	}
	fi, err := fs.Stat("baz")
	if err != nil {
		t.Fatalf("invalid file state: %v", err)
	}
	if !fi.IsDir() {
		t.Fatalf("not directory")
	}
}
