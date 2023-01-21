
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




#### Build own image version

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
