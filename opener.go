package ose

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"
)

func newBufferedReader(rc io.ReadCloser) io.ReadCloser {
	br := bufio.NewReader(rc)
	return NewReadCloser(br, func(_ io.Reader) error {
		return rc.Close()
	})
}

func newBufferedWriter(wc io.WriteCloser) io.WriteCloser {
	bw := bufio.NewWriter(wc)
	return NewWriteCloser(bw, func(_ io.Writer) error {
		err := bw.Flush()
		if err != nil {
			return fmt.Errorf("flush : %w", err)
		}
		err2 := wc.Close()
		if err2 != nil {
			return fmt.Errorf("clise : %w", err)
		}
		return nil
	})
}

type Opener struct {
	FallbackReader        io.Reader
	FallbackWriter        io.Writer
	TreatHyphenAsFileName bool
	Unbuffered            bool
	fs                    afero.Fs
	io                    IO
}

func NewOpener(fs afero.Fs, io IO) *Opener {
	return &Opener{
		FallbackReader: io.In(),
		FallbackWriter: io.Out(),
		fs:             fs,
		io:             io,
	}
}

func NewOpenerInThisWorld() *Opener {
	return NewOpener(GetFs(), GetIO())
}

func (o *Opener) shouldFallback(name string) bool {
	return name == "" || (!o.TreatHyphenAsFileName && name == "-")
}

func (o *Opener) openFile(name string, ff func(name string) (afero.File, error)) (io.ReadCloser, error) {
	if o.shouldFallback(name) {
		if o.Unbuffered {
			return NopReadCloser(o.FallbackReader), nil
		}
		return NopReadCloser(bufio.NewReader(o.FallbackReader)), nil
	}
	f, err := ff(name)
	if err != nil {
		return nil, err
	}
	if o.Unbuffered {
		return f, nil
	}
	return newBufferedReader(f), nil

}

func (o *Opener) OpenFile(name string, flag int, perm os.FileMode) (io.ReadCloser, error) {
	return o.openFile(name, func(name string) (afero.File, error) { return o.fs.OpenFile(name, flag, perm) })
}

func (o *Opener) Open(name string) (io.ReadCloser, error) {
	return o.openFile(name, o.fs.Open)
}

func (o *Opener) createFile(name string, ff func(name string) (afero.File, error)) (io.WriteCloser, error) {
	if o.shouldFallback(name) {
		if o.Unbuffered {
			return NopWriteCloser(o.FallbackWriter), nil
		}
		bw := bufio.NewWriter(o.FallbackWriter)
		return NewWriteCloser(bw, func(_ io.Writer) error {
			return bw.Flush()
		}), nil
	}
	f, err := ff(name)
	if err != nil {
		return nil, err
	}
	if o.Unbuffered {
		return f, nil
	}
	return newBufferedWriter(f), nil
}

func (o *Opener) CreateFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return o.createFile(name, func(name string) (afero.File, error) { return o.fs.OpenFile(name, flag, perm) })
}

func (o *Opener) Create(name string) (io.WriteCloser, error) {
	return o.createFile(name, o.fs.Create)
}

func (o *Opener) TempScope() *TempScope {
	return NewTempScope(o.fs)
}

func (o *Opener) CreateTempFile(dir, prefix, newname string, handler func(f io.WriteCloser) (bool, error)) (bool, error) {
	if o.shouldFallback(newname) {
		wc, err := o.Create(newname)
		if err != nil {
			return false, err
		}
		defer wc.Close()
		return handler(wc)
	}
	return o.TempScope().TempFileScope(dir, prefix, newname, func(f afero.File) (bool, error) {
		var wc io.WriteCloser = f
		if !o.Unbuffered {
			wc = newBufferedWriter(NopWriteCloser(f))
			defer wc.Close()
		}
		return handler(wc)
	})
}
