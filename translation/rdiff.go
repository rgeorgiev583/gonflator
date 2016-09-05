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

type Changer interface {
	Diff(target io.ReadSeeker) (<-chan Delta, error)
	Patch(target io.WriteCloser, diff <-chan Delta) error
}

type Translator interface {
	TranslateToRdiff(diff chan<- Delta) (translatedDiff <-chan Delta, err error)
	TranslateToRexec(diff chan<- Delta) (session *Session, err error)
}
