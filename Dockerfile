FROM fedora:36

RUN dnf install socat nmap-ncat strace -y && dnf remove all

COPY ./proxy /usr/bin/proxy
COPY ./connector /usr/bin/connector
