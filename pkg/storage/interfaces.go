package storage

import "io"

type Container interface {
	Download(oid string, dst io.WriterAt) error
	Upload(oid string, src io.ReadSeeker) error
}

type Factory interface {
	New() (Container, error)
}

type Filesystem interface {
	Open(name string) (io.ReadSeeker, error)
	TempFile() (TempFile, error)
}

type TempFile interface {
	io.Closer
	io.WriterAt
	Name() string
}
