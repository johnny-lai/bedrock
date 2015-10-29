FROM golang

# Install go
RUN go get -u github.com/derekparker/delve/cmd/dlv && \
	go get github.com/Masterminds/glide

# Install Kubernetes
RUN apt-get update && apt-get install -y lsb-release vim sudo

RUN curl -sSL https://get.kismatic.com/kubernetes/master.sh | sh

# Extra tooling
RUN apt-get install -y telnet jq

COPY sudoers /etc/sudoers

COPY Makefile /go/Makefile

COPY wait-for-pod.sh /bin/wait-for-pod.sh

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
