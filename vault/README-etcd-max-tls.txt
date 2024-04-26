https://github.com/hashicorp/vault/issues/26659





make tools
make dev

# enable TLS for etcd by uncommenting the TLS startup command in the docker-compose.yml file

docker-compose stop etcd0 etcd1 etcd2
docker-compose up etcd0 etcd1 etcd2


vault server -config=$HOME/work/devenvs/vault/configs/vault-config-etcd-with-tls.hcl
