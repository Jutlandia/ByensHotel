package filesystem

import (
	"io/fs"
	"net/http"
)

// File system errors.
var (
	ErrNotExist = fs.ErrNotExist
)

// FileSystem is a custom file system handler. It implements the http.FileSystem interface.
type FileSystem struct {
	FS http.FileSystem
}

// Open returns the file with the given name if it exists.
// A fs.ErrNotExist is returned if name is a directory.
func (fs FileSystem) Open(name string) (http.File, error) {
	f, err := fs.FS.Open(name)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, ErrNotExist
	}
	return f, nil
}
