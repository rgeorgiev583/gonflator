package delta

import "io"

type DeltaProtocol int

const (
	None Protocol = iota
	FUSE
	Rsync
	Git
)

type DeltaType int

const (
	DeltaUnmodified DeltaType = iota
	DeltaAdded
	DeltaDeleted
	DeltaModified
	DeltaRenamed
	DeltaCopied
)

type Delta struct {
	rsync.Operation
	OldPath, NewPath string
	Type             DeltaType
}

type OptionalDelta struct {
	Delta
	Err error
}

type Object interface {
	Diff(target io.ReadSeeker) (<-chan Delta, error)
	Patch(target io.WriteCloser, patch <-chan Delta) error
}
