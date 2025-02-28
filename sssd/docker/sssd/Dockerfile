FROM sssd-sssd-build-image:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN groupadd -r -g 900 sssd && useradd -r -u 900 -g sssd sssd

RUN mkdir -p /source/sssd && \
    curl -sSL https://github.com/SSSD/sssd/archive/deefe9ad82e8e0057aa77ea5be60a86d223900da.tar.gz | tar -xz -C /source/sssd --strip-components=1
RUN cd /source/sssd && \
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
    make -j && \
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

RUN cd /source/pam-setlocale && \
    make && \
    make install

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
