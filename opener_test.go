package osplus

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestOpener(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-test-")
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
	o := NewOpener()
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
	tmp, err := ioutil.TempDir("", "osplus-test-")
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
	o := NewOpener()
	wc, _, err := o.CreateTempFileWithDestination(fooPath, "", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(wc, "foo")
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
	wc.Close()
	bs, err := ioutil.ReadFile(fooPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "foo" {
		t.Fatalf("invalid content: %v", bs)
	}
}

func TestCancelOpenerViaTempFile(t *testing.T) {
	tmp, err := ioutil.TempDir("", "osplus-test-")
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
	o := NewOpener()
	wc, commit, err := o.CreateTempFileWithDestination(fooPath, "", "osplus-test-")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(wc, "foo")
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
	commit(false)
	wc.Close()
	_, err = os.Stat(fooPath)
	if err == nil {
		t.Fatalf("%s must not be exist", fooPath)
	}
}
