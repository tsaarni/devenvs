FROM ubuntu:19.10

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y ca-certificates inotify-tools lsof gdbserver lldb \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /tmp/* /var/tmp/* \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["sleep", "infinity"]
