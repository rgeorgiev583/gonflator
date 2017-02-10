// translator.go
package translator

type TargetProtocol int

type ShellType int

const (
	None TargetProtocol = iota
	Rsync
	SSH
)

type Translator interface {
	TranslateRdiff(rdiff chan<- Delta) translatedRdiff <-chan Delta
	TranslateRdiffToSh(rdiff chan<- Delta) session <-chan Session
}
