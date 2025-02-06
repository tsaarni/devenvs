





#############
#
# Run standalone vault raft node with fuse filesystem that generates delays
#



bin/vault server -config=$HOME/work/devenvs/vault/configs/vault-config-raft.hcl -log-level=debug

export VAULT_ADDR=http://127.0.0.1:8200

vault operator init -key-shares=1 -key-threshold=1 -format=json > vault-unseal-config.json
vault operator unseal $(jq -r .unseal_keys_b64[0] vault-unseal-config.json)

vault login $(jq -r .root_token vault-unseal-config.json)

vault secrets enable -path=secret kv-v2

vault kv put secret/hello foo="Hello at: $(date)"
vault kv get secret/hello

touch source/.block_operations







########
#
# chaos mesh
#

helm repo add chaos-mesh https://charts.chaos-mesh.org
helm search repo chaos-mesh

kubectl delete ns chaos-mesh
kubectl create ns chaos-mesh
helm install chaos-mesh chaos-mesh/chaos-mesh -n=chaos-mesh --set chaosDaemon.runtime=containerd --set chaosDaemon.socketPath=/run/containerd/containerd.sock --set controllerManager.replicaCount=1 --set dashboard.securityMode=false
￼
# check that pods are running￼
kubectl -n chaos-mesh get pod
￼




########
# Expose chaos-dashboard

# Install  contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# Create httpproxy for chaos-dashboard
cat <<EOF | kubectl apply -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: chaos-dashboard
  namespace: chaos-mesh
spec:
  virtualhost:
    fqdn: chaos-dashboard.127.0.0.195.nip.io
  routes:
    - services:
      - name: chaos-dashboard
        port: 2333
EOF






# test vault by creating a secret and reading it
kubectl exec -it vault-0 -- vault login $(jq -r .root_token vault-unseal-config.json)

# create kv secret engine
kubectl exec -it vault-0 -- vault secrets enable -path=secret kv-v2


# list mounted secret engines
kubectl exec -it vault-0 -- vault secrets list


kubectl exec -it vault-0 -- vault kv put secret/hello foo="Hello at: $(date)"
kubectl exec -it vault-0 -- vault kv get secret/hello

kubectl exec -it vault-0 -- sh -c 'echo "Hello at: $(date)" > /vault/data/hello.txt'
kubectl exec -it vault-0 -- cat /vault/data/hello.txt



kubectl -n chaos-mesh delete iochaos io-latency-example

cat <<EOF | kubectl apply -f -
apiVersion: chaos-mesh.org/v1alpha1
kind: IOChaos
metadata:
  name: io-latency-example
  namespace: chaos-mesh
spec:
  action: latency
  mode: one
  selector:
    labelSelectors:
      statefulset.kubernetes.io/pod-name: vault-0
    namespaces:
      - default
  volumePath: /vault/data
  path: '/vault/data/**/*'
  delay: '60s'
  percent: 100
  duration: '60s'
EOF



###########################
#
# output from chaos-mesh
#

unexpected fault address 0x73dd70601040
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x1 addr=0x73dd70601040 pc=0x1f2896a]

goroutine 989 gp=0xc003d01500 m=15 mp=0xc0048a2008 [running]:
runtime.throw({0xb864f99?, 0xc001dbe000?})
        /opt/hostedtoolcache/go/1.22.8/x64/src/runtime/panic.go:1023 +0x5c fp=0xc003951020 sp=0xc003950ff0 pc=0x43e99c
runtime.sigpanic()
        /opt/hostedtoolcache/go/1.22.8/x64/src/runtime/signal_unix.go:895 +0x285 fp=0xc003951080 sp=0xc003951020 pc=0x4575c5
