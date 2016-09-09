package remote

type Socket struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Communicator interface {
	Send(message string, endpoints *Socket) error
}

func DialURL(url url.URL, proto Protocol) (comm Communicator, err error) {
	if url.User.Username != "" {
		user := url.User.Username
	} else {
		user := os.Getenv("USER")
	}

	switch url.Scheme {
	default:
		return DialSSH(url.Host, user, "")
	}
}
