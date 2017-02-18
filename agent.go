package agent

import "io"

type AgentType int

const (
	None AgentType = iota
	Standard
	FUSE
)

type TargetType int

const (
	None TargetType = iota
	Augeas
	Rsync
	Daemon
	SSH
)

type Syncer interface {
	Diff() (<-chan Delta, error)
	Fetch()
}