go.etcd.io/bbolt.(*DB).meta(0xc0048a2008?)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/db.go:1094 +0x2a fp=0xc0039510c8 sp=0xc003951080 pc=0x1f2896a
go.etcd.io/bbolt.(*Tx).init(0xc001dbe000, 0xc001c266c8)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/tx.go:51 +0x8d fp=0xc003951128 sp=0xc0039510c8 pc=0x1f3066d
go.etcd.io/bbolt.(*DB).beginRWTx(0xc001c266c8)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/db.go:800 +0x170 fp=0xc003951198 sp=0xc003951128 pc=0x1f273d0
go.etcd.io/bbolt.(*DB).Begin(0xc0039511e8?, 0xa5?)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/db.go:721 +0x17 fp=0xc0039511b0 sp=0xc003951198 pc=0x1f26f77
github.com/hashicorp/raft-boltdb/v2.(*BoltStore).StoreLogs(0xc001c0d2e0, {0xc003598070, 0x1, 0x1})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft-boltdb/v2@v2.3.0/bolt_store.go:191 +0x7d fp=0xc0039512b0 sp=0xc0039511b0 pc=0x1f3b13d
github.com/hashicorp/raft.(*LogCache).StoreLogs(0xc003a36000, {0xc003598070, 0x1, 0xc003ba2540?})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/log_cache.go:67 +0x31 fp=0xc003951300 sp=0xc0039512b0 pc=0x1ef32b1
github.com/hashicorp/raft.(*Raft).dispatchLogs(0xc0001fc008, {0xc003951a90, 0x1, 0x0?})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/raft.go:1266 +0x3c5 fp=0xc003951450 sp=0xc003951300 pc=0x1f03fc5
github.com/hashicorp/raft.(*Raft).leaderLoop(0xc0001fc008)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/raft.go:931 +0x16bd fp=0xc003951d20 sp=0xc003951450 pc=0x1f0039d
github.com/hashicorp/raft.(*Raft).runLeader(0xc0001fc008)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/raft.go:575 +0x48c fp=0xc003951f60 sp=0xc003951d20 pc=0x1efdb8c
github.com/hashicorp/raft.(*Raft).run(0xc0001fc008)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/raft.go:152 +0x4a fp=0xc003951f98 sp=0xc003951f60 pc=0x1ef9baa
github.com/hashicorp/raft.(*Raft).run-fm()
        <autogenerated>:1 +0x25 fp=0xc003951fb0 sp=0xc003951f98 pc=0x1f16f65
github.com/hashicorp/raft.(*raftState).goFunc.func1()
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/state.go:149 +0x56 fp=0xc003951fe0 sp=0xc003951fb0 pc=0x1f137d6
runtime.goexit({})
        /opt/hostedtoolcache/go/1.22.8/x64/src/runtime/asm_amd64.s:1695 +0x1 fp=0xc003951fe8 sp=0xc003951fe0 pc=0x479d61
created by github.com/hashicorp/raft.(*raftState).goFunc in goroutine 984
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/state.go:147 +0x79






unexpected fault address 0x7b8bd7e01040
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x1 addr=0x7b8bd7e01040 pc=0x1f2896a]

goroutine 188 gp=0xc001d1f340 m=3 mp=0xc000180008 [running]:
runtime.throw({0xb864f99?, 0xc002dcdce0?})
        /opt/hostedtoolcache/go/1.22.8/x64/src/runtime/panic.go:1023 +0x5c fp=0xc003a8a9b0 sp=0xc003a8a980 pc=0x43e99c
runtime.sigpanic()
        /opt/hostedtoolcache/go/1.22.8/x64/src/runtime/signal_unix.go:895 +0x285 fp=0xc003a8aa10 sp=0xc003a8a9b0 pc=0x4575c5
go.etcd.io/bbolt.(*DB).meta(0xc000180008?)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/db.go:1094 +0x2a fp=0xc003a8aa58 sp=0xc003a8aa10 pc=0x1f2896a
go.etcd.io/bbolt.(*Tx).init(0xc002dcdce0, 0xc00390c6c8)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/tx.go:51 +0x8d fp=0xc003a8aab8 sp=0xc003a8aa58 pc=0x1f3066d
go.etcd.io/bbolt.(*DB).beginTx(0xc00390c6c8)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/db.go:753 +0x10b fp=0xc003a8ab28 sp=0xc003a8aab8 pc=0x1f270cb
go.etcd.io/bbolt.(*DB).Begin(0xc003a8ab80?, 0xb4?)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.3.10/db.go:723 +0x25 fp=0xc003a8ab40 sp=0xc003a8ab28 pc=0x1f26f85
github.com/hashicorp/raft-boltdb/v2.(*BoltStore).LastIndex(0x10?)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft-boltdb/v2@v2.3.0/bolt_store.go:149 +0x37 fp=0xc003a8abc8 sp=0xc003a8ab40 pc=0x1f3ab17
github.com/hashicorp/raft.(*LogCache).LastIndex(0x10?)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/log_cache.go:85 +0x1b fp=0xc003a8abe0 sp=0xc003a8abc8 pc=0x1ef341b
github.com/hashicorp/raft.(*Raft).shouldSnapshot(0xc003aba008)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/snapshot.go:111 +0x50 fp=0xc003a8ad88 sp=0xc003a8abe0 pc=0x1f119f0
github.com/hashicorp/raft.(*Raft).runSnapshots(0xc003aba008)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/snapshot.go:77 +0x3bd fp=0xc003a8af98 sp=0xc003a8ad88 pc=0x1f1181d
github.com/hashicorp/raft.(*Raft).runSnapshots-fm()
        <autogenerated>:1 +0x25 fp=0xc003a8afb0 sp=0xc003a8af98 pc=0x1f17025
