# Envoy Overload Manager

The [Overload Manager](https://www.envoyproxy.io/docs/envoy/latest/configuration/operations/overload_manager/overload_manager) provides mechanism to protect Envoy instances from overload of various resources.
It monitors a configurable set of resources and notifies registered listeners when triggers related to those resources fire.

## TCMalloc

TCMalloc (Thread-Caching Malloc) is a heap memory allocator developed by Google.
It is used by Envoy to manage dynamic memory allocation.

There are two implementations of TCMalloc

- [Google's tcmalloc](https://github.com/google/tcmalloc/) which is default for x86_64 architecture.
- [gperftools](https://github.com/gperftools/gperftools) tcmalloc is enabled with `--define tcmalloc=gperftools` when building Envoy, or default for non-x86_64 architectures.

This document focuses on the Google TCMalloc implementation.

## Heap Statistics

Here's a description of the fields returned by Envoy's Admin interface [`/memory`](https://www.envoyproxy.io/docs/envoy/latest/operations/admin#get--memory) endpoint, mapping to [TCMalloc properties](https://github.com/google/tcmalloc/blob/12f255231938d30493186b0a037feedd70f5a1c1/tcmalloc/malloc_extension.h#L374-L416) and [descriptions](https://www.envoyproxy.io/docs/envoy/latest/api-v3/admin/v3/memory.proto.html) what they represent.

| `/memory` endpoint     | TCMalloc Property                                        | Description                                                                                        |
| ---------------------- | -------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `allocated`            | `generic.current_allocated_bytes`                        | Memory currently allocated by the heap for Envoy.                                                  |
| `heap_size`            | `generic.heap_size` + `tcmalloc.pageheap_unmapped_bytes` | Total size of heap (not necessarily used at the moment) including both mapped and unmapped memory. |
| `pageheap_unmapped`    | `tcmalloc.pageheap_unmapped_bytes`                       | Memory that is released back to the OS.                                                            |
| `pageheap_free`        | `tcmalloc.pageheap_free_bytes`                           | Memory that is free in the page heap.                                                              |
| `total_thread_cache`   | `tcmalloc.current_total_thread_cache_bytes`              | Memory held in thread-local caches.                                                                |
| `total_physical_bytes` | `generic.physical_memory_used`                           | Total physical memory used by the process.                                                         |

Following table summarizes how fields in `/memory` elate to `/stats/prometheus` metrics:

| `/memory` endpoint     | `/stats/prometheus` metric          |
| ---------------------- | ----------------------------------- |
| `allocated`            | `envoy_server_memory_allocated`     |
| `heap_size`            | `envoy_server_memory_heap_size`     |
| `pageheap_unmapped`    |                                     |
| `pageheap_free`        |                                     |
| `total_thread_cache`   |                                     |
| `total_physical_bytes` | `envoy_server_memory_physical_size` |

As it can be seen, not all `/memory` fields have corresponding Prometheus metrics.
Due to this issue, the $\text{pressure}$ metric cannot be derived from Prometheus metrics alone.

> TODO: Consider adding missing metrics to metrics.

## Fixed Heap Resource Monitor

Envoy's Fixed Heap Resource Monitor is [one of the resource monitors](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/resource_monitor/resource_monitor#v3-config-resource-monitors) provided by Envoy.
It uses the heap statistics above to calculate memory pressure.
The memory pressure is determined using the following formula, based on the heap statistics described above:

```math
\begin{align*}
    \text{used\_memory} &= \text{heap\_size} - \text{pageheap\_unmapped} - \text{pageheap\_free} \\
    \text{pressure} &= \frac{\text{used\_memory}}{\text{max\_heap\_size}} \\
\end{align*}
```

This equation ensures that only actively used memory (excluding unmapped and free pages) counts towards the heap pressure calculation. The pressure value is used to determine if the system is approaching memory exhaustion, where:

- $\text{pressure} < 1.0$ (less than 100% memory usage) means the system is operating within its memory limits.
- $\text{pressure} \geq 1.0$ (100% memory usage or higher) means the system is exceeding its configured memory limit.


### Example

You can query the `/memory` endpoint of Envoy's admin interface to see these values. For example:

```bash
curl -s $ENVOY_ADDR:9901/memory | jq
```

```json
{
    "allocated": "5788928",
    "heap_size": "10485760",
    "pageheap_free": "131072",
    "pageheap_unmapped": "2940928",
    "total_physical_bytes": "14165422",
    "total_thread_cache": "782208"
}
```

With the above values, it is possible to calculate the memory pressure as follows:

```math
\begin{align*}
    \text{used\_memory} &= 10485760 - 2940928 - 131072 = 7413760 \\
    \text{max\_heap\_size} &= 20971520 \\
    \text{pressure} &= \frac{7413760}{20971520} = 0.35
\end{align*}
```


## Overload Actions

When the memory pressure exceeds certain thresholds, Envoy can take [actions](https://www.envoyproxy.io/docs/envoy/latest/configuration/operations/overload_manager/overload_manager#overload-actions) to mitigate the risk of running out of memory.
An example of such actions is `stop_accepting_requests`, which stops Envoy from accepting new connections and responding with `503 Service Unavailable` to incoming requests until the memory pressure drops back below the threshold.


In the above example, the memory pressure is `0.35`.
If thresholds are configured as follows:

```yaml
overload_manager:
  resource_monitors:
    - name: "envoy.resource_monitors.fixed_heap"
      typed_config:
        "@type": type.googleapis.com/envoy.extensions.resource_monitors.fixed_heap.v3.FixedHeapConfig
        max_heap_size_bytes: 20000000
  actions:
    - name: "envoy.overload_actions.stop_accepting_requests"
      triggers:
        - name: "envoy.resource_monitors.fixed_heap"
          threshold:
            value: 0.95
```

Then no overload action is triggered since the memory pressure `0.35` is below the threshold `0.95`.


## Releasing Memory

Envoy can call TCMalloc's `MallocExtension::ReleaseMemoryToSystem()` to release free memory back to the operating system. This can happen in two ways:


There are two ways this can happen:

- Overload action (`envoy.overload_actions.shrink_heap`):

  When memory pressure exceeds a configured threshold, this action can trigger memory release.

  - Amount released per call: Fixed at 100 MB.
  - Interval: Once every 10 seconds while the threshold is triggered.
  - Metric: Number of calls is tracked by the `shrink_count` metric.

- Periodic release (`MemoryAllocatorManager`):

  The [`MemoryAllocatorManager`](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/bootstrap/v3/bootstrap.proto#config-bootstrap-v3-memoryallocatormanager) option allows periodic memory release even when there is no overload condition.

  - Amount released per call: Configurable with `bytes_to_release` (default: `0`, which disables this feature)
  - Interval: Configurable with `memory_release_interval` (default: 1 second)
  - Metric: Successful releases are counted by the `released_by_timer` metric

Releasing is performed by TCMalloc by "unmapping" free pages from the process's address space.
Unmapping is done by calling `madvise()` with the `MADV_DONTNEED` and `MADV_REMOVE` flags.
After this, the virtual address space still includes the unmapped memory, but it no longer needs to be backed by physical memory.
The memory can be reused for future allocations, but may incur a page fault when accessed again.

Detailed memory allocation statistics can be obtained by calling TCMalloc's `MallocExtension::GetStats()` function.
Envoy currently does not expose these statistics via its admin interface or metrics.

> TODO: Consider adding support for displaying these statistics, for example through a new `/memory/tcmalloc` endpoint.

## Cgroup Memory Limits

When Envoy runs in a Kubernetes pod without memory limits, the `memory.current` cgroup metric can continue to report high memory usage even after TCMalloc has released memory back to the operating system.
This is because there might not be memory pressure on the cgroup level, so the kernel does not reclaim the memory immediately.

## Source Code References

For more details, see following source files:

Envoy

- Handler for `/memory` endpoint https://github.com/envoyproxy/envoy/blob/ac9d8ba9a8f239ccee911d8a40dd35d43ed63f72/source/server/admin/server_info_handler.cc#L43-L55
- Implementation of heap allocator stats https://github.com/envoyproxy/envoy/blob/ac9d8ba9a8f239ccee911d8a40dd35d43ed63f72/source/common/memory/stats.cc
- Fixed heap monitor algorithm for calculating memory pressure https://github.com/envoyproxy/envoy/blob/ac9d8ba9a8f239ccee911d8a40dd35d43ed63f72/source/extensions/resource_monitors/fixed_heap/fixed_heap_monitor.cc#L26-L45
- Shrink Heap action
  https://github.com/envoyproxy/envoy/blob/ac9d8ba9a8f239ccee911d8a40dd35d43ed63f72/source/common/memory/heap_shrinker.cc

TCMalloc

- Properties https://github.com/google/tcmalloc/blob/12f255231938d30493186b0a037feedd70f5a1c1/tcmalloc/malloc_extension.h#L374-L416
