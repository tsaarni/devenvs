FROM ubuntu:18.04 AS builder
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        build-essential \
        libssl-dev
COPY sslkeylog /
RUN make


FROM ubuntu:18.04
ENV DEBIAN_FRONTEND noninteractive

COPY --from=builder /libsslkeylog.so /
RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        ldap-utils \
        slapd \
        whois
# note: whois package includes mkpasswd which is used in ldif templates

COPY files /

EXPOSE 389 636

CMD ["/docker-entrypoint.sh"]
