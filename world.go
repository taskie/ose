package ose

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/spf13/afero"
)

var world World

func init() {
	SetWorld(NewRealWorld())
}

func GetWorld() World  { return world }
func SetWorld(w World) { world = w }

type IO interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

type realIO struct{}

func (realIO) In() io.Reader  { return os.Stdin }
func (realIO) Out() io.Writer { return os.Stdout }
func (realIO) Err() io.Writer { return os.Stderr }

type IOContainer struct {
	InR  io.Reader
	OutW io.Writer
	ErrW io.Writer
}

func NewIOContainer(in io.Reader, out, err io.Writer) *IOContainer {
	return &IOContainer{
		InR:  in,
		OutW: out,
		ErrW: err,
	}
}

func (c *IOContainer) In() io.Reader  { return c.InR }
func (c *IOContainer) Out() io.Writer { return c.OutW }
func (c *IOContainer) Err() io.Writer { return c.ErrW }

func NewStdio() IO {
	if runtime.GOOS == "windows" {
		outW := colorable.NewColorableStdout()
		errW := colorable.NewColorableStderr()
		return NewIOContainer(os.Stdin, outW, errW)
	}
	return realIO{}
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

type Env interface {
	Get(key string) string
	Lookup(key string) (string, bool)
	Set(key string, value string) error
	Clear()
}

type realEnv struct{}

func (realEnv) Get(key string) string              { return os.Getenv(key) }
func (realEnv) Lookup(key string) (string, bool)   { return os.LookupEnv(key) }
func (realEnv) Set(key string, value string) error { return os.Setenv(key, value) }
func (realEnv) Clear()                             { os.Clearenv() }

type MapEnv struct {
	m map[string]string
}

func NewMapEnv() *MapEnv {
	return &MapEnv{m: make(map[string]string)}
}

func (e *MapEnv) Get(key string) string            { return e.m[key] }
func (e *MapEnv) Lookup(key string) (string, bool) { v, ok := e.m[key]; return v, ok }
func (e *MapEnv) Set(key string, value string) error {
	e.m[key] = value
	return nil
}
func (e *MapEnv) Clear() { e.m = make(map[string]string) }

func (e *MapEnv) GetMap() map[string]string  { return e.m }
func (e *MapEnv) SetMap(m map[string]string) { e.m = m }

// see https://stackoverflow.com/questions/18970265
type Clock interface {
	Now() time.Time
	After(d time.Duration) <-chan time.Time
}

type realClock struct{}

func (realClock) Now() time.Time                         { return time.Now() }
func (realClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

type FakeClock struct {
	Time     time.Time
	Duration time.Duration
	Count    int64
}

func NewFakeClock(t time.Time, d time.Duration) *FakeClock {
	return &FakeClock{Time: t, Duration: d}
}

func (c *FakeClock) Now() time.Time {
	t := c.Time.Add(time.Duration(c.Count) * c.Duration)
	c.Count++
	return t
}
func (c *FakeClock) After(d time.Duration) <-chan time.Time { return time.After(c.Duration) }

type World interface {
	Fs() afero.Fs
	IO() IO
	Env() Env
	Clock() Clock
}

type WorldContainer struct {
	fs    afero.Fs
	io    IO
	env   Env
	clock Clock
}

func NewWorldContainer(fs afero.Fs, io IO, env Env, clock Clock) *WorldContainer {
	return &WorldContainer{fs: fs, io: io, env: env, clock: clock}
}

func (w *WorldContainer) Fs() afero.Fs { return w.fs }
func (w *WorldContainer) IO() IO       { return w.io }
func (w *WorldContainer) Env() Env     { return w.env }
func (w *WorldContainer) Clock() Clock { return w.clock }

type realWorld struct {
	fs afero.Fs
	io IO
}

func NewRealWorld() World {
	return &realWorld{fs: afero.NewOsFs(), io: NewStdio()}
}

func (w *realWorld) Fs() afero.Fs { return w.fs }
func (w *realWorld) IO() IO       { return w.io }
func (w *realWorld) Env() Env     { return realEnv{} }
func (w *realWorld) Clock() Clock { return realClock{} }

type FakeWorld struct {
	FakeFs    afero.Fs
	FakeIO    *BufIOContainer
	FakeEnv   *MapEnv
	FakeClock *FakeClock
}

func NewFakeWorld() *FakeWorld {
	return &FakeWorld{
		FakeFs:    afero.NewMemMapFs(),
		FakeIO:    NewBufIOContainer(),
		FakeEnv:   NewMapEnv(),
		FakeClock: NewFakeClock(time.Now(), time.Millisecond),
	}
}

func (w *FakeWorld) Fs() afero.Fs { return w.FakeFs }
func (w *FakeWorld) IO() IO       { return w.FakeIO }
func (w *FakeWorld) Env() Env     { return w.FakeEnv }
func (w *FakeWorld) Clock() Clock { return w.FakeClock }

func GetFs() afero.Fs { return world.Fs() }
func GetIO() IO       { return world.IO() }
func GetEnv() Env     { return world.Env() }
func GetClock() Clock { return world.Clock() }
