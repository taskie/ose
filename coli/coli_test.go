package coli_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/cobra"
	"github.com/taskie/ose"
	"github.com/taskie/ose/coli"
)

func TestColi(t *testing.T) {
	w := ose.NewFakeWorld()
	ose.SetWorld(w)
	cl := coli.NewColiInThisWorld()
	var prod coli.ColiCommandProducer = func(cl *coli.Coli, name string, path []string) *cobra.Command {
		return &cobra.Command{
			Use: name,
			Run: func(cmd *cobra.Command, args []string) {
				actualIn, err := ioutil.ReadAll(cmd.InOrStdin())
				if err != nil {
					t.Fatalf("some error occured (in): %v", err)
				}
				if !bytes.Equal([]byte("in"), actualIn) {
					t.Fatalf("invalid content: %v", actualIn)
				}
				cmd.Print("out")
				cmd.PrintErr("err")
			},
		}
	}
	cmd := prod(cl, "test", []string{})
	cl.Prepare(cmd)
	_, _ = w.FakeIO.InBuf.WriteString("in")
	err := cl.Execute(cmd)
	if err != nil {
		t.Fatalf("some error occured (execute): %v", err)
	}
	actualOut := w.FakeIO.OutBuf.Bytes()
	if !bytes.Equal([]byte("out"), actualOut) {
		t.Fatalf("invalid content: %v", actualOut)
	}
	actualErr := w.FakeIO.ErrBuf.Bytes()
	if !bytes.Equal([]byte("err"), actualErr) {
		t.Fatalf("invalid content: %v", actualErr)
	}
}
