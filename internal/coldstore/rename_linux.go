package coldstore

import "golang.org/x/sys/unix"

// renameNoReplace renames atomically and fails if newpath already exists.
func renameNoReplace(oldpath, newpath string) error {
	return unix.Renameat2(unix.AT_FDCWD, oldpath, unix.AT_FDCWD, newpath, unix.RENAME_NOREPLACE)
}
