# Dockerfile from etcd/Dockerfile-release.amd64 but with alpine base image
# for having /bin/sh

FROM alpine:3

ADD bin/etcd /usr/local/bin/
ADD bin/etcdctl /usr/local/bin/
ADD bin/etcdutl /usr/local/bin/

WORKDIR /var/etcd/
WORKDIR /var/lib/etcd/

EXPOSE 2379 2380

# Define default command.
CMD ["/usr/local/bin/etcd"]
