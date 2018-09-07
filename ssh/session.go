package session

import (
	"bytes"
)

type Session struct {
	*Socket
	Command string
}

func NewBufferedSession() *Session {
	return &Session{
		Socket: &Socket{
			Stdin:  &bytes.Buffer{},
			Stdout: &bytes.Buffer{},
			Stderr: &bytes.Buffer{},
		},
	}
}
