package remote

type Socket struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Communicator interface {
	Send(message string, endpoints *Socket) error
}
