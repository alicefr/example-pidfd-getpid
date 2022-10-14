FROM disk as builder
FROM pflaume as pflaume
FROM fedora:36

RUN dnf install -y qemu-system-x86-core strace && dnf clean all

RUN mkdir /usr/local/pflaume
COPY --from=builder /disk/disk.img /disk/disk.img
COPY --from=pflaume /pflaume/pfloader /usr/local/bin/pfloader
COPY --from=pflaume /pflaume/bpf.out  /usr/local/bin/bpf.out
COPY dockerfiles/qemu/start-qemu.sh /start-qemu.sh 
RUN mkdir -p /home/qemu && chown qemu:kvm -R /home/qemu && chmod +x /start-qemu.sh
ENTRYPOINT ["/start-qemu.sh"]

