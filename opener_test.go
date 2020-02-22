package ose_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"
	"github.com/taskie/ose"
)

func TestOpener(t *testing.T) {
	w := ose.NewFakeWorld()
	o := ose.NewOpener(w.Fs(), w.IO())
	fooPath := "foo"
	wc, err := o.Create(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(wc, "foo")
	wc.Close()

	rc, err := o.Open(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	rc.Close()
	if string(bs) != "foo" {
		t.Fatalf("invalid content: %v", bs)
	}
}

func TestOpenerViaTempFile(t *testing.T) {
	fooPath := "foo"
	w := ose.NewFakeWorld()
	o := ose.NewOpener(w.Fs(), w.IO())
	ok, err := o.CreateTempFile("", "opener", fooPath, func(f io.WriteCloser) (bool, error) {
		_, err := io.WriteString(f, "foo")
		return true, err
	})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("must be ok")
	}
	bs, err := afero.ReadFile(w.Fs(), fooPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "foo" {
		t.Fatalf("invalid content: %v", bs)
	}
}

func TestCancelOpenerViaTempFile(t *testing.T) {
	fooPath := "foo"
	w := ose.NewFakeWorld()
	o := ose.NewOpener(w.Fs(), w.IO())
	ok, err := o.CreateTempFile("", "opener", fooPath, func(f io.WriteCloser) (bool, error) {
		_, err := io.WriteString(f, "foo")
		return false, err
	})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("must be not ok")
	}
	_, err = w.Fs().Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
}
