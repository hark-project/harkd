package fixtures

import (
	"io"
	"os"
)

// NewFsFixture creates a new, empty FsFixture.
func NewFsFixture() *FsFixture {
	return &FsFixture{}
}

// FsFixture iplements Filesystem, where everything can be mocked out.
type FsFixture struct {
	MockOpen struct {
		CalledWith    string
		WillReturn    io.ReadCloser
		WillReturnErr error
	}
	MockStat struct {
		CalledWith    string
		WillReturn    os.FileInfo
		WillReturnErr error
	}
	MockWriteFile struct {
		CalledWithPath string
		CalledWithData []byte
		CalledWithPerm os.FileMode
		WillReturnErr  error
	}
}

func (fs *FsFixture) Open(path string) (io.ReadCloser, error) {
	fs.MockOpen.CalledWith = path
	return fs.MockOpen.WillReturn, fs.MockOpen.WillReturnErr
}

func (fs *FsFixture) Stat(path string) (os.FileInfo, error) {
	fs.MockStat.CalledWith = path
	return fs.MockStat.WillReturn, fs.MockStat.WillReturnErr
}

func (fs *FsFixture) WriteFile(path string, d []byte, perm os.FileMode) error {
	fs.MockWriteFile.CalledWithPath = path
	fs.MockWriteFile.CalledWithData = d
	fs.MockWriteFile.CalledWithPerm = perm

	return fs.MockWriteFile.WillReturnErr
}
