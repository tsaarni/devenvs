
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        augeas-lenses \
        autoconf \
        automake \
        autopoint \
        autotools-dev \
        bind9-dnsutils \
        bind9-host \
        bind9-libs \
        bsdextrautils \
        build-essential \
        ca-certificates \
        check \
        cifs-utils \
        curl \
        dbus \
        dbus-user-session \
        dconf-gsettings-backend \
        dconf-service \
        debhelper \
        debugedit \
        dh-apparmor \
        dh-autoreconf \
        dh-python \
        dh-strip-nondeterminism \
        distro-info-data \
        dnsutils \
        docbook-xml \
        docbook-xsl \
        dumb-init \
        dwz \
        faketime \
        file \
        gettext \
        gettext-base \
        gir1.2-gdm-1.0 \
        gir1.2-glib-2.0 \
        git \
        gnutls-bin \
        groff-base \
        icu-devtools \
        intltool-debian \
        krb5-config \
        krb5-multidev \
        ldap-utils \
        libapparmor1 \
        libarchive-zip-perl \
        libargon2-1 \
        libaugeas-dev \
        libaugeas0 \
        libbasicobjects-dev \
        libbasicobjects0 \
        libblkid-dev \
        libbz2-dev \
        libc-ares-dev \
        libc-ares2 \
        libcap-dev \
        libcbor0.8 \
        libcmocka-dev \
        libcmocka0 \
        libcollection-dev \
        libcollection4 \
        libcryptsetup12 \
        libcurl4-openssl-dev \
        libdbus-1-dev \
        libdconf1 \
        libdebhelper-perl \
        libdevmapper1.02.1 \
        libdhash-dev \
        libdhash1 \
        libdw1 \
        libedit2 \
        libelf1 \
        libexpat1-dev \
        libfaketime \
        libffi-dev \
        libfido2-1 \
        libfile-stripnondeterminism-perl \
        libgdm-dev \
        libgdm1 \
        libgirepository-1.0-1 \
        libglib2.0-bin \
        libglib2.0-data \
        libglib2.0-dev \
        libglib2.0-dev-bin \
        libgnutls-dane0 \
        libgssrpc4 \
        libicu-dev \
        libini-config-dev \
        libini-config5 \
        libip4tc2 \
        libjose-dev \
        libjs-jquery \
        libjs-sphinxdoc \
        libjs-underscore \
        libjson-c5 \
        libkadm5clnt-mit12 \
        libkadm5srv-mit12 \
        libkdb5-10 \
        libkeyutils-dev \
        libkmod2 \
        libkrad-dev \
        libkrb5-dev \
        libldap-dev \
        libldap2-dev \
        libldb-dev \
        libltdl-dev \
        libltdl7 \
        libmagic-mgc \
        libmagic1 \
        libmaxminddb0 \
        libmount-dev \
        libnfsidmap-dev \
        libnfsidmap1 \
        libnl-3-200 \
        libnl-3-dev \
        libnl-route-3-200 \
        libnl-route-3-dev \
        libnss-wrapper \
        libp11-kit-dev \
        libpam-systemd \
        libpam-wrapper \
        libpam0g-dev \
        libpath-utils1 \
        libpcre16-3 \
        libpcre2-16-0 \
        libpcre2-dev \
        libpcre2-posix3 \
        libpcre3-dev \
        libpcre32-3 \
        libpcrecpp0v5 \
        libpipeline1 \
        libpopt-dev \
        libpython3-dev \
        libpython3.10-dev \
        libref-array-dev \
        libref-array1 \
        libsasl2-dev \
        libselinux1-dev \
        libsemanage-dev \
        libsepol-dev \
        libsigsegv2 \
        libsmbclient \
        libsmbclient-dev \
        libsofthsm2 \
        libsub-override-perl \
        libsubunit-dev \
        libsubunit0 \
        libsystemd-dev \
        libtalloc-dev \
        libtdb-dev \
        libtevent-dev \
        libtool \
        libuchardet0 \
        libuid-wrapper \
        libunbound8 \
        libunistring-dev \
        libuv1 \
        libwbclient-dev \
        libxml2 \
        libxml2-dev \
        libxml2-utils \
        libxslt1.1 \
        lsb-release \
        m4 \
        man-db \
        openssh-client \
        openssh-server \
        pkg-config \
        po-debconf \
        python3-dev \
        python3-distutils \
        python3-lib2to3 \
        python3-pkg-resources \
        python3-setuptools \
        python3.10-dev \
        samba-dev \
        sgml-base \
        sgml-data \
        softhsm2 \
        softhsm2-common \
        sudo \
        systemd \
        systemd-sysv \
        systemtap-sdt-dev \
        uuid-dev \
        winbind \
        xml-core \
        xsltproc

RUN groupadd -r -g 900 sssd && useradd -r -u 900 -g sssd sssd

RUN mkdir -p /source/sssd && \
    curl -sSL https://github.com/SSSD/sssd/archive/79c384f.tar.gz | tar -xz -C /source/sssd --strip-components=1 && \
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
    && \
    make install && \
    rm /usr/lib/x86_64-linux-gnu/ldb/modules/ldb/memberof.la


RUN mkdir -p /source/syslog-redirector && \
    curl -sSL https://github.com/Nordix/syslog-redirector/archive/21d2d3a.tar.gz | tar -xz -C /source/syslog-redirector --strip-components=1 && \
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
    chown sssd:sssd /run/sssd && \
    mkdir -p /run/sssd /var/run/nscd/ && \
    chown sssd:sssd /run/sssd /var/run/nscd/ && \
    chmod 700 /run/sssd /var/run/nscd/ && \
    mkdir -p /run/sshd && \
    touch /run/sshd.pid && \
    chown -R sssd:sssd /run/sshd.pid /run/sshd /etc/ssh/ && \
    chmod u-s /usr/bin/passwd

CMD ["dumb-init", "--", "/docker-entrypoint.sh"]
