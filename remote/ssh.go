package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type SSHConn struct {
	Auth   ssh.AuthMethod
	Config *ssh.ClientConfig
	*ssh.Client
}

func DialSSH(addr string, username string, keyFilename string) (conn *SSHConn, err error) {
	conn = &SSHConn{}

	if keyFilename != "" {
		pemKey, err := ioutil.ReadFile(keyFilename)
		if err != nil {
			return
		}

		signer, err := ssh.ParsePrivateKey(pemKey)
		if err != nil {
			return
		}

		conn.Auth = ssh.PublicKeys(signer)
	} else {
		sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			return
		}

		conn.Auth = ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}

	conn.Config = &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			conn.Auth,
		},
	}

	conn.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", addr), conn.Config)
	if err != nil {
		return
	}

	return
}

func (conn *SSHConn) Send(message string, endpoints *Socket) (err error) {
	session, err := conn.NewSession()
	defer session.Close()
	if err != nil {
		return
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return
	}
	go io.Copy(stdin, endpoints.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		return
	}
	go io.Copy(endpoints.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		return
	}
	go io.Copy(endpoints.Stderr, stderr)

	return session.Run(name)
}
