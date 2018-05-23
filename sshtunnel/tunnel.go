package sshtunnel

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/golang/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func Tunnel(user string, rhost string, rport int, endpoint string) *exec.Cmd {
	options := []string{}
	options = append(options, "-l", user)
	options = append(options, rhost)
	options = append(options, "-N")
	options = append(options, "-p", strconv.Itoa(rport))
	options = append(options, "-L"+endpoint)
	cmd := exec.Command("/usr/bin/ssh", options...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Start()

	return cmd
}

func TunnelTest(user string, port int, rhost string, rport int, ehost string, eport int) {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: port,
	}

	serverEndpoint := &Endpoint{
		Host: rhost,
		Port: 22,
	}

	remoteEndpoint := &Endpoint{
		Host: ehost,
		Port: eport,
	}

	khc, err := knownhosts.New("/home/something/.ssh/known_hosts")
	if err != nil {
		panic(err)
	}
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		HostKeyCallback: khc,
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
