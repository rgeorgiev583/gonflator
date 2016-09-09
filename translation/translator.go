// translator.go
package translation

type Translator interface {
	TranslateRdiff(rdiff chan<- Delta) (translatedRdiff <-chan Delta, err error)
	TranslateRdiffToRexec(rdiff chan<- Delta) (session *Session, err error)
}
