FROM quay.io/coreos/etcd:v3.5.12 as etcd
FROM alpine:3.19
COPY --from=etcd /usr/local /usr/local