github.com/hashicorp/raft.(*raftState).goFunc.func1()
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/state.go:149 +0x56 fp=0xc003a8afe0 sp=0xc003a8afb0 pc=0x1f137d6
runtime.goexit({})
        /opt/hostedtoolcache/go/1.22.8/x64/src/runtime/asm_amd64.s:1695 +0x1 fp=0xc003a8afe8 sp=0xc003a8afe0 pc=0x479d61
created by github.com/hashicorp/raft.(*raftState).goFunc in goroutine 154
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/state.go:147 +0x79



goroutine 785 gp=0xc003904e00 m=7 mp=0xc001d80008 [running]:
runtime.throw({0xc1a5e33?, 0x1008ca913a8?})
        /opt/hostedtoolcache/go/1.23.3/x64/src/runtime/panic.go:1067 +0x48 fp=0xc003944b08 sp=0xc003944ad8 pc=0x475a08
runtime.sigpanic()
        /opt/hostedtoolcache/go/1.23.3/x64/src/runtime/signal_unix.go:931 +0x26c fp=0xc003944b68 sp=0xc003944b08 pc=0x477dec
go.etcd.io/bbolt/internal/common.(*Meta).Txid(...)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/internal/common/meta.go:128
go.etcd.io/bbolt.(*DB).meta(0xc001d80008?)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/db.go:1129 +0x2a fp=0xc003944bb0 sp=0xc003944b68 pc=0x1f680ca
go.etcd.io/bbolt.(*Tx).init(0xc0041061c0, 0xc0001a2248)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/tx.go:53 +0x8d fp=0xc003944c10 sp=0xc003944bb0 pc=0x1f6caed
go.etcd.io/bbolt.(*DB).beginTx(0xc0001a2248)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/db.go:795 +0x10b fp=0xc003944c80 sp=0xc003944c10 pc=0x1f666cb
go.etcd.io/bbolt.(*DB).Begin(0xc0001a2248, 0x0)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/db.go:758 +0x1cf fp=0xc003944d20 sp=0xc003944c80 pc=0x1f663af
github.com/hashicorp/raft-boltdb/v2.(*BoltStore).FirstIndex(0xc002d76301?)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft-boltdb/v2@v2.3.0/bolt_store.go:133 +0x37 fp=0xc003944da8 sp=0xc003944d20 pc=0x1f78eb7
github.com/hashicorp/raft.(*LogCache).FirstIndex(0xc005258720?)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/log_cache.go:81 +0x1b fp=0xc003944dc0 sp=0xc003944da8 pc=0x1f247fb
github.com/hashicorp/raft.oldestLog({0xddd5b40, 0xc00019b0c0})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/log.go:153 +0x84 fp=0xc003944e10 sp=0xc003944dc0 pc=0x1f23f04
github.com/hashicorp/raft.emitLogStoreMetrics({0xddd5b40, 0xc00019b0c0}, {0xc00341f2a0, 0x2, 0x2}, 0x2540be400, 0xc002d993b0)
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/log.go:183 +0x108 fp=0xc003944f80 sp=0xc003944e10 pc=0x1f24148
github.com/hashicorp/raft.(*Raft).runLeader.gowrap1()
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/raft.go:499 +0x94 fp=0xc003944fe0 sp=0xc003944f80 pc=0x1f2f474
runtime.goexit({})
        /opt/hostedtoolcache/go/1.23.3/x64/src/runtime/asm_amd64.s:1700 +0x1 fp=0xc003944fe8 sp=0xc003944fe0 pc=0x47e4a1
