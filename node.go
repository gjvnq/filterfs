package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

const DEFAULT_PERMISSIONS = 0660

type FNode struct {
	inode     *nodefs.Inode
	fd        int
	RealPath  string
	Name      string
	Size      uint64
	Atime     uint64
	Mtime     uint64
	Ctime     uint64
	Atimensec uint32
	Mtimensec uint32
	Ctimensec uint32
}

func (fs *FNode) OnUnmount() {
	Log.DebugF("OnUnmount")
}

func (fs *FNode) OnMount(conn *nodefs.FileSystemConnector) {
	Log.DebugF("OnMount")
}

func IsHidden(path string) bool {
	l := strings.Split(path, "/")
	for _, elem := range l {
		for _, rule := range HideList {
			match, err := filepath.Match(rule, elem)
			if match && err == nil {
				Log.DebugF("IsHidden: Node '%s' is hidden beacuse '%s' matches '%s'", path, elem, rule)
				return true
			}
		}
	}
	return false
}

func (n *FNode) IsHidden() bool {
	return IsHidden(n.RealPath)
}

func (n *FNode) StatFs() *fuse.StatfsOut {
	Log.DebugF("StatFs")
	if n.IsHidden() {
		return nil
	}
	stat := syscall.Statfs_t{}
	if err := syscall.Statfs(n.RealPath, &stat); err != nil {
		Log.ErrorF("syscall.Statfs(%s, stat): %+v", n.RealPath, err)
		return nil
	}
	fuse_stat := new(fuse.StatfsOut)
	fuse_stat.FromStatfsT(&stat)
	return fuse_stat
}

func (n *FNode) SetInode(node *nodefs.Inode) {
	Log.DebugF("SetInode (%+v)", *node)
	n.inode = node
}

func (n *FNode) Deletable() bool {
	return n.fd != 0
}

func (n *FNode) Inode() *nodefs.Inode {
	Log.DebugF("Inode (n=%v)", *n)
	return n.inode
}

func (n *FNode) OnForget() {
	Log.DebugF("OnForget")
	if n.fd != 0 {
		syscall.Close(n.fd)
	}
}

func (n *FNode) Lookup(out *fuse.Attr, name string, context *fuse.Context) (ret_node *nodefs.Inode, ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Lookup", &_start)
	Log.DebugF("Lookup (n=%v; out=%v; name=%v; context=%v)", *n, *out, name, *context)

	// Set basics
	new_node := &FNode{}
	new_node.RealPath = n.RealPath + "/" + name
	// Is hidden?
	if new_node.IsHidden() {
		return nil, fuse.ENOENT
	}

	// Is dir?
	info, err := os.Stat(new_node.RealPath)
	if err != nil {
		// Exists?
		if os.IsNotExist(err) {
			return nil, fuse.ENOENT
		}

		Log.ErrorF("os.Stat(%s): %+v", new_node.RealPath, err)
		return nil, fuse.EIO
	}

	// Finish
	child := n.Inode().NewChild(name, info.IsDir(), new_node)
	child.Node().GetAttr(out, nil, context)
	return child, fuse.OK
}

func (n *FNode) Access(mode uint32, context *fuse.Context) (code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Access", &_start)
	Log.DebugF("Access (mode=%d)", mode)

	if n.IsHidden() {
		return fuse.ENOENT
	}
	err := syscall.Access(n.RealPath, mode)
	if err != nil {
		Log.ErrorF("syscall.Access(%s, %d) = %+v", n.RealPath, mode, err)
	}
	return fuse.ToStatus(err)
}

func (n *FNode) Readlink(c *fuse.Context) ([]byte, fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Readlink", &_start)

	Log.DebugF("Readlink")
	buf := make([]byte, 999)

	if n.IsHidden() {
		return nil, fuse.ENOENT
	}

	num, err := syscall.Readlink(n.RealPath, buf)
	if err != nil {
		Log.ErrorF("syscall.Readlink(%s, buf): %+v", n.RealPath, err)
		return nil, fuse.EIO
	}
	buf = buf[:num]
	return buf, fuse.OK
}

func (n *FNode) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Mknod", &_start)

	Log.DebugF("Mknod (name=%s, mode=%d, dev=%d) n.RealPath=%s", name, mode, dev, n.RealPath)
	if n.IsHidden() || IsHidden(name) {
		return nil, fuse.ENOENT
	}

	new_node := FNode{}
	new_node.RealPath = n.RealPath + "/" + name
	is_dir := (mode & uint32(os.ModeDir)) != 0
	new_node.SetInode(n.Inode().NewChild(name, is_dir, &new_node))
	return new_node.Inode(), fuse.ToStatus(syscall.Mknod(new_node.RealPath, mode, int(dev)))
}

func (n *FNode) Mkdir(name string, mode uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Mkdir", &_start)

	Log.DebugF("Mkdir (name=%s, mode=%d) n.RealPath=%s", name, mode, n.RealPath)
	if n.IsHidden() || IsHidden(name) {
		return nil, fuse.ENOENT
	}

	new_node := FNode{}
	new_node.RealPath = n.RealPath + "/" + name
	is_dir := true
	new_node.SetInode(n.Inode().NewChild(name, is_dir, &new_node))
	return new_node.Inode(), fuse.ToStatus(syscall.Mkdir(new_node.RealPath, mode))
}

func (n *FNode) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Unlink")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}
func (n *FNode) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Rmdir")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}
func (n *FNode) Symlink(name string, content string, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Symlink")
	if n.IsHidden() {
		return nil, fuse.ENOENT
	}
	return nil, fuse.ENOSYS
}

