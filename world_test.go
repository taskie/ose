package ose_test

import (
	"testing"

	"github.com/taskie/ose"
)

func TestNewRealWorld(t *testing.T) {
	_ = ose.NewRealWorld()
}

func TestNewFakeWorld(t *testing.T) {
	var _ ose.World = ose.NewFakeWorld()
}
