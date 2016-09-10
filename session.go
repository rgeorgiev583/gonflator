package session

type Session struct {
	*Socket
	Command string
}