func (n *FNode) Rename(oldName string, newParent nodefs.Node, newName string, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Rename")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) Link(name string, existing nodefs.Node, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Link")
	if n.IsHidden() {
		return nil, fuse.ENOENT
	}
	return nil, fuse.ENOSYS
}

func (n *FNode) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Create")
	if n.IsHidden() {
		return nil, nil, fuse.ENOENT
	}
	return nil, nil, fuse.ENOSYS
}

func (n *FNode) Open(flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Open", &_start)
	Log.DebugF("Open (n.RealPath = %s, flags = %d)", n.RealPath, flags)
	if n.IsHidden() {
		return nil, fuse.ENOENT
	}

	var err error
	n.fd, err = syscall.Open(n.RealPath, int(flags), DEFAULT_PERMISSIONS)
	if err != nil {
		Log.ErrorF("syscall.Open(%s, %d, %d) = %d, %v", n.RealPath, flags, DEFAULT_PERMISSIONS, n.fd, err)
	}
	return nodefs.NewDefaultFile(), fuse.ToStatus(err)
}

func (n *FNode) Flush(file nodefs.File, openFlags uint32, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Flush")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) OpenDir(context *fuse.Context) (ret_dirs []fuse.DirEntry, ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("OpenDir", &_start)

	if n.IsHidden() {
		return nil, fuse.ENOENT
	}

	ans := make([]fuse.DirEntry, 0)
	files, err := ioutil.ReadDir(n.RealPath)
	if err != nil {
		Log.ErrorF("OpenDir (RealPath=%s): %+v", n.RealPath, err)
		return nil, fuse.ENOSYS
	}
	for _, file := range files {
		new_file := fuse.DirEntry{}
		new_file.Name = file.Name()
		new_file.Mode = uint32(file.Mode())
		if !IsHidden(file.Name()) {
			ans = append(ans, new_file)
		}
	}

	Log.DebugF("OpenDir (context=%v)", *context)
	return ans, fuse.OK
}

func (n *FNode) GetXAttr(attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("GetXAttr", &_start)

	Log.DebugF("GetXAttr (attribute=%v)", attribute)
	if n.IsHidden() {
		return nil, fuse.ENOENT
	}
	return nil, fuse.ENOSYS
}

func (n *FNode) RemoveXAttr(attr string, context *fuse.Context) fuse.Status {
	_start := time.Now()
	defer PrintCallDuration("RemoveXAttr", &_start)

	Log.DebugF("RemoveXAttr")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) SetXAttr(attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	_start := time.Now()
	defer PrintCallDuration("SetXAttr", &_start)

	Log.DebugF("SetXAttr")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) ListXAttr(context *fuse.Context) (attrs []string, code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("ListXAttr", &_start)

	Log.DebugF("ListXAttr")
	if n.IsHidden() {
		return nil, fuse.ENOENT
	}
	return nil, fuse.ENOSYS
}

func (n *FNode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("GetAttr", &_start)

	Log.DebugF("GetAttr (n=%v; out=%v; file=%v; context=%v)", *n, *out, file, *context)

	if n.IsHidden() {
		return fuse.ENOENT
	}

	path := n.RealPath
	var stat syscall.Stat_t
	// Get real parameters
	if err := syscall.Lstat(path, &stat); err != nil {
		Log.ErrorF("syscall.Lstat(%s, stat): %+v", path, err)
		return fuse.ToStatus(err)
	}
	out.FromStat(&stat)
	return fuse.OK
}

func (n *FNode) Chmod(file nodefs.File, perms uint32, context *fuse.Context) (code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Chmod", &_start)

	Log.DebugF("Chmod")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) Chown(file nodefs.File, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Chown", &_start)

	Log.DebugF("Chown")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) Truncate(file nodefs.File, size uint64, context *fuse.Context) (code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Truncate", &_start)

	Log.DebugF("Truncate")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) Utimens(file nodefs.File, atime *time.Time, mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Utimens", &_start)

	Log.DebugF("Utimens")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) Fallocate(file nodefs.File, off uint64, size uint64, mode uint32, context *fuse.Context) (code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Fallocate", &_start)

	Log.DebugF("Fallocate")
	if n.IsHidden() {
		return fuse.ENOENT
	}
	return fuse.ENOSYS
}

func (n *FNode) Read(file nodefs.File, dest []byte, off int64, context *fuse.Context) (ret_res fuse.ReadResult, ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Read", &_start)
	if n.IsHidden() {
		return nil, fuse.ENOENT
	}

	// Safety
	if n.fd == 0 {
		Log.WarningF("File '%s' has not been opened yet", n.RealPath)
		return fuse.ReadResultData(nil), fuse.EBADF
	}

	// Read
	_, err := syscall.Pread(n.fd, dest, off)
	if err != nil {
		Log.ErrorF("syscall.Pread(%d, [%d]byte{...}, %d) = _, %v", n.fd, len(dest), off, err)
		return fuse.ReadResultData(nil), fuse.ToStatus(err)
	}
	return fuse.ReadResultData(dest), fuse.OK
}

func (n *FNode) Write(file nodefs.File, data []byte, off int64, context *fuse.Context) (written uint32, code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Write", &_start)
	Log.DebugF("Write")
	if n.IsHidden() {
		return 0, fuse.ENOENT
	}

	// Safety
	if n.fd == 0 {
		Log.WarningF("File '%s' has not been opened yet", n.RealPath)
		return 0, fuse.EBADF
	}
	// Read
	num, err := syscall.Pwrite(n.fd, data, off)
	if err != nil {
		Log.ErrorF("syscall.Pwrite(%d, [%d]byte{...}, %d) = %d, %v", n.fd, len(data), off, num, err)
		return uint32(num), fuse.ToStatus(err)
	}
	return uint32(num), fuse.OK
}
