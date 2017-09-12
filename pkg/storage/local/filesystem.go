package local

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sjansen/hoggle/pkg/storage"
)

type Filesystem struct{}

func (fs *Filesystem) Open(name string) (io.ReadSeeker, error) {
	return os.Open(name)
}

func (fs *Filesystem) TempFile() (storage.TempFile, error) {
	return ioutil.TempFile("", "hoggle-")
}
