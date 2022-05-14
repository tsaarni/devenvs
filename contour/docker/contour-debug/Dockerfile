FROM ubuntu:19.10

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y ca-certificates golang git \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /tmp/* /var/tmp/* \
    && rm -rf /var/lib/apt/lists/*

RUN go get -u github.com/go-delve/delve/cmd/dlv
RUN ln -s /root/go/bin/dlv /usr/bin

ENTRYPOINT ["sleep", "infinity"]
