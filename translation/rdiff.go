package translation

import "io"

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

type Session struct {
	Socket
	Command string
}

type Object interface {
	Diff(target io.ReadSeeker) (<-chan Delta, error)
	Patch(target io.WriteCloser, patch <-chan Delta) error
}

type Translator interface {
	TranslateRdiff(rdiff chan<- Delta) (translatedRdiff <-chan Delta, err error)
	TranslateRdiffToRexec(rdiff chan<- Delta) (session *Session, err error)
}
