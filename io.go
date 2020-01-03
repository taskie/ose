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
	if runtime.GOOS == "windows" {
		outW := colorable.NewColorableStdout()
		errW := colorable.NewColorableStderr()
		return NewIOContainer(os.Stdin, outW, errW)
	}
	return NewIOContainer(os.Stdin, os.Stdout, os.Stderr)
}

type BufIOContainer struct {
	InBuf  *bytes.Buffer
	OutBuf *bytes.Buffer
	ErrBuf *bytes.Buffer
}

func (i *BufIOContainer) In() io.Reader  { return i.InBuf }
func (i *BufIOContainer) Out() io.Writer { return i.OutBuf }
func (i *BufIOContainer) Err() io.Writer { return i.ErrBuf }

func NewBufIOContainer() *BufIOContainer {
	return &BufIOContainer{
		InBuf:  new(bytes.Buffer),
		OutBuf: new(bytes.Buffer),
		ErrBuf: new(bytes.Buffer),
	}
}
