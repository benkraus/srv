//go:build !linux

package coldstore

import "errors"

// renameNoReplace relies on Linux renameat2. Other platforms build this
// package only for shared types (srv-front), never to move cold-store objects.
func renameNoReplace(oldpath, newpath string) error {
	return errors.ErrUnsupported
}
