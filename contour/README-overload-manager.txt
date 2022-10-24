

docker run --volume $PWD/configs:/configs:ro --network host --rm envoyproxy/envoy-distroless:v1.22-latest --config-path /configs/envoy-overload-manager.yaml --log-level debug
python3 -m http.server 8081

http http://localhost:8080/bigfile
http http://localhost:8080/smallfile


make container
docker tag ghcr.io/projectcontour/contour:$(git rev-parse --short HEAD) localhost/contour:latest
kind load docker-image --name contour localhost/contour:latest

cat <<EOF | kubectl -n projectcontour patch daemonset envoy --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: shutdown-manager
        image: localhost/contour:latest
        imagePullPolicy: Never
      initContainers:
      - name: envoy-initconfig
        image: localhost/contour:latest
        imagePullPolicy: Never
EOF

kubectl -n projectcontour delete pod -l app=envoy --force


kubectl -n projectcontour logs daemonset/envoy envoy -f

sudo cat /proc/$(pidof envoy)/root/config/envoy.json | jq .


kubectl apply -f manifests/shell.yaml


kubectl exec shell -- http http://$(kubectl -n projectcontour get pod -lapp=envoy -o=jsonpath={..podIP}):8002/ready

kubectl exec shell -- http http://$(kubectl -n projectcontour get pod -lapp=envoy -o=jsonpath={..podIP}):8002/stats | grep -E "^overload|^server.memory"
sudo curl --silent --unix-socket /proc/$(pidof envoy)/root/admin/admin.sock http://localhost/stats | grep -E "^overload|^server.memory"

sudo curl --silent --unix-socket /proc/$(pidof envoy)/root/admin/admin.sock http://localhost/memory


watch --interval 1 kubectl -n projectcontour describe service envoy



Envoy uses https://github.com/google/tcmalloc as its malloc() implementation to manage heap memory.

tcmalloc exposes some statistics that are available at
$ sudo curl --silent --unix-socket /proc/$(pidof envoy)/root/admin/admin.sock http://localhost/memory
{
 "allocated": "8811712",
 "heap_size": "16777216",
 "pageheap_unmapped": "3358720",
 "pageheap_free": "8192",
 "total_thread_cache": "129760",
 "total_physical_bytes": "21089194"
}

https://github.com/envoyproxy/envoy/blob/43cd7847892892a0be3f01e2c7e8189bfb0a6321/source/server/admin/server_info_handler.cc#L48-L53
How these map to tcmalloc metrics and how are they documented?

The tcmalloc metrics are documented here
https://github.com/google/tcmalloc/blob/e091b551403f42417cb374110467a2f0ad661dc2/tcmalloc/malloc_extension.h#L225-L267
https://github.com/google/tcmalloc/blob/ac7a54fc916587a680b99ec7ceef1fd0e5b0953e/tcmalloc/malloc_extension.h#L427-L444

allocated: generic.current_allocated_bytes
  // "generic.current_allocated_bytes"
  //      Number of bytes currently allocated by application

heap_size: generic.heap_size + tcmalloc.pageheap_unmapped_bytes
  // "generic.heap_size"
  //      Number of bytes in the heap ==
  //            current_allocated_bytes +
  //            fragmentation +
  //            freed (but not released to OS) memory regions
  // "tcmalloc.pageheap_unmapped_bytes"
  //      Number of bytes in free, unmapped pages in page heap.
  //      These are bytes that have been released back to the OS,
  //      possibly by one of the MallocExtension "Release" calls.
  //      They can be used to fulfill allocation requests, but
  //      typically incur a page fault.  They always count towards
  //      virtual memory usage, and depending on the OS, typically
  //      do not count towards physical memory usage.

pageheap_unmapped: tcmalloc.pageheap_unmapped_bytes
  // "tcmalloc.pageheap_unmapped_bytes"
  //      Number of bytes in free, unmapped pages in page heap.
  //      These are bytes that have been released back to the OS,
  //      possibly by one of the MallocExtension "Release" calls.
  //      They can be used to fulfill allocation requests, but
  //      typically incur a page fault.  They always count towards
  //      virtual memory usage, and depending on the OS, typically
  //      do not count towards physical memory usage.

pageheap_free: tcmalloc.pageheap_free_bytes
  // "tcmalloc.pageheap_free_bytes"
  //      Number of bytes in free, mapped pages in page heap.  These
  //      bytes can be used to fulfill allocation requests.  They
  //      always count towards virtual memory usage, and unless the
  //      underlying memory is swapped out by the OS, they also count
  //      towards physical memory usage.

total_thread_cache   -> tcmalloc.current_total_thread_cache_bytes
  // "tcmalloc.current_total_thread_cache_bytes"
  //      Number of bytes used across all thread caches.

total_physical_bytes: generic.physical_memory_used
  // Overall (including malloc internals)


Envoy calculates pressure as follows
https://github.com/envoyproxy/envoy/blob/4b5eee6b8ec3fc84fd1f7e6ee684947724843e47/source/extensions/resource_monitors/fixed_heap/fixed_heap_monitor.cc#L24-L34


physical = generic.heap_size
unmapped = tcmalloc.pageheap_unmapped_bytes

used = physical - unmapped
pressure = used / configured_max_heap_size_bytes
