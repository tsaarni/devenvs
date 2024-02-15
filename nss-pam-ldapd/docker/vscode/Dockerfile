FROM ubuntu:22.04
ENV DEBIAN_FRONTEND noninteractive

ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        autoconf \
        automake \
        bind9-host \
        build-essential \
        ca-certificates \
        cracklib-runtime \
        curl \
        gdb \
        git \
        iputils-ping \
        ldap-utils \
        less \
        libldap-dev \
        libpam-dev \
        locales \
        openssh-server \
        silversearcher-ag \
        sudo \
        tini \
        vim-tiny
RUN locale-gen en_US.UTF-8

COPY files /

RUN groupadd --gid $USER_GID $USERNAME && \
    useradd --create-home --shell /bin/bash --uid $USER_UID --gid $USER_GID $USERNAME && \
    echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

USER $USERNAME

CMD ["/usr/bin/tini", "--", "sleep", "infinity"]
