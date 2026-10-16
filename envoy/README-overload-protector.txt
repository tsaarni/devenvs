

Use allocated memory for fixed heap pressure
https://github.com/envoyproxy/envoy/pull/41442



Envoy is not freeing memory in a meaningful way to the system
https://github.com/envoyproxy/envoy/issues/39006

memory heap not released long after overload manager actions triggered
https://github.com/envoyproxy/envoy/issues/21923

CGROUP aware resource monitor on memory
https://github.com/envoyproxy/envoy/issues/38718

Fixed heap monitor: discount pageheap free bytes
https://github.com/envoyproxy/envoy/pull/22585




# https://google.github.io/tcmalloc/design.html
# https://github.com/google/tcmalloc/blob/master/docs/stats.md

# https://docs.kernel.org/admin-guide/cgroup-v2.html


kubectl port-forward envoy-0 9901:9901


http http://localhost:9901/stats/prometheus
http http://localhost:9901/memory
http http://localhost:9901/memory/tcmalloc



bazel build -c fastbuild //source/exe:envoy-static
cp -a bazel-bin/source/exe/envoy-static /home/tsaarni/envoy-tcmalloc

kubectl delete pod -l app=prometheus --force
kubectl delete pod envoy-0 --force


cd ~/work/devenvs/envoy$
mv ~/envoy-tcmalloc  .
kubectl exec -it envoy-0 -- /host/envoy-tcmalloc -c /host/configs/envoy-overload-manager.yaml --log-level info


####################
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


# Restart any cgroups stats and prometheus metrics
kubectl delete pod -l app=envoy --force
kubectl delete pod -l app=prometheus --force


# Grafana dashboard
http://127.0.0.135.nip.io:3000





# Traffic generation

cd ~/work/echoclient
go run ./cmd/echoclient get -url http://127.0.0.135.nip.io/ -concurrency 10000 -duration 30s

go run ./cmd/echoclient upload -url http://127.0.0.135.nip.io/upload?throttle=1K -chunk 10M -size 100M -repetitions 0 -concurrency 100 -repetitions 0 -duration 30s



# Madvise tracing


code /build/.cache/bazel/_bazel_vscode/2d35de14639eaad1ac7060a4dd7e3351/external/com_github_google_tcmalloc/tcmalloc/parameters.cc

ABSL_CONST_INIT std::atomic<MadvisePreference> Parameters::madvise_(
    MadvisePreference::kFreeAndDontNeed);
//    MadvisePreference::kDontNeed);

                absl::Seconds(60)
//                absl::Seconds(300)


grep "define MADV_" /usr/include/x86_64-linux-gnu/bits/mman-linux.h

sudo bpftrace - <<EOF
tracepoint:syscalls:sys_enter_madvise
/ pid == $(pidof envoy-tcmalloc) /
{
  printf("%s pid=%d comm=%s start=0x%lx len=0x%lx behavior=%d\n",
         strftime("%H:%M:%S", nsecs), pid, comm, args->start, args->len_in, args->behavior);
}
EOF


# Smaps parsing

sudo pmap -x $(pidof envoy)
sudo bash -c 'while true; do date; pmap -X $(pidof envoy) | grep tcmalloc; sleep 3; done'
sudo ~/wiki/files/bin/parse_smaps.py -p envoy





### https://docs.kernel.org/admin-guide/mm/transhuge.html


# defaults on ubuntu 24.04.3 LTS

cat /sys/kernel/mm/transparent_hugepage/enabled /sys/kernel/mm/transparent_hugepage/defrag /sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none

always [madvise] never
always defer defer+madvise [madvise] never
511

'/sys/kernel/mm/transparent_hugepage/enabled:': No such file or directory
/sys/kernel/mm/transparent_hugepage/defrag:: No such file or directory
'/sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none:': No such file or directory

cat /proc/sys/vm/overcommit_memory
1


sudo bash -c "echo madvise >/sys/kernel/mm/transparent_hugepage/enabled; echo madvise >/sys/kernel/mm/transparent_hugepage/defrag; echo 511 >/sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none"



sudo bash -c "echo never >/sys/kernel/mm/transparent_hugepage/enabled; echo never >/sys/kernel/mm/transparent_hugepage/defrag; echo 0 >/sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none"


# change according to
# https://google.github.io/tcmalloc/tuning.html

sudo bash -c "echo always >/sys/kernel/mm/transparent_hugepage/enabled; echo defer+madvise >/sys/kernel/mm/transparent_hugepage/defrag; echo 0 >/sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none"


grep -i huge /proc/meminfo
AnonHugePages:    704512 kB
ShmemHugePages:  2088960 kB
FileHugePages:         0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
Hugetlb:               0 kB

$ cat /sys/kernel/mm/transparent_hugepage/hpage_pmd_size
2097152
$ pidof khugepaged
148


Adaptive huge-page subrelease for non-moving memory allocators in warehouse-scale computers
https://dl.acm.org/doi/10.1145/3459898.3463905

Understanding Malloc Stats
https://google.github.io/tcmalloc/stats.html



sudo bash -c "echo never > /sys/kernel/mm/transparent_hugepage/defrag; echo never > /sys/kernel/mm/transparent_hugepage/enabled"


# reclaim manually

ENVOY_CGROUP_PATH=/sys/fs/cgroup/$(cat /proc/$(pidof envoy-tcmalloc)/cgroup | grep '0::' | cut -d: -f3)

