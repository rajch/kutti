package core

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

const (
	commandSetHostName = "sudo kutti-setupscripts/set-hostname.sh %s"
)

type sshclient struct {
	config *ssh.ClientConfig
}

func (sc *sshclient) RunWithResults(address string, command string) (string, error) {
	client, err := ssh.Dial("tcp", address, sc.config)
	if err != nil {
		return "", fmt.Errorf("Could not connect to address %s:%v ", address, err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Could not create session at address %s:%v ", address, err)
	}
	defer session.Close()

	resultdata, err := session.Output(command)
	if err != nil {
		return string(resultdata), fmt.Errorf("Command '%s' at address %s produced an error:%v ", command, address, err)
	}

	return string(resultdata), nil
}

func newSSHClient(username string, password string) *sshclient {
	return &sshclient{
		config: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
}
