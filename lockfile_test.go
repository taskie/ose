package osplus

import (
	"testing"
)

func TestTryToLockFile(t *testing.T) {
	lockFile := NewLockFile("foo.lock")
	err := lockFile.TryToLock()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := lockFile.TryToUnlock()
		if err != nil {
			t.Fatal(lockFile.Error())
		}
	}()
}

func TestLockFile(t *testing.T) {
	lockFile := NewLockFile("foo.lock")
	lockFile.Lock()
	defer func() {
		lockFile.Unlock()
		if lockFile.Error() != nil {
			t.Fatal(lockFile.Error())
		}
	}()
}
