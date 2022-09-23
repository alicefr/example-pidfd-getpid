package main

import (
	"net"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const socketAddr = "proxy.sock"

func main() {
	// Create fd to connect to the privileged daemon
	fd, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Pid=%d FD=%d", os.Getpid(), fd)

	if err := os.RemoveAll(socketAddr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", socketAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatal("accept error:", err)
	}
	defer conn.Close()
	log.Info("Accepted connection")
	for {
		reply := make([]byte, 1024)
		_, err = conn.Read(reply)
		log.Info("got:", string(reply))
		// Test writing on the privileged daemon
		_, err := syscall.Write(fd, reply)
		if err != nil {
			log.Fatal("Failed writing priv daemon: %v:", err)
		}

	}
}
