// +build !windows

package contenthash

import (
	"os"
	"syscall"

	"github.com/stevvooe/continuity/sysx"
	"github.com/tonistiigi/fsutil"
)

func chmodWindowsTarEntry(perm os.FileMode) os.FileMode {
	return perm
}

func setUnixOpt(path string, fi os.FileInfo, stat *fsutil.Stat) error {
	s := fi.Sys().(*syscall.Stat_t)

	stat.Uid = s.Uid
	stat.Gid = s.Gid

	if !fi.IsDir() {
		if s.Mode&syscall.S_IFBLK != 0 ||
			s.Mode&syscall.S_IFCHR != 0 {
			stat.Devmajor = int64(major(uint64(s.Rdev)))
			stat.Devminor = int64(minor(uint64(s.Rdev)))
		}
	}

	attrs, err := sysx.LListxattr(path)
	if err != nil {
		return err
	}
	if len(attrs) > 0 {
		stat.Xattrs = map[string][]byte{}
		for _, attr := range attrs {
			v, err := sysx.LGetxattr(path, attr)
			if err == nil {
				stat.Xattrs[attr] = v
			}
		}
	}
	return nil
}

func major(device uint64) uint64 {
	return (device >> 8) & 0xfff
}

func minor(device uint64) uint64 {
	return (device & 0xff) | ((device >> 12) & 0xfff00)
}
