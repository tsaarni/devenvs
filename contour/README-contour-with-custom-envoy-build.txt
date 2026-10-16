

# Replace Envoy in Contour with a custom build for testing

kubectl -n projectcontour set image daemonset/envoy envoy=ubuntu:24.04

# Replace command, args, and run container as root to allow writing to /usr/local/bin
kubectl -n projectcontour patch daemonset/envoy --type='json' -p='[{"op":"replace","path":"/spec/template/spec/containers/1/command","value":["/bin/sh","-c","sleep 9999999999"]},{"op":"replace","path":"/spec/template/spec/containers/1/args","value":[]},{"op":"replace","path":"/spec/template/spec/containers/1/securityContext","value":{"runAsUser":0,"runAsNonRoot":false,"allowPrivilegeEscalation":true}}]'


kubectl -n projectcontour delete pod -lapp=envoy --force



bazel build -c fastbuild //source/exe:envoy-static

# Copy custom Envoy binary into the Envoy pod (fails until the pod is up)
#   log to console and file
kubectl -n projectcontour cp bazel-bin/source/exe/envoy-static $(kubectl -n projectcontour get pods -l app=envoy -o jsonpath='{.items[0].metadata.name}'):/usr/local/bin/envoy -c envoy 






kubectl -n projectcontour exec -it ds/envoy -c envoy -- envoy -c /config/envoy.json --service-cluster projectcontour --service-node envoy-debug --log-level debug


kubectl -n projectcontour exec  ds/envoy -c envoy -- pkill -9 envoy
kubectl -n projectcontour exec  ds/envoy -c envoy -- envoy -c /config/envoy.json --service-cluster projectcontour --service-node envoy-debug --log-level debug 2>&1 | tee envoy.log


