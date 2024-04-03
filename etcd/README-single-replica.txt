
https://github.com/etcd-io/etcd/issues/17438
https://github.com/etcd-io/etcd/discussions/17394

docker-compose -f docker-compose-single-replica.yaml up

docker exec etcd-etcd-client-1 etcdctl put /key1 value1
docker exec etcd-etcd-client-1 etcdctl get --prefix --keys-only ""
docker exec etcd-etcd-client-1 etcdctl get --prefix ""

docker exec etcd-etcd-client-1 etcdctl endpoint health --write-out=table
docker exec etcd-etcd-client-1 etcdctl endpoint status --write-out=table
docker exec etcd-etcd-client-1 etcdctl member list --write-out=table
docker exec etcd-etcd-client-1 etcdctl alarm list

# write a lot of keys
docker exec etcd-etcd-client-1 sh -c 'for i in $(seq 1 1000); do etcdctl put /key$i value$i; done'


# load test
cd ~/work/etcd
go run ./tools/benchmark put






rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml


sudo rm -rf tmp
mkdir -p tmp/etcd0/
cp -r 20240201_181731_*/var/lib/kubelet/pods/365393b6-f4b1-40b0-aecd-d5785fc6c361/volumes/kubernetes.io~csi/pvc-236a63e9-7c32-4f0a-a9ad-de43c15987af/mount/member tmp/etcd0/member/
chmod -R og-rwx tmp/etcd0/

sudo rm -rf tmp
mkdir -p tmp/etcd0/
cp -r 20240201_181959_*/var/lib/kubelet/pods/620425c8-85d9-4a37-b82b-c3322b91a52c/volumes/kubernetes.io~csi/pvc-772ee157-b830-4e5e-8118-86e9624c7c33/mount/member tmp/etcd0/member/
chmod -R og-rwx tmp/etcd0/

docker-compose -f docker-compose-releases.yaml down -v
docker-compose -f docker-compose-releases.yaml up etcd0

docker exec -it etcd-etcd0-1 ash
etcdctl --cert=/certs/root.pem --key=/certs/root-key.pem --insecure-skip-tls-verify=true get --prefix --keys-only ""
etcdctl --cert=/certs/root.pem --key=/certs/root-key.pem --insecure-skip-tls-verify=true put /key1 value1
etcdctl --cert=/certs/root.pem --key=/certs/root-key.pem --insecure-skip-tls-verify=true get --print-value-only /key1
