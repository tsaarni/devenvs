
# start new cluster
kind delete cluster --name contour
kind create cluster --config configs/kind-cluster-config.yaml --name contour

# deploy contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# point contour service towards the host
sed "s/REPLACE_ADDRESS_HERE/$(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')/" ~/work/devenvs/contour/manifests/contour-endpoints-dev.yaml | kubectl apply -f -

# shutdown contour inside the cluster
kubectl -n projectcontour scale deployment --replicas=0 contour
kubectl -n projectcontour rollout restart daemonset envoy

# Download contour certs
kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key


rm secrets.yaml
cert=$(base64 -w 0 certs/echoserver.pem)
key=$(base64 -w 0 certs/echoserver-key.pem)
for i in {1..2000}; do
  cat <<EOF >> secrets.yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: my-secret-$i
type: kubernetes.io/tls
data:
  tls.crt: $cert
  tls.key: $key
EOF
done

kubectl apply -f secrets.yaml


rm httpproxies.yaml
for i in {1..2000}; do
  cat <<EOF >> httpproxies.yaml
---
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
  name: echoserver-$i
spec:
  virtualhost:
    fqdn: echoserver-$i.127-0-0-101.nip.io
    tls:
      secretName: my-secret-$i
  routes:
  - services:
    - name: echoserver
      port: 80
EOF
done

kubectl apply -f httpproxies.yaml



kubectl delete -f secrets.yaml
kubectl delete -f httpproxies.yaml




make build
./contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key
./contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key --kubernetes-client-qps=500 --kubernetes-client-burst=500



go tool pprof -http=:8081 ./contour 'http://localhost:6060/debug/pprof/profile?seconds=60'
go tool pprof -http=:8081 ./contour 'http://localhost:6060/debug/pprof/heap'

http://localhost:8081/ui/flamegraph


kubectl get httpproxy echoserver-1967
http --verify=no https://echoserver-1967.127-0-0-101.nip.io



openssl s_client -connect echoserver-1967.127-0-0-101.nip.io:443 -servername echoserver-1967.127-0-0-101.nip.io | grep NotBefore
rm certs/echoserver*
certyaml -d certs configs/certs.yaml
kubectl create secret tls my-secret-1967 --cert=certs/echoserver.pem --key=certs/echoserver-key.pem --dry-run=client -o yaml | kubectl apply -f -
openssl s_client -connect echoserver-1967.127-0-0-101.nip.io:443 -servername echoserver-1967.127-0-0-101.nip.io | grep NotBefore



procpath record -p 604979 -d procpath.sqlite
procpath plot -d procpath.sqlite -q cpu -q rss








### Change sync period

cat > syncperiod.patch <<EOF
diff --git a/cmd/contour/serve.go b/cmd/contour/serve.go
index 179719f3..1574b9a5 100644
--- a/cmd/contour/serve.go
+++ b/cmd/contour/serve.go
@@ -210,11 +210,14 @@ func NewServer(log logrus.FieldLogger, ctx *serveContext) (*Server, error) {
                return nil, fmt.Errorf("unable to create scheme: %w", err)
        }

+       syncPeriod := 10 * time.Second
+
        // Instantiate a controller-runtime manager.
        options := manager.Options{
                Scheme:                 scheme,
                MetricsBindAddress:     "0",
                HealthProbeBindAddress: "0",
+               SyncPeriod:             &syncPeriod,
        }
        if ctx.LeaderElection.Disable {
                log.Info("Leader election disabled")
EOF

patch -p1 < syncperiod.patch
