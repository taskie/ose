package osplus_test

import (
	"testing"

	"github.com/taskie/ose/osplus"
)

func testTryToLockFile(t *testing.T) {
	lockFile := osplus.NewLockFile("osplus-test-foo.lock")
	err := lockFile.TryToLock()
	if err != nil {
		t.Fatal(err)
	}
	err = lockFile.TryToLock()
	if err == nil {
		t.Fatal("TryToLock (twice): must fail")
	}
	defer func() {
		err := lockFile.TryToUnlock()
		if err != nil {
			t.Fatal(lockFile.Error())
		}
		err = lockFile.TryToUnlock()
		if err == nil {
			t.Fatal("TryToUnlock (twice): must fail")
		}
	}()
}

func testLockFile(t *testing.T) {
	lockFile := osplus.NewLockFile("osplus-test-foo.lock")
	lockFile.Lock()
	defer func() {
		lockFile.Unlock()
		if lockFile.Error() != nil {
			t.Fatal(lockFile.Error())
		}
	}()
}
