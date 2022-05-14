FROM ubuntu:18.04
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install --yes --no-install-recommends \
        ca-certificates \
        cracklib-runtime \
        curl \
        less \
        libnss-sss \
        libpam-sss \
        openssh-server \
        sssd \
        sssd-tools

RUN curl -fsSL https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/v1.5.4505/files/docker/systemctl3.py -o /usr/bin/systemctl3.py && \
    curl -fsSL https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/v1.5.4505/files/docker/journalctl3.py -o /usr/bin/journalctl3.py && \
    chmod +x /usr/bin/systemctl3.py /usr/bin/journalctl3.py && \
    ln -s /usr/bin/systemctl3.py /usr/bin/systemctl && \
    ln -s /usr/bin/journalctl3.py /usr/bin/journalctl

COPY files /

CMD ["/docker-entrypoint.sh"]
