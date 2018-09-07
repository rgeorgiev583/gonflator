package fuse

import (
	"hash/adler32"
	"log"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/rgeorgiev583/gonflator"
)

const (
	ReadOnly = 1 << iota
	RestrictPermissions
)

func getFuseErrorCode(err error) fuse.Status {
	gonflatorError, ok := err.(gonflator.Error)
	if !ok {
		return fuse.ENOSYS
	}

	switch gonflatorError.Code {
	case gonflator.CouldNotInitialize:
		return fuse.EBUSY

	case gonflator.NotFound:
		return fuse.ENOENT

	case gonflator.NoError:
		return fuse.OK

	case gonflator.ENOMEM:
		return fuse.Status(syscall.ENOMEM)

	case gonflator.EINTERNAL:
		return fuse.EIO

	case gonflator.EPATHX:
		return fuse.EINVAL

	case gonflator.ENOMATCH:
		return fuse.ENOENT

	case gonflator.EMMATCH:
		return fuse.Status(syscall.EOVERFLOW)

	case gonflator.EMVDESC:
		return fuse.EINVAL

	case gonflator.EBADARG:
		return fuse.EINVAL
	}

	return fuse.EIO
}

type setValue struct{}

type ConfigurationServerOptions uint

type ConfigurationServer struct {
	pathfs.FileSystem
	Provider gonflator.ConfigurationProvider
	Options  ConfigurationServerOptions

	directoryCache map[string]setValue
}

func (server *ConfigurationServer) isDir(name string) (res bool, err error) {
	if _, ok := server.directoryCache[name]; ok {
		res = true
		return
	}

	isTree, err := server.Provider.IsTree(name)
	if err != nil {
		return
	}

	if isTree {
		res = true
		return
	}

	return
}

func (server *ConfigurationServer) checkForWritePermissions(context *fuse.Context) fuse.Status {
	if server.Options&ReadOnly != 0 {
		return fuse.EROFS
	}

	if (server.Options&RestrictPermissions != 0) && int(context.Uid) != syscall.Geteuid() {
		return fuse.EACCES
	}

	return fuse.OK
}

func (server *ConfigurationServer) checkThatParentExistsButNotSelf(name string) fuse.Status {
	_, err := server.Provider.GetSetting(name)
	if err == nil {
		return fuse.Status(syscall.EEXIST)
	}

	_, err = server.Provider.GetSetting(filepath.Dir(name))
	if err != nil {
		return getFuseErrorCode(err)
	}

	return fuse.OK
}

func (server *ConfigurationServer) getMode(isDir bool) (mode uint32) {
	mode = syscall.S_IRUSR
	if !(server.Options&ReadOnly != 0) {
		mode |= syscall.S_IWUSR
	}
	if !(server.Options&RestrictPermissions != 0) {
		mode |= syscall.S_IRGRP | syscall.S_IROTH

		if !(server.Options&ReadOnly != 0) {
			mode |= syscall.S_IWGRP | syscall.S_IWOTH
		}
	}
	if isDir {
		mode |= fuse.S_IFDIR
	} else {
		mode |= fuse.S_IFREG
	}
	return
}

func (server *ConfigurationServer) getAttr(name string, value string, isDir bool, owner *fuse.Owner) *fuse.Attr {
	nameChecksum := adler32.Checksum([]byte(name))
	valueLength := len(value)
	currentTime := time.Now()
	currentTimeSeconds := currentTime.Unix()
	currentTimeNanoseconds := currentTime.Nanosecond()
	mode := server.getMode(isDir)

	return &fuse.Attr{
		Ino:       uint64(nameChecksum),
		Size:      uint64(valueLength),
		Atime:     uint64(currentTimeSeconds),
		Mtime:     uint64(currentTimeSeconds),
		Ctime:     uint64(currentTimeSeconds),
		Atimensec: uint32(currentTimeNanoseconds),
		Mtimensec: uint32(currentTimeNanoseconds),
		Ctimensec: uint32(currentTimeNanoseconds),
		Mode:      mode,
		Owner:     *owner,
	}
}

func (server *ConfigurationServer) String() string {
	return "GonflationFS(" + server.Provider.Name() + ")"
}

func (server *ConfigurationServer) GetAttr(name string, context *fuse.Context) (attr *fuse.Attr, code fuse.Status) {
	isDir, err := server.isDir(name)
	if err != nil {
		code = getFuseErrorCode(err)
		return
	}

	var value string
	if !isDir {
		value, err = server.Provider.GetSetting(name)
		if err != nil {
			code = getFuseErrorCode(err)
			return
		}
	}

	attr = server.getAttr(name, value, isDir, &context.Owner)
	return
}

func (server *ConfigurationServer) Truncate(name string, size uint64, context *fuse.Context) fuse.Status {
	code := server.checkForWritePermissions(context)
	if code != fuse.OK {
		return code
	}

	isDir, err := server.isDir(name)
	if err != nil {
		return getFuseErrorCode(err)
	}
	if isDir {
		return fuse.Status(syscall.EISDIR)
	}

	value, err := server.Provider.GetSetting(name)
	if err != nil {
		return getFuseErrorCode(err)
	}

	err = server.Provider.SetSetting(name, value[:size])
	if err != nil {
		return getFuseErrorCode(err)
	}

	err = server.Provider.Save()
	if err != nil {
		return getFuseErrorCode(err)
	}

	return fuse.OK
}

