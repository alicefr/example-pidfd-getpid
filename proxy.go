package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const socketAddr = "proxy.sock"

func fromClientToDaemon(ctx context.Context, conn net.Conn, fd int) {
	for {
		select {
		case <-ctx.Done():
			break
		default: // Read from the client
			reply := make([]byte, 1024)
			_, err := conn.Read(reply)
			if err != nil {
				log.Fatal("reading from client: %v:", err)
			}
			log.Info("got from client:", string(reply))
			// Write to the privileged daemon
			_, err = syscall.Write(fd, reply)
			if err != nil {
				log.Fatal("failed writing privileged daemon: %v:", err)
			}
		}
	}
}

func fromDaemonToClient(ctx context.Context, conn net.Conn, fd int) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			// Read from the privileged daemon
			reply := make([]byte, 1024)
			_, err := syscall.Read(fd, reply)
			if err != nil {
				log.Fatal("failed reading from the daemon: %v", err)
			}
			log.Info("got from the daemon:", string(reply))
			// Write to the client
			_, err = conn.Write(reply)
			if err != nil {
				log.Fatal("failed writing to client: %v:", err)
			}
		}
	}
}

func main() {
	// Create fd to connect to the privileged daemon
	fd, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Pid=%d FD=%d", os.Getpid(), fd)
	ctx, cancel := context.WithCancel(context.Background())
	if err := os.RemoveAll(socketAddr); err != nil {
		log.Fatal(err)
	}
	defer cancel()

	l, err := net.Listen("unix", socketAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	defer os.Remove(socketAddr)

	conn, err := l.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	defer conn.Close()
	log.Info("Accepted connection")
	go fromClientToDaemon(ctx, conn, fd)
	go fromDaemonToClient(ctx, conn, fd)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	<-c
	log.Info("Terminating")
}
