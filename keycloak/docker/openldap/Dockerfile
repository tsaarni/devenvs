FROM ubuntu:24.04 AS builder
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        build-essential \
        libssl-dev

COPY sslkeylog /src/sslkeylog
RUN cd /src/sslkeylog && \
    make

COPY check-password /src/check-password
RUN cd /src/check-password && \
    make


FROM ubuntu:24.04
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        ldap-utils \
        slapd \
        whois
    # note: whois package includes mkpasswd which is used in ldif templates

COPY --from=builder /src/sslkeylog/libsslkeylog.so /
COPY --from=builder /src/check-password/check-password.so /usr/lib/ldap/

COPY files /

EXPOSE 389 636

CMD ["/docker-entrypoint.sh"]
