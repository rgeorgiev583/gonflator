# Golang CONfiguration File Server
A synthetic filesystem (a.k.a. file server) written in Go that provides a simple
API for configuring and managing software on a remote host.

## Purpose
To provide a simple, portable and uniform means for configuring and managing a
remote machine.
Every object in the API is a file, and every action is implemented as a
modification of the filesystem (FUSE call when mounted as a FUSE filesystem, or
a Git diff when using a Git hook).

## Usage
Every action that **sysconfs** supports is accomplished by modifying the
filesystem in some way.
In order for all actions to be fully reversible, it is possible for changes to
    be tracked using a VCS.  Additionally, a Git hook is provided for updating
    target machines before pushing changes to upstream.

    ## Installation
    Simply execute the following command line:

        $ go get github.com/rgeorgiev583/gonfs
