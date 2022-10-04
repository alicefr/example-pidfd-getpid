# Use pidfd_getfd in container

This repository contains an example how to use the syscall [pidfd_getfd](https://man7.org/linux/man-pages/man2/pidfd_getfd.2.html) (available from kernel 5.6) to connect an unprivileged proxy to a privileged daemon running in two separate containers.

## Build
Compile and build the image `getfd`:
```bash
make image
```
The image getfd contains the `proxy` and `connector` binary. The proxy is the process that will be connected to the privileged daemon (in the example socat). The connector is the program that connects the proxy to the privileged daemon.

```bash
docker run --name unprivileged -td getfd /usr/bin/proxy
pid=$(docker inspect --format "{{.State.Pid}}" unprivileged)
docker run -ti -d --name privileged \
  --pid host \
  --cap-add SYS_PTRACE \
  getfd \
  socat UNIX-LISTEN:/tmp/priv.sock -
docker exec -ti privileged connector -pid=$pid -fd=3
Succesfully connected proxy to pr-helper
```
Test the connection:
From the unprivileged container
```bash
docker exec -ti unprivileged nc -Uv proxy.sock
Ncat: Version 7.93 ( https://nmap.org/ncat )
Ncat: Connected to proxy.sock.
test the connection
```
Checking the log of the privileged container:
```bash
$ docker logs privileged
test the connection
```
