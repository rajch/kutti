package sshclient

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/containerd/console"
	"github.com/rajch/kutti/internal/pkg/kuttilog"
	"github.com/rajch/kutti/pkg/core"
	"golang.org/x/crypto/ssh"
)

const (
	commandSetHostName = "sudo rw-installscripts/set-hostname.sh %s"
)

type sshclient struct {
	config *ssh.ClientConfig
}

// RunWithResults connects to the specified address, runs the specified command, and
// fetches the results.
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

// RunInterativeShell connects to the specified address and runs an interactive
// shell.
func (sc *sshclient) RunInterativeShell(address string) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := sc.runclient(ctx, address); err != nil {
			kuttilog.Print(0, err)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}

// Copied almost verbatim from https://gist.github.com/atotto/ba19155295d95c8d75881e145c751372
// Thanks, Ato Araki (atotto@github)
func (sc *sshclient) runclient(ctx context.Context, address string) error {
	conn, err := ssh.Dial("tcp", address, sc.config)
	if err != nil {
		return fmt.Errorf("cannot connect to %v: %v", address, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %v", err)
	}
	defer session.Close()

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	/*
		fd := int(os.Stdin.Fd())
		state, err := terminal.MakeRaw(fd)
		if err != nil {
			return fmt.Errorf("terminal make raw: %s", err)
		}
		defer terminal.Restore(fd, state)
	*/
	current := console.Current()
	defer current.Reset()

	err = current.SetRaw()
	if err != nil {
		return fmt.Errorf("terminal make raw: %s", err)
	}

	// fd2 := int(os.Stdout.Fd())
	// w, h, err := terminal.GetSize(fd2)
	// if err != nil {
	// 	return fmt.Errorf("terminal get size: %s", err)
	// }

	ws, err := current.Size()
	if err != nil {
		return fmt.Errorf("terminal get size: %s", err)
	}

	h := int(ws.Height)
	w := int(ws.Width)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	term := os.Getenv("TERM")
	if term == "" {
		term = "xterm-256color"
	}
	if err := session.RequestPty(term, h, w, modes); err != nil {
		return fmt.Errorf("session xterm: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("session shell: %s", err)
	}

	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("ssh: %s", err)
	}
	return nil
}

// New creates a new SSH client with password authentication, and no host key check
func New(username string, password string) core.SSHClient {
	return &sshclient{
		config: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			Timeout:         5 * time.Second,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
	}
}
