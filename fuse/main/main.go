package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"github.com/rgeorgiev583/gonflator/augeas"
	"github.com/rgeorgiev583/gonflator/fuse"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	provider, err := augeas.NewConfigurationProvider("/", "", 0)
	if err != nil {
		log.Fatalln(err)
	}

	var options fuse.ConfigurationServerOptions
	server := fuse.NewConfigurationServer(provider, options)

	opts := &nodefs.Options{}
	mountpoint := os.Args[1]

	nfs := pathfs.NewPathNodeFs(server, &pathfs.PathNodeFsOptions{})
	state, _, err := nodefs.MountRoot(mountpoint, nfs.Root(), opts)
	if err != nil {
		panic(fmt.Sprintf("cannot mount %v", err)) // ugh - benchmark has no error methods.
	}

	state.Serve()
}
