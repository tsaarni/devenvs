




docker-compose -f docker-compose-limited-disk-space.yaml rm --force
docker volume rm vault_etcd-data
docker-compose -f docker-compose-limited-disk-space.yaml up


export VAULT_ADDR=http://127.0.0.1:8200

# initialize vault

http POST $VAULT_ADDR/v1/sys/init secret_shares:=1 secret_threshold:=1 > init.json
export UNSEAL_KEY=$(jq -r .keys[0] init.json)
export VAULT_TOKEN=$(jq -r .root_token init.json)


http POST $VAULT_ADDR/v1/sys/unseal key=$UNSEAL_KEY


# Enable userpass auth method and create user
http POST $VAULT_ADDR/v1/sys/auth/userpass X-Vault-Token:$VAULT_TOKEN type=userpass
http POST $VAULT_ADDR/v1/auth/userpass/users/joe password=joe X-Vault-Token:$VAULT_TOKEN


http POST $VAULT_ADDR/v1/sys/auth/userpass X-Vault-Token:$VAULT_TOKEN type=userpass config:='{"token_type": "batch"}'


# Login with userpass
http POST $VAULT_ADDR/v1/auth/userpass/login/joe password=joe

# Create masses of tokens
siege --concurrent=100 "$VAULT_ADDR/v1/auth/userpass/login/joe POST password=joe"





# List token accessors
http LIST $VAULT_ADDR/v1/auth/token/accessors X-Vault-Token:$VAULT_TOKEN
http LIST $VAULT_ADDR/v1/auth/token/accessors X-Vault-Token:$VAULT_TOKEN | jq -r .data.keys[] | wc -l

# Try to revoke all tokens
http LIST $VAULT_ADDR/v1/auth/token/accessors X-Vault-Token:$VAULT_TOKEN | jq -r .data.keys[] | xargs -I{} http -v --ignore-stdin $VAULT_ADDR/v1/auth/token/lookup-accessor X-Vault-Token:$VAULT_TOKEN accessor={}


# token tidy
http POST $VAULT_ADDR/v1/auth/token/tidy X-Vault-Token:$VAULT_TOKEN

# fetch metrics
http $VAULT_ADDR/v1/sys/metrics?format=prometheus X-Vault-Token:$VAULT_TOKEN|grep -E "^(vault_token_create_count|vault_expire_num_leases|vault_expire_num_irrevocable_leases)"



watch -n 1 "for c in vault-etcd0-1 vault-etcd1-1 vault-etcd2-1; do docker exec -it \$c df -h /etcd-data; done"



# Check etcd data

docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl get --keys-only --prefix ""
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl endpoint status --write-out=table

docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl get --keys-only --prefix /vault/sys/token/accessor | wc -l
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl get --keys-only --prefix /vault/sys/token/id       | wc -l
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl get --keys-only --prefix /vault/sys/expire         | wc -l


# compact and defragment etcd

rev=$(docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl endpoint status --write-out=json | jq '.[] | (.. | .revision? | select(.!=null))')
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl compact $rev
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl defrag --cluster


# delete all token accessors directly from etcd
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl del --prefix=true /vault/sys/token/accessor
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl del --prefix=true /vault/sys/token/id/



/vault/sys/expire/id/auth/userpass/login/joe/h41be6d8c3aeeda4c799367915ccb7062fb455e2c4b492bef62e07b82f916d34f
/vault/sys/token/accessor/4d737c4ec12999a01a418549c1941b8444670fe1
/vault/sys/token/accessor/911ebc80c838d8c86b6937b2e90290d9c11707c6
/vault/sys/token/id/h41be6d8c3aeeda4c799367915ccb7062fb455e2c4b492bef62e07b82f916d34f
/vault/sys/token/id/h7b0991ca73c9bd6ca73f0c5cb4cc294394fac04df61ed9e876c86f426f684ac1


ETCD_QUOTA_BACKEND_BYTES=5000000 docker-compose up --force-recreate --detach etcd0 etcd1 etcd2
ETCD_QUOTA_BACKEND_BYTES=0       docker-compose up --force-recreate --detach etcd0 etcd1 etcd2

