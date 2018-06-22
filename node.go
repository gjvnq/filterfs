package main

import (
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type FNode struct {
	inode     *nodefs.Inode
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

func (n *FNode) StatFs() *fuse.StatfsOut {
	Log.DebugF("StatFs")
	return nil
}

func (n *FNode) SetInode(node *nodefs.Inode) {
	Log.DebugF("SetInode (%+v)", *node)
	n.inode = node
}

func (n *FNode) Deletable() bool {
	Log.DebugF("Deletable")
	return true
}

func (n *FNode) Inode() *nodefs.Inode {
	Log.DebugF("Inode (n=%v)", *n)
	return n.inode
}

func (n *FNode) OnForget() {
	Log.DebugF("OnForget")
}

func (n *FNode) Lookup(out *fuse.Attr, name string, context *fuse.Context) (ret_node *nodefs.Inode, ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Lookup", &_start)

	Log.DebugF("Lookup (n=%v; out=%v; name=%v; context=%v)", *n, *out, name, *context)
	return nil, fuse.ENOSYS
}

func (n *FNode) Access(mode uint32, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Access")
	return fuse.ENOSYS
}

func (n *FNode) Readlink(c *fuse.Context) ([]byte, fuse.Status) {
	Log.DebugF("Readlink")
	return nil, fuse.ENOSYS
}

func (n *FNode) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Mknod")
	return nil, fuse.ENOSYS
}
func (n *FNode) Mkdir(name string, mode uint32, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Mkdir")
	return nil, fuse.ENOSYS
}
func (n *FNode) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Unlink")
	return fuse.ENOSYS
}
func (n *FNode) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Rmdir")
	return fuse.ENOSYS
}
func (n *FNode) Symlink(name string, content string, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Symlink")
	return nil, fuse.ENOSYS
}

func (n *FNode) Rename(oldName string, newParent nodefs.Node, newName string, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Rename")
	return fuse.ENOSYS
}

func (n *FNode) Link(name string, existing nodefs.Node, context *fuse.Context) (newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Link")
	return nil, fuse.ENOSYS
}

func (n *FNode) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, newNode *nodefs.Inode, code fuse.Status) {
	Log.DebugF("Create")
	return nil, nil, fuse.ENOSYS
}

func (n *FNode) Open(flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	Log.DebugF("Open")
	return nil, fuse.OK
}

func (n *FNode) Flush(file nodefs.File, openFlags uint32, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Flush")
	return fuse.ENOSYS
}

func (n *FNode) OpenDir(context *fuse.Context) (ret_dirs []fuse.DirEntry, ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("OpenDir", &_start)

	Log.DebugF("OpenDir (context=%v)", *context)
	return nil, fuse.ENOSYS
}

func (n *FNode) GetXAttr(attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	Log.DebugF("GetXAttr (attribute=%v)", attribute)
	return nil, fuse.ENOSYS
}

func (n *FNode) RemoveXAttr(attr string, context *fuse.Context) fuse.Status {
	Log.DebugF("RemoveXAttr")
	return fuse.ENOSYS
}

func (n *FNode) SetXAttr(attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	Log.DebugF("SetXAttr")
	return fuse.ENOSYS
}

func (n *FNode) ListXAttr(context *fuse.Context) (attrs []string, code fuse.Status) {
	Log.DebugF("ListXAttr")
	return nil, fuse.ENOSYS
}

func (n *FNode) GetAttr(out *fuse.Attr, file nodefs.File, context *fuse.Context) (ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("GetAttr", &_start)

	Log.DebugF("GetAttr (n=%v; out=%v; file=%v; context=%v)", *n, *out, file, *context)
	return fuse.ENOSYS
}

func (n *FNode) Chmod(file nodefs.File, perms uint32, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Chmod")
	return fuse.ENOSYS
}

func (n *FNode) Chown(file nodefs.File, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Chown")
	return fuse.ENOSYS
}

func (n *FNode) Truncate(file nodefs.File, size uint64, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Truncate")
	return fuse.ENOSYS
}

func (n *FNode) Utimens(file nodefs.File, atime *time.Time, mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Utimens")
	return fuse.ENOSYS
}

func (n *FNode) Fallocate(file nodefs.File, off uint64, size uint64, mode uint32, context *fuse.Context) (code fuse.Status) {
	Log.DebugF("Fallocate")
	return fuse.ENOSYS
}

func (n *FNode) Read(file nodefs.File, dest []byte, off int64, context *fuse.Context) (ret_res fuse.ReadResult, ret_code fuse.Status) {
	_start := time.Now()
	defer PrintCallDuration("Read", &_start)
	return nil, fuse.ENOSYS
}

func (n *FNode) Write(file nodefs.File, data []byte, off int64, context *fuse.Context) (written uint32, code fuse.Status) {
	Log.DebugF("Write")
	return 0, fuse.ENOSYS
}
