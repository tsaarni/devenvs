# issue and analysis about EDS performance problem
# https://github.com/projectcontour/contour/issues/6743


# change to separate snapshot cache for EDS and comment about LinearCache
https://github.com/projectcontour/contour/pull/6250/files#diff-1900033d7ec953a2e318fa5369de7c7d199c4e478873a17e61b7f3776771bb1e


kubectl apply -f manifests/echoserver.yaml



function echoserver() {
    local action=$1
    local num=$2
    cat <<EOF | kubectl $action -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
    name: echoserver-$num
    annotations:
        date: "$(date)"
spec:
    virtualhost:
        fqdn: echoserver-$num.127-0-0-101.nip.io
    routes:
        - services:
            - name: echoserver-$num
              port: 80
---
apiVersion: v1
kind: Service
metadata:
    name: echoserver-$num
    annotations:
        date: "$(date)"
spec:
    selector:
        app: echoserver
    ports:
        - protocol: TCP
          port: 80
          targetPort: http-api
---
EOF
kubectl get httpproxies,services -o wide
}


function httpproxy() {
    local action=$1
    local httpproxy_name=$2
    local service_name=$3
    cat <<EOF | kubectl $action -f -
apiVersion: projectcontour.io/v1
kind: HTTPProxy
metadata:
    name: ${httpproxy_name}
    annotations:
        date: "$(date)"
spec:
    virtualhost:
        fqdn: ${httpproxy_name}.127-0-0-101.nip.io
    routes:
        - services:
            - name: ${service_name}
              port: 80
EOF
kubectl get httpproxies,services -o wide
}

function service() {
    local action=$1
    local service_name=$2
    cat <<EOF | kubectl $action -f -
apiVersion: v1
kind: Service
metadata:
    name: ${service_name}
    annotations:
        date: "$(date)"
spec:
    selector:
        app: echoserver
    ports:
        - protocol: TCP
          port: 80
          targetPort: http-api
EOF
kubectl get httpproxies,services -o wide
}

echoserver create 0001
echoserver delete 0001
echoserver apply 0001


kubectl delete pod -l app.kubernetes.io/name=echoserver --wait=false


http http://echoserver-0000.127-0-0-101.nip.io
http http://echoserver-7999.127-0-0-101.nip.io



# contour metrics
http http://localhost:8000/metrics
http http://localhost:8000/metrics | grep EndpointDiscoveryService

# envoy metrics
kubectl -n projectcontour port-forward daemonset/envoy 9001:9001
http http://localhost:9001/stats



kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key

kubectl -n projectcontour get secret envoycert -o jsonpath='{..ca\.crt}' | base64 -d > envoy-ca.crt
kubectl -n projectcontour get secret envoycert -o jsonpath='{..tls\.crt}' | base64 -d > envoy-tls.crt
kubectl -n projectcontour get secret envoycert -o jsonpath='{..tls\.key}' | base64 -d > envoy-tls.key


###kubectl -n projectcontour port-forward deployment/contour 8001:8001



go run github.com/projectcontour/contour/cmd/contour cli eds --cafile=envoy-ca.crt --cert-file=envoy-tls.crt --key-file=envoy-tls.key --contour="localhost:8001"

kubectl -n projectcontour rollout restart daemonset envoy

kubectl rollout restart deployment echoserver
kubectl get pod -o wide



# Patch the code for grpc-json-sniffer

go get github.com/tsaarni/grpc-json-sniffer


patch -p1 <<EOF
diff --git a/internal/xds/server.go b/internal/xds/server.go
index 90d0e330..56c66926 100644
--- a/internal/xds/server.go
+++ b/internal/xds/server.go
@@ -16,6 +16,7 @@ package xds
 import (
 	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
 	"github.com/prometheus/client_golang/prometheus"
+	sniffer "github.com/tsaarni/grpc-json-sniffer"
 	"google.golang.org/grpc"
 )

@@ -29,9 +30,11 @@ func NewServer(registry *prometheus.Registry, opts ...grpc.ServerOption) *grpc.S
 		metrics = grpc_prometheus.NewServerMetrics()
 		registry.MustRegister(metrics)

+		interceptor, _ := sniffer.NewGrpcJsonInterceptor()
+
 		opts = append(opts,
-			grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
-			grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
+			grpc.ChainStreamInterceptor(metrics.StreamServerInterceptor(), interceptor.StreamServerInterceptor()),
+			grpc.ChainUnaryInterceptor(metrics.UnaryServerInterceptor(), interceptor.UnaryServerInterceptor()),
 		)
 	}


EOF




# Add to .vscode/launch.json

{
        "version": "0.2.0",
        "configurations": [
                        "env": {
                                "GRPC_JSON_SNIFFER_FILE": "grpc_capture.json",
                                "GRPC_JSON_SNIFFER_ADDR": "localhost:8080",
                        },
     ]
}



apps/eds-message-parser.py /home/tsaarni/work/contour-worktree/eds-performance-fix/grpc_capture.json




# Run with unmodified contour

cd ~/work/contour
export GRPC_JSON_SNIFFER_FILE=grpc_capture.json
export GRPC_JSON_SNIFFER_ADDR=localhost:8080


kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key

http http://echoserver.127-0-0-101.nip.io



apps/eds-message-parser.py /home/tsaarni/work/contour/grpc_capture.json
