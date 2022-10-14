FROM fedora:36

RUN dnf install -y qemu-pr-helper && dnf clean all
RUN mkdir /usr/local/pflaume
COPY ./connector /usr/bin/connector
COPY dockerfiles/pr-helper/multipath.conf /etc/multipath.conf
COPY --from=pflaume /pflaume/pflaume /usr/local/pflaume/pflaume
COPY --from=pflaume /pflaume/t.out /usr/local/pflaume/t.out
ENTRYPOINT ["/usr/bin/qemu-pr-helper"]