# disarm etcd alarms
docker exec -e ETCDCTL_API=3 vault-etcd0-1 etcdctl alarm disarm

docker-compose rm --stop --force etcd0 etcd1 etcd2
sudo rm -rf tmp/



vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:32:19.979Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:32:19.979Z [ERROR] token: failed to mark token as revoked
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:32:19.979Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:32:19.980Z [ERROR] token: failed to mark token as revoked
vault-vault-1  | 2024-02-16T17:32:19.980Z [ERROR] expiration: failed to revoke lease: lease_id=auth/userpass/login/joe/h00e5cb6c28265a7d321e6a15cf6226eb64fd194edafe572e7b77f9d47a7fd1b4 error="failed to revoke token: failed to revoke entry: failed to persist entry: etcdserver: mvcc: database space exceeded" attempts=1 next_attempt=13.378555207s
vault-vault-1  | 2024-02-16T17:32:19.980Z [ERROR] expiration: failed to revoke lease: lease_id=auth/userpass/login/joe/hf0be94a792e615ca3d45c774b43d9949f08c92b663148a9488b75be93692ec78 error="failed to revoke token: failed to revoke entry: failed to persist entry: etcdserver: mvcc: database space exceeded" attempts=1 next_attempt=20.621594164s
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:32:19.987Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:32:19.987Z [ERROR] token: failed to mark token as revoked
vault-vault-1  | 2024-02-16T17:32:19.987Z [ERROR] expiration: failed to revoke lease: lease_id=auth/userpass/login/joe/h6a6596fc1a63fc842b0bfe89a42d6dec9e16100bf43a44335901d8e25602cba1 error="failed to revoke token: failed to revoke entry: failed to persist entry: etcdserver: mvcc: database space exceeded" attempts=1 next_attempt=26.674189063s
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:32:19.992Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:32:19.992Z [ERROR] token: failed to mark token as revoked
vault-vault-1  | 2024-02-16T17:32:19.992Z [ERROR] expiration: failed to revoke lease: lease_id=auth/userpass/login/joe/h28ccf86f9541dceb8c5867b996517e26038157f0cb489bc5f107a38c963cb1e7 error="failed to revoke token: failed to revoke entry: failed to persist entry: etcdserver: mvcc: database space exceeded" attempts=1 next_attempt=28.849546981s
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:32:19.992Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}


vault-vault-1  | 2024-02-16T17:32:20.054Z [ERROR] expiration: failed to revoke lease: lease_id=auth/userpass/login/joe/hc80d67b929ba2974520b91b4a4db7c613cb29a43c7f31997cbfab8aa0d681dbc error="failed to revoke token: failed to revoke entry: failed to persist entry: etcdserver: mvcc: database space exceeded" attempts=1 next_attempt=10.193121426s
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:32:20.056Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:32:20.056Z [ERROR] token: failed to mark token as revoked
vault-vault-1  | 2024-02-16T17:32:30.258Z [INFO]  expiration: revoked lease: lease_id=auth/userpass/login/joe/hc80d67b929ba2974520b91b4a4db7c613cb29a43c7f31997cbfab8aa0d681dbc




vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:36:13.365Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:36:13.366Z [ERROR] core: error in barrier auto rotation: error="failed to persist keyring: etcdserver: mvcc: database space exceeded"
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:41:13.375Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}
vault-vault-1  | 2024-02-16T17:41:13.375Z [ERROR] core: error in barrier auto rotation: error="failed to persist keyring: etcdserver: mvcc: database space exceeded"
vault-vault-1  | {"level":"warn","ts":"2024-02-16T17:46:13.366Z","logger":"etcd-client","caller":"v3@v3.5.7/retry_interceptor.go:62","msg":"retrying of unary invoker failed","target":"etcd-endpoints://0xc000556000/etcd0:2379","attempt":0,"error":"rpc error: code = ResourceExhausted desc = etcdserver: mvcc: database space exceeded"}



Successful deletion path

