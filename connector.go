package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/oraoto/go-pidfd"
)

var pid = flag.Int("pid", 0, "proxy pid")
var fd = flag.Int("fd", 0, "proxy fd")

const sock = "/tmp/priv.sock"

func main() {
	flag.Parse()
	if *pid == 0 {
		flag.PrintDefaults()
		log.Fatal("Provide the pid of the proxy")
	}
	if *fd == 0 {
		flag.PrintDefaults()
		log.Fatal("Provide the fd of the proxy")
	}
	if _, err := os.Stat(sock); err != nil {
		log.Fatal(err)
	}
	log.Infof("Connect process pid=%d with fd=%d ", *pid, *fd)
	p, err := pidfd.Open(*pid, 0)
	if err != nil {
		panic(err)
	}
	listenFd, err := p.GetFd(*fd, 0)
	if err != nil {
		panic(err)
	}
	serverAddr := &unix.SockaddrUnix{
		Name: sock,
	}
	err = unix.Connect(listenFd, serverAddr)
	if err != nil {
		log.Fatalf("fail connecting: %v", err)
	}
	fmt.Printf("Succesfully connected proxy to pr-helper")
}
