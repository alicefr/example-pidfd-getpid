FROM fedora:36

RUN dnf install -y qemu-pr-helper && dnf clean all
COPY ./connector /usr/bin/connector
COPY dockerfiles/pr-helper/multipath.conf /etc/multipath.conf
COPY --from=pflaume /pflaume/pflaume /usr/bin/pflaume
ENTRYPOINT ["/usr/bin/qemu-pr-helper"]
