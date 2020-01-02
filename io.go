package osplus

import (
	"bytes"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
)

type IO interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

type IOContainer struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func (c *IOContainer) In() io.Reader  { return c.InR }
func (c *IOContainer) Out() io.Writer { return c.OutW }
func (c *IOContainer) Err() io.Writer { return c.ErrW }

func NewIOContainer(in io.Reader, out, err io.Writer) *IOContainer {
	return &IOContainer{
		InR:  in,
		OutW: out,
		ErrW: err,
	}
}

func Stdio() IO {
	io := NewIOContainer(os.Stdin, os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		io.OutW = colorable.NewColorableStdout()
		io.ErrW = colorable.NewColorableStderr()
	}
	return io
}

type FakeIO struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (i *FakeIO) In() io.Reader  { return i.InBuf }
func (i *FakeIO) Out() io.Writer { return i.OutBuf }
func (i *FakeIO) Err() io.Writer { return i.ErrBuf }

func NewFakeIO() *FakeIO {
	return &FakeIO{
		InBuf:  new(bytes.Buffer),
		OutBuf: new(bytes.Buffer),
		ErrBuf: new(bytes.Buffer),
	}
}
