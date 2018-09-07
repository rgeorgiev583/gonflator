package fuse

import (
	"syscall"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type fuseConfigurationSetting struct {
	nodefs.File

	server  *ConfigurationServer
	context *fuse.Context
	flags   uint32
	name    string
	value   string
}

func (setting *fuseConfigurationSetting) String() string {
	return setting.value
}

func (setting *fuseConfigurationSetting) Read(dest []byte, off int64) (res fuse.ReadResult, code fuse.Status) {
	if setting.flags&syscall.O_WRONLY != 0 {
		code = fuse.EBADF
		return
	}

	content := []byte(setting.value)

	if off >= int64(len(setting.value)) || len(dest) > len(content) {
		res = fuse.ReadResultData(dest)
		code = fuse.EINVAL
		return
	}

	copy(dest, content[off:off+int64(len(dest))])
	res = fuse.ReadResultData(dest)
	return
}

func (setting *fuseConfigurationSetting) Write(dest []byte, off int64) (written uint32, code fuse.Status) {
	if setting.flags&syscall.O_RDONLY != 0 {
		code = fuse.EBADF
		return
	}

	lenDest := len(dest)
	setting.value = setting.value[:off] + string(dest) + setting.value[off+int64(lenDest):]
	written = uint32(lenDest)
	return
}

func (setting *fuseConfigurationSetting) Flush() fuse.Status {
	err := setting.server.Provider.SetSetting(setting.name, setting.value)
	if err != nil {
		return getFuseErrorCode(err)
	}

	err = setting.server.Provider.Save()
	if err != nil {
		return getFuseErrorCode(err)
	}

	return fuse.OK
}

func (setting *fuseConfigurationSetting) Fsync(flags int) fuse.Status {
	return setting.Flush()
}

func (setting *fuseConfigurationSetting) Truncate(size uint64) fuse.Status {
	if setting.flags&syscall.O_RDONLY != 0 {
		return fuse.EBADF
	}

	if size > uint64(len(setting.value)) {
		return fuse.Status(syscall.EFBIG)
	}

	setting.value = setting.value[:size]
	return fuse.OK
}

func (setting *fuseConfigurationSetting) GetAttr(out *fuse.Attr) fuse.Status {
	isDir, err := setting.server.isDir(setting.name)
	if err != nil {
		return getFuseErrorCode(err)
	}

	out = setting.server.getAttr(setting.name, setting.value, isDir, &setting.context.Owner)
	return fuse.OK
}
