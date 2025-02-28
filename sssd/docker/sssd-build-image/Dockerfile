FROM ubuntu:24.04
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        xz-utils \
        wget

# SSSD build dependencies derived from the SSSD debian source package.
# https://packages.debian.org/sid/sssd
# Note that this is not the version of SSSD that we are building,
# but the dependencies are likely the same so we use it as a reference.

COPY files /

RUN wget http://deb.debian.org/debian/pool/main/s/sssd/sssd_2.9.5-3.debian.tar.xz -O /tmp/sssd.tar.xz && \
    /extract-build-deps.sh /tmp/sssd.tar.xz | xargs apt-get install --yes --no-install-recommends


# Additional dependencies for the SSSD build, added in newer versions of SSSD
RUN apt-get install --yes --no-install-recommends \
    libcap-dev


# Additinoal tools for debugging and development
RUN apt-get install --yes --no-install-recommends \
    bear \
    ca-certificates \
    curl \
    dumb-init \
    gdb \
    ldap-utils \
    less \
    locales \
    openssh-server \
    rsyslog \
    silversearcher-ag \
    sudo \
    vim-tiny

RUN locale-gen en_US.UTF-8
