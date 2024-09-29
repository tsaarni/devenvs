FROM ubuntu:22.04
ENV DEBIAN_FRONTEND noninteractive

ARG USERNAME=tsaarni
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# basic tools
RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        autoconf \
        automake \
        bear \
        bind9-host \
        build-essential \
        ca-certificates \
        curl \
        dumb-init \
        file \
        gdb \
        git \
        jq \
        iputils-ping \
        ldap-utils \
        less \
        libltdl-dev \
        locales \
        openssh-server \
        pkg-config \
        python3 \
        python3-pip \
        python3-venv \
        rsyslog \
        silversearcher-ag \
        sudo \
        vim-tiny \
        whois

# golang
RUN curl -s https://go.dev/dl/?mode=json | \
    jq -r '.[0].files[] | select(.os == "linux" and .arch == "amd64") | .filename' | \
    xargs -I{} curl -L https://go.dev/dl/{} -o go.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

# packages specific to this project
RUN apt-get install --yes --no-install-recommends \
        groff-base \
        libargon2-dev \
        libevent-dev \
        libsasl2-dev \
        libldap-dev \
        libpam-dev \
        libssl-dev \
        zlib1g-dev

RUN locale-gen en_US.UTF-8

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --create-home --shell /bin/bash --uid $USER_UID --gid $USER_GID $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME

CMD ["/usr/bin/dumb-init", "--", "sleep", "infinity"]
