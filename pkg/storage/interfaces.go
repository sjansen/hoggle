package storage

import "io"

type Factory interface {
	New() (Container, error)
}

type Container interface {
	Download(oid string, dst io.WriterAt) error
	Upload(oid string, src io.ReadSeeker) error
}
