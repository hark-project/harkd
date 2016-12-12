package fs

import (
	"io"
	"io/ioutil"
	"os"
)

// Filesystem is an interface that wraps some standard file i/o operations.
//
// Other types can embed a Filesystem instead of directly making calls to OS:
// this makes them more testable.
type Filesystem interface {
	Open(string) (io.ReadCloser, error)
	Stat(string) (os.FileInfo, error)

	WriteFile(string, []byte, os.FileMode) error
}

// NewFilesystem constructs a new Filesystem backed by the real system filesystem.
func NewFilesystem() Filesystem {
	return fileSystem{}
}

type fileSystem struct {
}

func (fs fileSystem) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (fs fileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (fs fileSystem) WriteFile(path string, d []byte, perm os.FileMode) error {
	return ioutil.WriteFile(path, d, perm)
}
