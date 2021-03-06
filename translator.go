// translator.go
package translator

type TargetProtocol int

const (
	None TargetProtocol = iota
	Rsync
	SSH
)

type Translator interface {
	TranslateRdiff(rdiff chan<- Delta) (translatedRdiff <-chan Delta, err error)
	TranslateRdiffToRexec(rdiff chan<- Delta) (session *Session, err error)
}
