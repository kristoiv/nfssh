package main

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func setupSshClient() (*ssh.Client, net.Conn, error) {
	agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, nil, err
	}
	agentClient := agent.NewClient(agentConn)
	client, err := ssh.Dial("tcp", Config.Host, &ssh.ClientConfig{
		User: Config.Username,
		Auth: []ssh.AuthMethod{ssh.PublicKeysCallback(agentClient.Signers)},
	})
	return client, agentConn, err
}
