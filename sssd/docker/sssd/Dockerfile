
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# SSSD build dependencies derived from the SSSD debian source package
RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        autoconf \
        automake \
        autopoint \
        autopoint \
        bc \
        bind9-dnsutils \
        build-essential \
        check \
        cifs-utils \
        docbook-xml \
        docbook-xsl \
        krb5-config \
        ldap-utils \
        libaugeas-dev \
        libcap-dev \
        libc-ares-dev \
        libcmocka-dev \
        libcollection-dev \
        libcurl4-openssl-dev \
        libdbus-1-dev \
        libdhash-dev \
        libfido2-dev \
        libgdm-dev \
        libglib2.0-dev \
        libini-config-dev \
        libjansson-dev \
        libjose-dev \
        libkeyutils-dev \
        libkrad-dev \
        libkrb5-dev \
        libldap2-dev \
        libldb-dev \
        libltdl-dev \
        libnfsidmap-dev \
        libnl-3-dev \
        libnl-route-3-dev \
        libnss-wrapper \
        libp11-kit-dev \
        libpam-dev \
        libpam-pwquality \
        libpam-runtime \
        libpam-wrapper \
        libpcre2-dev \
        libpopt-dev \
        libsasl2-dev \
        libselinux1-dev \
        libsemanage-dev \
        libsmbclient-dev \
        libssl-dev \
        libsubid-dev \
        libsystemd-dev \
        libtalloc-dev \
        libtdb-dev \
        libtevent-dev \
        libtool \
        libuid-wrapper \
        libunistring-dev \
        libxml2-utils \
        openssh-client \
        openssl \
        pkgconf \
        python3 \
        python3-dev \
        python3-setuptools \
        samba-dev \
        softhsm2 \
        systemd-dev \
        systemtap-sdt-dev \
        uuid-dev \
        xml-core \
        xsltproc

RUN apt-get install --yes --no-install-recommends \
        ca-certificates \
        curl \
        dumb-init \
        openssh-server

RUN groupadd -r -g 900 sssd && useradd -r -u 900 -g sssd sssd

RUN mkdir -p /source/sssd && \
    curl -sSL https://github.com/SSSD/sssd/releases/download/2.10.0-beta2/sssd-2.10.0-beta2.tar.gz | tar -xz -C /source/sssd --strip-components=1 && \
    cd /source/sssd && \
    autoreconf -i && \
    ./configure \
        --with-sssd-user=sssd \
        --disable-static \
        --disable-rpath \
        --prefix=/usr \
        --sysconfdir=/etc \
        --enable-pammoddir=/usr/lib/x86_64-linux-gnu/security/ \
        --enable-nsslibdir=/usr/lib/x86_64-linux-gnu/ \
        --with-systemdunitdir=/lib/systemd/system \
        --without-python2-bindings \
        --with-smb-idmap-interface-version=6 \
    && \
    make install

RUN mkdir -p /source/syslog-redirector && \
    curl -sSL https://github.com/Nordix/syslog-redirector/archive/8e12a28.tar.gz | tar -xz -C /source/syslog-redirector --strip-components=1 && \
    cd /source/syslog-redirector && \
    CFLAGS="-O2 -Wall -Werror" make

COPY files  /

RUN cd /source/nss-fake && \
    make && \
    make install

RUN cd /source/pam-userenv && \
    make && \
    make install

RUN cd /source/getlogin-fake && \
    make

RUN chown sssd:sssd /etc/sssd/sssd.conf && \
    chmod 600 /etc/sssd/sssd.conf && \
    mkdir -p /run/sssd /var/run/nscd/ && \
    chown sssd:sssd /run/sssd /var/run/nscd/ && \
    chmod 700 /run/sssd /var/run/nscd/ && \
    mkdir -p /run/sshd && \
    touch /run/sshd.pid && \
    chown -R sssd:sssd /run/sshd.pid /run/sshd /etc/ssh/ && \
    chmod u-s /usr/bin/passwd

CMD ["dumb-init", "--", "/docker-entrypoint.sh"]
