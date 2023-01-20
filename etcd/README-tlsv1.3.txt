

rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml

docker-compose up

# scan TLS versions
docker exec -it etcd-shell-1 sslyze etcd0:2379
docker exec -it etcd-shell-1 sslyze etcd0:2380
