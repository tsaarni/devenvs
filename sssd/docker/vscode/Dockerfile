FROM ubuntu:22.04
ENV DEBIAN_FRONTEND noninteractive

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        autoconf \
        automake \
        bear \
        bind9-host \
        build-essential \
        ca-certificates \
        cracklib-runtime \
        curl \
        gdb \
        gettext \
        git \
        iputils-ping \
        ldap-utils \
        less \
        libldap-dev \
        libpam-dev \
        locales \
        openssh-server \
        rsyslog \
        silversearcher-ag \
        sudo \
        dumb-init \
        vim-tiny \
        zlib1g-dev \
        libssl-dev

RUN locale-gen en_US.UTF-8

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --create-home --shell /bin/bash --uid $USER_UID --gid $USER_GID $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

RUN curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh && \
    usermod -aG docker $USERNAME

USER $USERNAME

CMD ["/usr/bin/dumb-init", "--", "sleep", "infinity"]
