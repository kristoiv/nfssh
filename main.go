package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt, os.Kill)
mainLoop:
	for {
		client, agentConn, err := setupSshClient()
		if err != nil {
			log.Println("Ssh connection error, trying again in five seconds", err)
			agentConn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		// Setup ssh tunnels to the remote host
		quit := make(chan struct{})

		err = setupTunnel(quit, client, Config.LocalNfsdPort, Config.RemoteNfsdPort)
		if err != nil {
			log.Println("Unable to create tunnel, trying again in five seconds", err)
			client.Close()
			agentConn.Close()
			close(quit)
			time.Sleep(5 * time.Second)
			continue
		}

		err = setupTunnel(quit, client, Config.LocalMountdPort, Config.RemoteMountdPort)
		if err != nil {
			log.Println("Unable to create tunnel, trying again in five seconds", err)
			client.Close()
			agentConn.Close()
			close(quit)
			time.Sleep(5 * time.Second)
			continue
		}

		if !nfsIsMounted() {
			mountNfs()
		}

		// Create a channel we can select on for when the ssh connection closes
		clientWaitChan := make(chan struct{})
		go func() {
			client.Wait()
			close(clientWaitChan)
		}()

		go func() {
			for {
				waiter := time.After(1 * time.Second)
				select {
				case <-waiter:
				case <-quit:
					return
				}
				fmt.Printf("\r[METER] 10s avg. throughput: %s/second. 1min avg. throughput: %s/second.", Meter.GetHumanReadablePer10Seconds(), Meter.GetHumanReadablePerMinute())
			}
		}()

		select {
		case <-killSignal:
			if nfsIsMounted() {
				unmountNfs()
			}
			break mainLoop
		case <-clientWaitChan:
			log.Println("Ssh connection closed, restarting in five seconds")
			agentConn.Close()
			close(quit)
			time.Sleep(5 * time.Second)
		}
	}
	log.Fatalln("Interrupted! Unmounted the NFS and quiting..")
}