func (server *ConfigurationServer) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	if (server.Options&ReadOnly != 0) && (mode&fuse.W_OK != 0) {
		return fuse.EROFS
	}

	_, err := server.Provider.GetSetting(name)
	if err != nil {
		return getFuseErrorCode(err)
	}

	if mode&fuse.X_OK != 0 {
		return fuse.EACCES
	}

	return fuse.OK
}

func (server *ConfigurationServer) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	code := server.checkForWritePermissions(context)
	if code != fuse.OK {
		return code
	}

	code = server.checkThatParentExistsButNotSelf(name)
	if code != fuse.OK {
		return code
	}

	err := server.Provider.Save()
	if err != nil {
		return getFuseErrorCode(err)
	}

	server.directoryCache[name] = setValue{}
	return fuse.OK
}

func (server *ConfigurationServer) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	code := server.checkForWritePermissions(context)
	if code != fuse.OK {
		return code
	}

	code = server.checkThatParentExistsButNotSelf(name)
	if code != fuse.OK {
		return code
	}

	if mode&fuse.S_IFREG != 0 {
		err := server.Provider.SetSetting(name, "")
		if err != nil {
			return getFuseErrorCode(err)
		}

		err = server.Provider.Save()
		if err != nil {
			return getFuseErrorCode(err)
		}

		return fuse.OK
	} else if (mode&syscall.S_IFBLK != 0) || (mode&syscall.S_IFCHR != 0) || (mode&fuse.S_IFIFO != 0) || (mode&syscall.S_IFSOCK != 0) {
		return fuse.EPERM
	}

	return fuse.EINVAL
}

func (server *ConfigurationServer) Rename(oldName string, newName string, context *fuse.Context) fuse.Status {
	code := server.checkForWritePermissions(context)
	if code != fuse.OK {
		return code
	}

	if _, ok := server.directoryCache[oldName]; ok {
		delete(server.directoryCache, oldName)
		server.directoryCache[newName] = setValue{}
	}

	_, err := server.Provider.GetSetting(oldName)
	if err != nil {
		return getFuseErrorCode(err)
	}

	_, err = server.Provider.GetSetting(filepath.Dir(newName))
	if err != nil {
		return getFuseErrorCode(err)
	}

	err = server.Provider.MoveTree(oldName, newName)
	if err != nil {
		return getFuseErrorCode(err)
	}

	if _, ok := server.directoryCache[oldName]; ok {
		delete(server.directoryCache, oldName)
		server.directoryCache[newName] = setValue{}
	}

	err = server.Provider.Save()
	if err != nil {
		return getFuseErrorCode(err)
	}

	return fuse.OK
}

func (server *ConfigurationServer) Rmdir(name string, context *fuse.Context) fuse.Status {
	code := server.checkForWritePermissions(context)
	if code != fuse.OK {
		return code
	}

	isDir, err := server.isDir(name)
	if err != nil {
		return getFuseErrorCode(err)
	}
	if !isDir {
		return fuse.ENOTDIR
	}

	isTree, err := server.Provider.IsTree(name)
	if err != nil {
		return getFuseErrorCode(err)
	}
	if isTree {
		return fuse.Status(syscall.ENOTEMPTY)
	}

	delete(server.directoryCache, name)
	return fuse.OK
}

func (server *ConfigurationServer) Unlink(name string, context *fuse.Context) fuse.Status {
	code := server.checkForWritePermissions(context)
	if code != fuse.OK {
		return code
	}

	isDir, err := server.isDir(name)
	if err != nil {
		return getFuseErrorCode(err)
	}
	if isDir {
		return fuse.Status(syscall.EISDIR)
	}

	_, err = server.Provider.GetSetting(name)
	if err != nil {
		return getFuseErrorCode(err)
	}

	err = server.Provider.ClearSetting(name)
	if err != nil {
		return getFuseErrorCode(err)
	}

	err = server.Provider.Save()
	if err != nil {
		return getFuseErrorCode(err)
	}

	return fuse.OK
}

func (server *ConfigurationServer) OnMount(nodeFs *pathfs.PathNodeFs) {
	err := server.Provider.Load()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (server *ConfigurationServer) OnUnmount() {
	err := server.Provider.Save()
	if err != nil {
		log.Println(err.Error())
	}
}

func (server *ConfigurationServer) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	value, err := server.Provider.GetSetting(name)
	if err != nil {
		code = getFuseErrorCode(err)
		return
	}

	file = &fuseConfigurationSetting{nodefs.NewDefaultFile(), server, context, flags, name, value}
	return
}

func (server *ConfigurationServer) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	value, err := server.Provider.GetSetting(name)
	if err != nil {
		code = getFuseErrorCode(err)
		return
	}

	file = &fuseConfigurationSetting{nodefs.NewDefaultFile(), server, context, flags, name, value}
	return
}

func (server *ConfigurationServer) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, code fuse.Status) {
	isDir, err := server.isDir(name)
	if err != nil {
		code = getFuseErrorCode(err)
		return
	}
	if !isDir {
		code = fuse.ENOTDIR
		return
	}

	entries, err := server.Provider.ListSettings(name)
	if err != nil {
		code = getFuseErrorCode(err)
		return
	}

	for _, entry := range entries {
		mode := server.getMode(isDir)
		baseName := path.Base(entry)
		nameChecksum := adler32.Checksum([]byte(name))

		dirEntry := fuse.DirEntry{
			Mode: mode,
			Name: baseName,
			Ino:  uint64(nameChecksum),
		}
		stream = append(stream, dirEntry)
	}
	return
}
