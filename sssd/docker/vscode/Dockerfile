FROM sssd-sssd-build-image:latest

ENV DEBIAN_FRONTEND=noninteractive

ARG USERNAME=ubuntu

RUN echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME && \
    chmod 0440 /etc/sudoers.d/$USERNAME

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        bear \
        bind9-host \
        curl \
        gettext \
        git \
        iputils-ping \
        ldap-utils \
        ltrace \
        less \
        silversearcher-ag \
        strace \
        sudo \
        tini \
        vim-tiny
RUN locale-gen en_US.UTF-8

RUN curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh && \
    usermod -aG docker $USERNAME

COPY files /

USER $USERNAME

CMD ["/usr/bin/dumb-init", "--", "sleep", "infinity"]
