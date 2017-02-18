// translator.go
package translator

type Object interface {
	Diff(target io.ReadSeeker) (<-chan Delta, error)
	Patch(target io.WriteCloser, patch <-chan Delta) error
}

type Translator interface {
	TranslateRdiff(rdiff chan<- Delta) translatedRdiff <-chan Delta
}

type Interpreter interface {
	TranslateRdiffToSh(rdiff chan<- Delta) session <-chan Session
}