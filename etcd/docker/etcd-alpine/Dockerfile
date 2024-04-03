# https://github.com/etcd-io/etcd/releases
# NOTE: Newer images after etcd:v3.5.6 do not have /bin/sh anymore and we cannot use env vars in entrypoint command line.
#       Therefore we build our own image.
#FROM quay.io/coreos/etcd:v3.5.12 as etcd
FROM quay.io/coreos/etcd:v3.5.7 as etcd
FROM alpine:3.19
COPY --from=etcd /usr/local /usr/local
