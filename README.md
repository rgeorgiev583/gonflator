# SYStem CONfiguration File Server
A synthetic filesystem that provides a simple API for configuring and managing the whole system on a remote host.

## Purpose
To provide a simple, portable and uniform means for configuring and manging a remote machine.
Every object in the API is a file, and every action is implemented as a series of filesystem operations.

## Usage
Every action that **sysconfs** supports is accomplished by modifying the filesystem in some way.
In order for all actions to be fully reversible, changes are tracked using Git.

## Installation
Simply execute the following command line:

$ go get "github.com/rgeorgiev583/sysconfs"
