
### https://github.com/etcd-io/etcd/pull/15156


# generate certificates
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml

# remove previous containers and volumes if any
docker-compose rm -f
docker volume rm etcd_etcd-data



docker-compose up

# scan TLS versions with sslyze
docker exec -it etcd-shell-1 sslyze etcd0:2379
docker exec -it etcd-shell-1 sslyze etcd0:2380




#### Build image

cd etcd

make         # executable will be in bin/etcd
make test    # run all tests, e2e tests depend on bin/etcd


# or
make test-unit
make test-integration
make test-e2e


# build container
(cd ~/work/etcd && make)
docker-compose build
docker-compose up


# cleanup
docker-compose rm -f
docker volume rm etcd_etcd-data




#### Test on a command line

make
bin/etcd --config-file ~/work/devenvs/etcd/configs/etcd-config-tls-version.yaml



bin/etcd --data-dir /tmp/etcd-datadir --listen-peer-urls https://localhost:2380 --listen-client-urls https://localhost:2379 --advertise-client-urls https://localhost:2379 --peer-trusted-ca-file=/home/tsaarni/work/devenvs/etcd/certs/ca.pem --peer-cert-file=/home/tsaarni/work/devenvs/etcd/certs/etcd.pem --peer-key-file=/home/tsaarni/work/devenvs/etcd/certs/etcd-key.pem --client-cert-auth --trusted-ca-file=/home/tsaarni/work/devenvs/etcd/certs/ca.pem --cert-file=/home/tsaarni/work/devenvs/etcd/certs/etcd.pem --key-file=/home/tsaarni/work/devenvs/etcd/certs/etcd-key.pem --min-tls-version TLS10 --max-tls-version TLS13


sslyze localhost:2379
sslyze localhost:2380
