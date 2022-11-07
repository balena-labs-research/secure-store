package rclone

import (
	"context"
	"errors"
	"sync"

	_ "github.com/rclone/rclone/backend/crypt"
	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/cmd/mount"
	"github.com/rclone/rclone/cmd/mountlib"
	"github.com/rclone/rclone/fs/rc"
	"github.com/rclone/rclone/vfs/vfsflags"
)

var (
	liveMounts = map[string]*mountlib.MountPoint{}
	mountMu    sync.Mutex
)

func CreateRcloneMount(ctx context.Context, in rc.Params) (out rc.Params, err error) {
	mountPoint, err := in.GetString("mountPoint")
	if err != nil {
		return nil, err
	}

	vfsOpt := vfsflags.Opt
	err = in.GetStructMissingOK("vfsOpt", &vfsOpt)
	if err != nil {
		return nil, err
	}

	mountOpt := mountlib.Opt
	err = in.GetStructMissingOK("mountOpt", &mountOpt)
	if err != nil {
		return nil, err
	}

	mountType, err := in.GetString("mountType")

	mountMu.Lock()
	defer mountMu.Unlock()

	if err != nil {
		mountType = ""
	}
	_, mountFn := mountlib.ResolveMountMethod(mountType)
	if mountFn == nil {
		return nil, errors.New("mount option specified is not registered, or is invalid")
	}

	// Get Fs.fs to be mounted from fs parameter in the params
	fdst, err := rc.GetFs(ctx, in)
	if err != nil {
		return nil, err
	}

	mnt := mountlib.NewMountPoint(mountFn, mountPoint, fdst, &mountOpt, &vfsOpt)

	_, err = mnt.Mount()

	if err != nil {
		return nil, err
	}

	// Add mount to list if mount point was successfully created
	liveMounts[mountPoint] = mnt

	return nil, nil
}
