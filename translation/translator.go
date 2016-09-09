// translator.go
package translation

type TargetProtocol int

const (
	None TargetProtocol = iota
	SSH
	Gonflated
)

type Translator interface {
	TranslateRdiff(rdiff chan<- Delta) (translatedRdiff <-chan Delta, err error)
	TranslateRdiffToRexec(rdiff chan<- Delta) (session *Session, err error)
}
