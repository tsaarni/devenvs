


### Memory footprint


# create secrets

for i in {1..100}; do
  cat <<EOF | kubectl create -f -
apiVersion: v1
kind: Secret
type: helm.sh/release.v1
metadata:
  name: my-helm-secret-$i
data:
  release: $(head -c 1MB /dev/urandom | base64 -w 0)
EOF
done



for i in {1..100}; do
  kubectl delete secret my-helm-secret-$i
done


ps_mem.py -p $(pgrep -f "contour serve")

while true; do curl --silent http://localhost:8000/metrics|grep ^go_memstats_heap_inuse_bytes; sleep .5; done

gcore $(pgrep -f "contour serve")
strings core.603923 |grep helm




- Predicates
  https://medium.com/@gallettilance/10-things-you-should-know-before-writing-a-kubernetes-controller-83de8f86d659
  "If you’re looking to filter out objects to limit memory utilization this is not the place to do that. This is where the cache comes in handy."

- SelectorByObject
  https://github.com/kubernetes-sigs/controller-runtime/blob/master/designs/use-selectors-at-cache.md
  https://github.com/kubernetes-sigs/controller-runtime/pull/1435
  Works on server-side but cannot do complicated queries like "secrets with type kubernetes.io/tls and Opaque"

- cache.TransformFunc
  https://stackoverflow.com/questions/72932546/filtering-items-are-stored-in-kubernetes-shared-informer-cache
  https://github.com/kubernetes-sigs/controller-runtime/pull/1805
  used for e.g stripping managedFields and annotations metadata
  "These caches can have tens of thousands of objects. We have found that some our controllers are being OOM-killed more frequently recently since the addition of managedFields."

- namespace filtering
  MultiNamespacedCacheBuilder(watchedNamespaces) or manager.Options{Namespace: watchedNamespace}

- metadata-only watches



Other reading
- Kubernetes Controllers at Scale: Clients, Caches, Conflicts, and Patches Explained
  https://medium.com/@timebertt/kubernetes-controllers-at-scale-clients-caches-conflicts-patches-explained-aa0f7a8b4332
- A deep dive into Kubernetes informers
  https://aly.arriqaaq.com/kubernetes-informers/
- Cert-manager
  - Investigate improving resource consumption and performance in clusters with large amount of resources #5220
    https://github.com/cert-manager/cert-manager/issues/5220
  - Design: reduce cert-manager controller's memory consumption #5639
    https://github.com/cert-manager/cert-manager/pull/5639
