FROM golang:1.22
ENV DEBIAN_FRONTEND noninteractive

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        bind9-host \
        ca-certificates \
        curl \
        git \
        httpie \
        iputils-ping \
        jq \
        less \
        locales \
        siege \
        silversearcher-ag \
        sudo \
        tini \
        vim-tiny

RUN locale-gen en_US.UTF-8

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --create-home --shell /bin/bash --uid $USER_UID --gid $USER_GID $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME

# run as user when installing go tools
RUN go install -v github.com/go-delve/delve/cmd/dlv@v1.22.0 && \
    go install -v golang.org/x/tools/gopls@latest

CMD ["/usr/bin/tini", "--", "sleep", "infinity"]
