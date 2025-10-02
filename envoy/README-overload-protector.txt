

Envoy is not freeing memory in a meaningful way to the system
https://github.com/envoyproxy/envoy/issues/39006

memory heap not released long after overload manager actions triggered
https://github.com/envoyproxy/envoy/issues/21923

CGROUP aware resource monitor on memory
https://github.com/envoyproxy/envoy/issues/38718

Fixed heap monitor: discount pageheap free bytes
https://github.com/envoyproxy/envoy/pull/22585






docker compose -f docker-compose-overload.yaml rm -f
docker compose -f docker-compose-overload.yaml up

# Grafana dashboard
http://localhost:3000/




http http://localhost:9901/stats/prometheus

http http://localhost:9901/memory




while true; do http http://localhost:8080; sleep 1; done



dd if=/dev/zero bs=100M count=1 | http -v POST http://localhost:8080/upload?throttle=1K




kind delete cluster --name=envoy
kind create cluster --name=envoy --config=configs/kind-cluster-config.yaml



kubectl apply -f manifests/echoserver.yaml
kubectl apply -f manifests/deploy-envoy-overload-protector.yaml

kubectl apply -f manifests/observability-stack.yaml

http://127.0.0.135.nip.io:3000


http http://127.0.0.135.nip.io/



cd ~/work/echoclient
go run ./cmd/echoclient get -addr http://127.0.0.135.nip.io -concurrency 1000

go run ./cmd/echoclient upload -addr http://127.0.0.135.nip.io/upload?throttle=1K -chunksize 10M -totalsize 100M -repetitions 0 -concurrency 10

http http://127.0.0.135.nip.io


sudo bpftrace -e 'tracepoint:syscalls:sys_enter_madvise /args->behavior == 4/ { printf("madvise DONTNEED called\n"); }'
sudo bpftrace -e "tracepoint:syscalls:sys_enter_madvise /pid == $(pidof envoy) && args->behavior == 4/ { printf(\"madvise DONTNEED called\n\"); }"




grep "define MADV_" /usr/include/x86_64-linux-gnu/bits/mman-linux.h

sudo bpftrace - <<EOF
tracepoint:syscalls:sys_enter_madvise
/ pid == $(pidof envoy) /
{
  printf("%s pid=%d comm=%s start=0x%lx len=0x%lx behavior=%d\n",
         strftime("%H:%M:%S", nsecs), pid, comm, args->start, args->len_in, args->behavior);
}
EOF


sudo pmap -x $(pidof envoy)
sudo bash -c 'while true; do date; pmap -X $(pidof envoy) | grep tcmalloc; sleep 3; done'