created by github.com/hashicorp/raft.(*Raft).runLeader in goroutine 236
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/raft.go:499 +0x3d3









2025-02-03T18:25:56.065Z [WARN]  storage.raft: heartbeat timeout reached, starting election: last-leader-addr=vault-1.vault-internal:8201 last-leader-id=vault-1
2025-02-03T18:25:56.065Z [INFO]  storage.raft: entering candidate state: node="Node at vault-0.vault-internal:8201 [Candidate]" term=5
unexpected fault address 0x717a8d401040
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x1 addr=0x717a8d401040 pc=0x1f680ca]

goroutine 585 gp=0xc00333f6c0 m=19 mp=0xc004b80008 [running]:
runtime.throw({0xc1a5e33?, 0x1000089292e?})
        /opt/hostedtoolcache/go/1.23.3/x64/src/runtime/panic.go:1067 +0x48 fp=0xc004837b80 sp=0xc004837b50 pc=0x475a08
runtime.sigpanic()
        /opt/hostedtoolcache/go/1.23.3/x64/src/runtime/signal_unix.go:931 +0x26c fp=0xc004837be0 sp=0xc004837b80 pc=0x477dec
go.etcd.io/bbolt/internal/common.(*Meta).Txid(...)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/internal/common/meta.go:128
go.etcd.io/bbolt.(*DB).meta(0xc004b80008?)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/db.go:1129 +0x2a fp=0xc004837c28 sp=0xc004837be0 pc=0x1f680ca
go.etcd.io/bbolt.(*Tx).init(0xc002e6a540, 0xc003b8c6c8)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/tx.go:53 +0x8d fp=0xc004837c88 sp=0xc004837c28 pc=0x1f6caed
go.etcd.io/bbolt.(*DB).beginTx(0xc003b8c6c8)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/db.go:795 +0x10b fp=0xc004837cf8 sp=0xc004837c88 pc=0x1f666cb
go.etcd.io/bbolt.(*DB).Begin(0xc003b8c6c8, 0x0)
        /home/runner/go/pkg/mod/go.etcd.io/bbolt@v1.4.0-beta.0/db.go:758 +0x1cf fp=0xc004837d98 sp=0xc004837cf8 pc=0x1f663af
github.com/hashicorp/raft-boltdb/v2.(*BoltStore).Get(0xe?, {0x1469a8f0, 0xb, 0xb})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft-boltdb/v2@v2.3.0/bolt_store.go:274 +0x65 fp=0xc004837e40 sp=0xc004837d98 pc=0x1f7a205
github.com/hashicorp/raft-boltdb/v2.(*BoltStore).GetUint64(0x4?, {0x1469a8f0?, 0x0?, 0x0?})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft-boltdb/v2@v2.3.0/bolt_store.go:296 +0x18 fp=0xc004837e70 sp=0xc004837e40 pc=0x1f7a578
github.com/hashicorp/raft.HasExistingState({0xddd5b40, 0xc003c26100}, {0xddc3028?, 0xc003c24300?}, {0xdda4430, 0xc003b95920})
        /home/runner/go/pkg/mod/github.com/hashicorp/raft@v1.7.1/api.go:456 +0x5b fp=0xc004837ec8 sp=0xc004837e70 pc=0x1f181fb
github.com/hashicorp/vault/physical/raft.(*RaftBackend).HasState(0xc003aad008?)
        /home/runner/work/vault/vault/physical/raft/raft.go:1107 +0xa5 fp=0xc004837f38 sp=0xc004837ec8 pc=0x1fa6d65
github.com/hashicorp/vault/vault.(*Core).periodicCheckKeyUpgrades.func1()
        /home/runner/work/vault/vault/vault/ha.go:977 +0x2b5 fp=0xc004837fe0 sp=0xc004837f38 pc=0x3db1a55
runtime.goexit({})
        /opt/hostedtoolcache/go/1.23.3/x64/src/runtime/asm_amd64.s:1700 +0x1 fp=0xc004837fe8 sp=0xc004837fe0 pc=0x47e4a1
created by github.com/hashicorp/vault/vault.(*Core).periodicCheckKeyUpgrades in goroutine 123
        /home/runner/work/vault/vault/vault/ha.go:944 +0x1ab