cat $ENVOY_CGROUP_PATH/memory.current

sudo bash -c "echo 100M > $ENVOY_CGROUP_PATH/memory.reclaim"

cat $ENVOY_CGROUP_PATH/memory.current


https://github.com/facebookincubator/senpai
https://www.usenix.org/system/files/lisa21_slides_down.pdf


https://linuxera.org/cpu-memory-management-kubernetes-cgroupsv2/

https://serverfault.com/questions/1166165/how-do-memory-limits-in-kubernetes-work-with-cgroup-v2-memory-high
https://kubernetes.io/blog/2023/05/05/qos-memory-resources/
https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/

1. run envoy
2. run senpai on background

ENVOY_CGROUP_PATH=/sys/fs/cgroup/$(cat /proc/$(pidof envoy-tcmalloc)/cgroup | grep '0::' | cut -d: -f3)
sudo python3 senpai.py $ENVOY_CGROUP_PATH

3. trigger load

4. read what senpai set

cat $ENVOY_CGROUP_PATH/memory.high


Manually trigger memory pressure

sudo bash -c "echo 274308224 > $ENVOY_CGROUP_PATH/memory.high"



process_smaps_anonymous_bytes
process_smaps_anon_hugepages_bytes
process_smaps_rss_bytes


I'm curious what the motivation would be to proactively reclaim instead of letting the system handle it?  A few possible reasons I could imagine:

* Cgroup-aware memory monitor relies on memory actually being reclaimed to report accurate pressure values.
* Provide users with a more accurate view of actual memory usage.


Another consideration while using `memory.current` is that cgroup memory monitor is still affected by tcmalloc's 5-minute subrelease delay. If there is a large spikie in allocations, many allocated pages will "age out" and expire simultaneously, causing Envoy to report memory pressure for several minutes when there is no real pressure. I'm not sure whether this is real issue or just something exaggerated by my simple test case.



date > 1-tcmalloc-dump-at-start.txt
echo "" >> 1-tcmalloc-dump-at-start.txt
http http://localhost:9901/memory/tcmalloc >> 1-tcmalloc-dump-at-start.txt


date
go run ./cmd/echoclient get -url http://127.0.0.135.nip.io/ -concurrency 10000 -duration 30s


date > 2-tcmalloc-dump-during-load.txt
echo "" >> 2-tcmalloc-dump-during-load.txt
http http://localhost:9901/memory/tcmalloc >> 2-tcmalloc-dump-during-load.txt

date > 3-tcmalloc-dump-after-load.txt
echo "" >> 3-tcmalloc-dump-after-load.txt
http http://localhost:9901/memory/tcmalloc >> 3-tcmalloc-dump-after-load.txt

date > 4-tcmalloc-dump-5-minutes-after-load.txt
echo "" >> 4-tcmalloc-dump-5-minutes-after-load.txt
http http://localhost:9901/memory/tcmalloc >> 4-tcmalloc-dump-5-minutes-after-load.txt





bazel test -c fastbuild //test/extensions/resource_monitors/fixed_heap:fixed_heap_monitor_test


bazel test -c fastbuild //test/extensions/resource_monitors/fixed_heap:fixed_heap_monitor_test --define tcmalloc=gperftools
bazel test -c fastbuild //test/extensions/resource_monitors/fixed_heap:fixed_heap_monitor_test --define=tcmalloc=disabled



---------------------- | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `allocated`            | `generic.current_allocated_bytes`                        | Memory currently allocated by the heap for Envoy.                                                  |
| `heap_size`            | `generic.heap_size` + `tcmalloc.pageheap_unmapped_bytes` | Total size of heap (not necessarily used at the moment) including both mapped and unmapped memory. |
| `pageheap_unmapped`    | `tcmalloc.pageheap_unmapped_bytes`                       | Memory that is released back to the OS.                                                            |
| `pageheap_free`        | `tcmalloc.pageheap_free_bytes`                           | Memory that is free in the page heap.                                                              |
| `total_thread_cache`   | `tcmalloc.current_total_thread_cache_bytes`              | Memory held in thread-local caches.                                                                |
| `total_physical_bytes` | `generic.physical_memory_used`                           | Total physical memory used by the process.                                                         |




used_memory = generic.heap_size + tcmalloc.pageheap_unmapped_bytes - tcmalloc.pageheap_unmapped_bytes - tcmalloc.pageheap_free_bytes - tcmalloc.current_total_thread_cache_bytes



heap_size =  stats.system_bytes - stats.unmapped_bytes


// Return approximate number of bytes in use by app.
uint64_t InUseByApp(const TCMallocStats& stats) {
  return StatSub(stats.pageheap.system_bytes,
                 stats.thread_bytes + stats.central_bytes +
                     stats.transfer_bytes + stats.per_cpu_bytes +
                     stats.sharded_transfer_bytes + stats.pageheap.free_bytes +
                     stats.pageheap.unmapped_bytes);
}




totalPageHeapFree() + transferCacheFree() + centralCacheFree() + threadCacheFree() +   shardedTransferCacheFree();

tcmalloc.pageheap_free_bytes                  + tcmalloc.transfer_cache_free + tcmalloc.central_cache_free + tcmalloc.thread_cache_free + tcmalloc.sharded_transfer_cache_free

free_bytes (Total bytes on normal freelists)       stats.sharded_transfer_bytes    stats.central_bytes         stats.thread_bytes           stats.sharded_transfer_bytes
