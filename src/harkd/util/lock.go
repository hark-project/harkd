package util

import (
	"sync"
	"time"

	"github.com/nightlyone/lockfile"

	"harkd/errors"
)

const lockMaxRetries = 5
const lockDelayInterval = time.Millisecond * time.Duration(50)

// Lock implements a lock that has both a mutex (for exclusive access within a
// single process) and a file lock (for exclusive access between processes).
type Lock interface {
	Lock() error
	Unlock() error
}

// NewLock creates a new Lock whose file-level lock will use the provided lock
// file path.
func NewLock(lockFilePath string) (Lock, error) {
	return newLock(lockFilePath)
}

type lock struct {
	fileLock lockfile.Lockfile
	mutex    *sync.Mutex
}

func newLock(lockFilePath string) (*lock, error) {
	lockfile, err := lockfile.New(lockFilePath)
	if err != nil {
		// TODO(cera) - Wrap this
		return nil, err
	}

	return &lock{lockfile, new(sync.Mutex)}, nil
}

// Lock attempts to take the lock.
// It first takes the file lock, then the mutex.
func (l *lock) Lock() error {
	// Take the file lock, then the mutex
	for i := 0; ; i++ {
		err := l.fileLock.TryLock()
		if err == nil {
			// We now have the file lock
			break
		} else if i >= lockMaxRetries {
			return errors.ErrStateLock(err)
		}

		time.Sleep(lockDelayInterval)
	}

	l.mutex.Lock()

	return nil
}

func (l *lock) Unlock() error {
	// Release the file lock, then the mutex
	err := l.fileLock.Unlock()
	if err != nil {
		return err
	}
	l.mutex.Unlock()
	return nil
}
