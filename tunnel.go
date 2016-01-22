package main

import (
	"log"
	"net"
	"strconv"
	"sync"

	"golang.org/x/crypto/ssh"
)

func setupTunnel(quit <-chan struct{}, client *ssh.Client, lport, rport int) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(lport))
	if err != nil {
		return err
	}
	go func() {
		defer listener.Close()

		connChan := make(chan net.Conn)
		errChan := make(chan error)

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					errChan <- err
					return
				}
				connChan <- conn
			}
		}()

		for {
			select {
			case conn := <-connChan:
				go tunnelConnectionToPortUsingClient(conn, rport, client)
			case err := <-errChan:
				log.Println("Listener failed", err)
			case <-quit:
				return
			}
		}
	}()
	return nil
}

func tunnelConnectionToPortUsingClient(localConn net.Conn, rport int, client *ssh.Client) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+strconv.Itoa(rport))
	if err != nil {
		panic(err)
	}

	remoteConn, err := client.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	copyConn := func(writer, reader net.Conn) {
		defer wg.Done()
		CopyAndMeasureThroughput(writer, reader)
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
	go func() {
		wg.Wait()
		localConn.Close()
		remoteConn.Close()
	}()
}
