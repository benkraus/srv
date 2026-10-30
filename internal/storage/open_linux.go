package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// OpenNoSymlinks opens a clean absolute path, pinning each directory
// before opening its child. Unlike openat2, openat works with systemd's
// RestrictSUIDSGID filter; O_NOFOLLOW still rejects symlinks at every component.
func OpenNoSymlinks(path string, flags int) (*os.File, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, fmt.Errorf("storage path must be absolute and clean: %q", path)
	}
	if path == "/" {
		fd, err := unix.Open(path, flags|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err != nil {
			return nil, err
		}
		return os.NewFile(uintptr(fd), path), nil
	}
	fd, err := unix.Open("/", unix.O_PATH|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(path[1:], "/")
	for i, part := range parts {
		openFlags := unix.O_PATH | unix.O_DIRECTORY | unix.O_NOFOLLOW | unix.O_CLOEXEC
		if i == len(parts)-1 {
			openFlags = flags | unix.O_NOFOLLOW | unix.O_CLOEXEC
		}
		next, err := unix.Openat(fd, part, openFlags, 0)
		unix.Close(fd)
		if err != nil {
			return nil, err
		}
		fd = next
	}
	// O_PATH|O_NOFOLLOW opens the symlink inode itself instead of rejecting it.
	if flags&unix.O_PATH != 0 {
		var st unix.Stat_t
		if err := unix.Fstat(fd, &st); err != nil {
			unix.Close(fd)
			return nil, err
		}
		if st.Mode&unix.S_IFMT == unix.S_IFLNK {
			unix.Close(fd)
			return nil, unix.ELOOP
		}
	}
	return os.NewFile(uintptr(fd), path), nil
}
