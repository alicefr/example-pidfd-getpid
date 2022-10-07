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

func fromClientToDaemon(ctx context.Context, conn *net.UnixConn, fd int) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			msg, oob := make([]byte, 1024), make([]byte, 128)
			// Read from the client
			n, oobn, _, _, err := conn.ReadMsgUnix(msg, oob)
			if err != nil {
				log.Fatalf("reading from client: %v", err)
			}
			// Write to the privileged daemon
			err = syscall.Sendmsg(fd, msg[:n], oob[:oobn], nil, 0)
			if err != nil {
				log.Fatalf("failed writing privileged daemon: %v:", err)
			}
		}
	}
}

func fromDaemonToClient(ctx context.Context, conn *net.UnixConn, fd int) {
	for {
		select {
		case <-ctx.Done():
			break
		default:
			// Read from the privileged daemon
			reply := make([]byte, 1024)
			n, err := syscall.Read(fd, reply)
			if err != nil {
				log.Fatalf("failed reading from the daemon: %v, read bytes: %d", err, n)
			}
			// Write to the client
			_, err = conn.Write(reply)
			if err != nil {
				log.Fatalf("failed writing to client: %v:", err)
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
	defer cancel()

	syscall.Unlink(socketAddr)
	addr, err := net.ResolveUnixAddr("unix", socketAddr)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	defer syscall.Unlink(socketAddr)

	conn, err := l.AcceptUnix()
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