etcd.(*EtcdBackend).Delete (/workspace/physical/etcd/etcd3.go:229)
vault.(*sealUnwrapper).Delete (/workspace/vault/sealunwrapper.go:149)
physical.(*Cache).Delete (/workspace/sdk/physical/cache.go:214)
physical.(*StorageEncoding).Delete (/workspace/sdk/physical/encoding.go:83)
vault.(*AESGCMBarrier).Delete (/workspace/vault/barrier_aes_gcm.go:922)
vault.(*ForwardedWriter).Delete (/workspace/vault/forwarded_writer_oss.go:92)
logical.(*StorageView).Delete (/workspace/sdk/logical/storage_view.go:84)
vault.(*BarrierView).Delete (/workspace/vault/barrier_view.go:93)
vault.(*TokenStore).revokeInternal.func1 (/workspace/vault/token_store.go:1819) <<<<<<<----------
runtime.deferreturn (/usr/local/go/src/runtime/panic.go:602)
vault.(*TokenStore).revokeInternal (/workspace/vault/token_store.go:1959)
vault.(*TokenStore).revokeTreeInternal (/workspace/vault/token_store.go:2063)
vault.(*TokenStore).revokeTree (/workspace/vault/token_store.go:1983)
vault.(*ExpirationManager).revokeEntry (/workspace/vault/expiration.go:1946)
vault.(*ExpirationManager).revokeCommon (/workspace/vault/expiration.go:1017)
vault.(*ExpirationManager).Revoke (/workspace/vault/expiration.go:919)
vault.(*TokenStore).lookupInternal (/workspace/vault/token_store.go:1730)
vault.(*TokenStore).handleTidy.func1.1 (/workspace/vault/token_store.go:2336)
vault.(*TokenStore).handleTidy.func1 (/workspace/vault/token_store.go:2432)
runtime.goexit (/usr/local/go/src/runtime/asm_amd64.s:1695)





#################
#
# To be checked
#


https://etcd.io/docs/v3.5/op-guide/maintenance/

https://developer.hashicorp.com/vault/docs/concepts/lease

https://developer.hashicorp.com/vault/api-docs/system/leases
https://developer.hashicorp.com/vault/api-docs/system/lease-count-quotas

https://developer.hashicorp.com/vault/docs/concepts/policies

https://developer.hashicorp.com/vault/docs/internals/limits





### devcontainer

cd ~/work/vault
mkdir -p .devcontainer .vscode
cp ~/work/devenvs/vault/configs/devcontainer.json .devcontainer


mkdir -p .vscode
cp ~/work/devenvs/vault/configs/launch.json .vscode














####

diff --git a/physical/etcd/etcd3.go b/physical/etcd/etcd3.go
index 8501a8b40f..00212dee21 100644
--- a/physical/etcd/etcd3.go
+++ b/physical/etcd/etcd3.go
@@ -186,6 +186,7 @@ func (c *EtcdBackend) Put(ctx context.Context, entry *physical.Entry) error {
        ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
        defer cancel()
        _, err := c.etcd.Put(ctx, path.Join(c.path, entry.Key), string(entry.Value))
+       fmt.Printf("Put (res: %t): %s\n", err == nil, entry.Key)
        return err
 }

@@ -198,6 +199,7 @@ func (c *EtcdBackend) Get(ctx context.Context, key string) (*physical.Entry, err
        ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
        defer cancel()
        resp, err := c.etcd.Get(ctx, path.Join(c.path, key))
+       fmt.Printf("Get (res: %t): %s\n", err == nil, key)
        if err != nil {
                return nil, err
        }
@@ -223,6 +225,7 @@ func (c *EtcdBackend) Delete(ctx context.Context, key string) error {
        ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
        defer cancel()
        _, err := c.etcd.Delete(ctx, path.Join(c.path, key))
+       fmt.Printf("Delete (res: %t): %s\n", err == nil, key)
        if err != nil {
                return err
        }
@@ -239,6 +242,7 @@ func (c *EtcdBackend) List(ctx context.Context, prefix string) ([]string, error)
        defer cancel()
        prefix = path.Join(c.path, prefix) + "/"
        resp, err := c.etcd.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithKeysOnly())
+       fmt.Printf("List prefix (res: %t): %s\n", err == nil, prefix)
        if err != nil {
                return nil, err
        }
~

