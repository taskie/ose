package osplus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type LockFile struct {
	Path              string
	FileMode          os.FileMode
	AcquiringInterval time.Duration
	id                string
	body              []byte
	mutex             sync.Mutex
	error             error
}

type LockFileIDType string

const (
	LockFileIDTypeNone LockFileIDType = ""
	LockFileIDTypePID                 = "PID"
)

func NewLockFile(path string) *LockFile {
	lockFile, err := NewLockFileWithAutoGeneratedId(path, LockFileIDTypePID)
	if err != nil {
		panic(err)
	}
	return lockFile
}

func NewLockFileWithId(path string, id string) *LockFile {
	body := []byte{}
	if id != "" {
		body = []byte(id + "\n")
	}
	return &LockFile{
		Path:     path,
		FileMode: os.FileMode(0644),
		id:       id,
		body:     body,
	}
}

func NewLockFileWithAutoGeneratedId(path string, idType LockFileIDType) (*LockFile, error) {
	var id string
	switch idType {
	case LockFileIDTypeNone:
		id = ""
	case LockFileIDTypePID:
		id = fmt.Sprintf("%d", os.Getpid())
	default:
		return &LockFile{}, fmt.Errorf("invalid auto generated id type: %s", idType)
	}
	return NewLockFileWithId(path, id), nil
}

func (lockFile *LockFile) ID() string {
	return lockFile.id
}

func (lockFile *LockFile) tryToLockImpl(noMutexLock bool) error {
	if !noMutexLock {
		lockFile.mutex.Lock()
		defer lockFile.mutex.Unlock()
	}
	mode := os.FileMode(0644)
	if lockFile.FileMode > 0 {
		mode = lockFile.FileMode
	}
	file, err := os.OpenFile(lockFile.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	defer file.Close()
	if len(lockFile.body) != 0 {
		_, err = file.Write(lockFile.body)
	}
	return err
}

func (lockFile *LockFile) TryToLock() error {
	return lockFile.tryToLockImpl(false)
}

func (lockFile *LockFile) Lock() {
	lockFile.mutex.Lock()
	defer lockFile.mutex.Unlock()
	for {
		err := lockFile.tryToLockImpl(true)
		lockFile.error = err
		if err == nil {
			break
		}
		time.Sleep(lockFile.AcquiringInterval)
	}
}

func (lockFile *LockFile) tryToUnlockImpl(noMutexLock bool) error {
	if !noMutexLock {
		lockFile.mutex.Lock()
		defer lockFile.mutex.Unlock()
	}
	if lockFile.body != nil {
		actualBody, err := ioutil.ReadFile(lockFile.Path)
		if err != nil {
			return err
		}
		if !bytes.Equal(actualBody, lockFile.body) {
			return fmt.Errorf("invalid lockfile content")
		}
	}
	err := os.Remove(lockFile.Path)
	return err
}

func (lockFile *LockFile) TryToUnlock() error {
	return lockFile.tryToUnlockImpl(false)
}

func (lockFile *LockFile) Unlock() {
	lockFile.mutex.Lock()
	defer lockFile.mutex.Unlock()
	err := lockFile.tryToUnlockImpl(true)
	lockFile.error = err
}

func (lockFile *LockFile) Error() error {
	return lockFile.error
}