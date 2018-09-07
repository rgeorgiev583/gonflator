# Design of the **sysconfs** file server
**sysconfs** is a synthetic Unix filesystem that provides a high-level,
abstract, object- and configuration-oriented representation of the state of all
software on a machine (be it physical or virtual), and supports doing this for
multiple machines simultaneously.

## Requirements
* a Unix-like OS on the admin machine (on which the filesystem is mounted)
* FUSE (for Linux/*BSD)
* golang
* Git

## Architecture
The software that implements **sysconfs** consists of three parts:
* **sysconfsd** - the daemon that provides the file server as a userspace service
* **sysconfctl** - a CLI utility for easy management of the **sysconfsd** daemon
* a Git hook for easily deploying changes to the target machine - as simple as
  a `git push` which translates the pushed changes in the filesystem to SSH
  commands which execute the relevant commands to implement said changes.

## What the heck is a state

A state of a machine is defined as the representation of all configuration
on a machine at a given moment.  This includes persistent configuration on the
physical file system (aka configuration files) as well as runtime configuration
(aka procfs, sysfs, udev, tmpfs, /run, /tmp) and whatever else can be exported
as a filesystem).

## Structure
The basic (and stable) filesystem hierarchy is structured as follows:

* / - root node
* /conf - sysconfs current filesystem configuration
* /*node* - hostname or IP address of the target machine
* /*node*/conf - persistent configuration
* /*node*/conf/sys - persistent system configuration
* /*node*/conf/*pkgname* - persistent software configuration for package *pkgname*
* /*node*/run - runtime configuration
* /*node*/run/sys - runtime system configuration
* /*node*/run/*pkgname* - runtime software configuration for package *pkgname*

The configuration for `sys` and *`pkgname`* nodes depends on the particular
configuration server implementation.

## 
