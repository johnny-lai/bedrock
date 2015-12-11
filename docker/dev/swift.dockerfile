FROM johnnylai/swift:2.2

# Extra tooling
RUN apt-get install -y vim telnet jq

COPY docker/dev/sudoers /etc/sudoers

COPY docker/dev/Makefile /go/Makefile

COPY docker/dev/entrypoint.sh /entrypoint.sh

COPY scripts/cluster.sh /bin/cluster.sh

ENTRYPOINT ["/entrypoint.sh"]
