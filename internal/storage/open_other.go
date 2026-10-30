//go:build !linux

package storage

import (
	"errors"
	"os"
)

// OpenNoSymlinks relies on Linux O_PATH directory pinning. Other platforms
// build this package only for shared types (srv-front), never to touch storage.
func OpenNoSymlinks(path string, flags int) (*os.File, error) {
	return nil, errors.ErrUnsupported
}
