

Envoy is not freeing memory in a meaningful way to the system
https://github.com/envoyproxy/envoy/issues/39006

memory heap not released long after overload manager actions triggered
https://github.com/envoyproxy/envoy/issues/21923

CGROUP aware resource monitor on memory
https://github.com/envoyproxy/envoy/issues/38718

Fixed heap monitor: discount pageheap free bytes
https://github.com/envoyproxy/envoy/pull/22585




kubectl port-forward envoy-0 9901:9901



http http://localhost:9901/stats/prometheus
http http://localhost:9901/memory

# https://github.com/google/tcmalloc/blob/master/docs/stats.md
http http://localhost:9901/memory/tcmalloc




###################3
#
# Running kind cluster with envoy overload protector
#


kind delete cluster --name=envoy
kind create cluster --name=envoy --config=configs/kind-cluster-config.yaml



kubectl apply -f manifests/echoserver.yaml
kubectl apply -f manifests/observability-stack.yaml

kubectl apply -f manifests/envoy-poller.yaml
kubectl logs envoy-poller -f


kubectl delete statefulsets.apps envoy --force


# Envoy official image
kubectl apply -f manifests/deploy-envoy-overload-protector.yaml


# Custom build

kubectl apply -f manifests/deploy-envoy-placeholder.yaml
kubectl exec -it envoy-0 -- bash

cp -a bazel-bin/source/exe/envoy-static /home/tsaarni/
mv ~/envoy-static ~/work/devenvs/envoy/envoy

kubectl exec -it envoy-0 -- /host/envoy-tcmalloc -c /host/configs/envoy-overload-manager.yaml --log-level info
kubectl exec -it envoy-0 -- /host/envoy-gperftools -c /host/configs/envoy-overload-manager.yaml --log-level info




# Grafana dashboard
http://127.0.0.135.nip.io:3000


# Reset collected metrics
kubectl delete pod -l app=prometheus --force



# Poll echoserver
http http://127.0.0.135.nip.io/
while true; do http http://127.0.0.135.nip.io; sleep 5; done


# Traffic generation

cd ~/work/echoclient
go run ./cmd/echoclient get -url http://127.0.0.135.nip.io -concurrency 1000 -duration 30s

go run ./cmd/echoclient upload -url http://127.0.0.135.nip.io/upload?throttle=1K -chunk 10M -size 100M -repetitions 0 -concurrency 100 -repetitions 0 -duration 30s



# Madvise tracing


grep "define MADV_" /usr/include/x86_64-linux-gnu/bits/mman-linux.h

sudo bpftrace - <<EOF
tracepoint:syscalls:sys_enter_madvise
/ pid == $(pidof envoy) /
{
  printf("%s pid=%d comm=%s start=0x%lx len=0x%lx behavior=%d\n",
         strftime("%H:%M:%S", nsecs), pid, comm, args->start, args->len_in, args->behavior);
}
EOF


# Smaps parsing

sudo pmap -x $(pidof envoy)
sudo bash -c 'while true; do date; pmap -X $(pidof envoy) | grep tcmalloc; sleep 3; done'
sudo ~/wiki/files/bin/parse_smaps.py -p envoy






# defaults on ubuntu 24.04.3 LTS

cat /sys/kernel/mm/transparent_hugepage/enabled /sys/kernel/mm/transparent_hugepage/defrag /sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none

'/sys/kernel/mm/transparent_hugepage/enabled:': No such file or directory
/sys/kernel/mm/transparent_hugepage/defrag:: No such file or directory
'/sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none:': No such file or directory

cat /proc/sys/vm/overcommit_memory
1

# change according to
# https://google.github.io/tcmalloc/tuning.html

sudo bash -c "echo always >/sys/kernel/mm/transparent_hugepage/enabled"
sudo bash -c "echo defer+madvise >/sys/kernel/mm/transparent_hugepage/defrag"
sudo bash -c "echo 0 >/sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none"
