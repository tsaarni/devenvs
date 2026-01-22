



# Add json sniffer as dependency by running following in source directory

go get github.com/tsaarni/grpc-json-sniffer


# Apply following changes to the code

diff --git a/internal/xds/server.go b/internal/xds/server.go
index 90d0e330f..32e51ebe6 100644
--- a/internal/xds/server.go
+++ b/internal/xds/server.go
@@ -16,6 +16,7 @@ package xds
 import (
        grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
        "github.com/prometheus/client_golang/prometheus"
+       sniffer "github.com/tsaarni/grpc-json-sniffer"
        "google.golang.org/grpc"
 )

@@ -29,9 +30,14 @@ func NewServer(registry *prometheus.Registry, opts ...grpc.ServerOption) *grpc.S
                metrics = grpc_prometheus.NewServerMetrics()
                registry.MustRegister(metrics)

+               interceptor, err := sniffer.NewGrpcJsonInterceptor()
+               if err != nil {
+                       panic(err)
+               }
+
                opts = append(opts,
-                       grpc.StreamInterceptor(metrics.StreamServerInterceptor()),
-                       grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
+                       grpc.ChainStreamInterceptor(metrics.StreamServerInterceptor(), interceptor.StreamServerInterceptor()),
+                       grpc.ChainUnaryInterceptor(metrics.UnaryServerInterceptor(), interceptor.UnaryServerInterceptor()),
                )
        }


go mod tidy       # to add new dep to go.mod


# Start a new cluster and install contour

kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml


# Build new contour image with json sniffer and deploy it to kind cluster

cd ~/work/contour
make container

docker tag ghcr.io/projectcontour/contour:$(git describe --tags --exact-match 2>/dev/null || git rev-parse --short=8 --verify HEAD) localhost/contour:latest

kind load docker-image localhost/contour:latest --name contour



# Patch the contour deployment and envoy daemonset to use the new image and enable json sniffer

cat <<EOF | kubectl -n projectcontour patch deployment contour --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: contour
        image: localhost/contour:latest
        imagePullPolicy: Never
        env:
        - name: GRPC_JSON_SNIFFER_ADDR
          value: "0.0.0.0:12345"
        - name: GRPC_JSON_SNIFFER_FILE
          value: "/capture/grpc-capture.json"
        volumeMounts:
        - name: capture-file
          mountPath: /capture
      volumes:
      - name: capture-file
        emptyDir:
          medium: Memory
EOF

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



kubectl -n projectcontour scale deployment/contour --replicas 1
kubectl -n projectcontour port-forward deployments/contour 12345


# Open with browser

http://localhost:12345



# To enable debug logging in envoy, run following

cat <<EOF | kubectl -n projectcontour patch daemonset envoy --patch-file=/dev/stdin
spec:
  template:
    spec:
      containers:
      - name: envoy
        args:
        - -c
        - /config/envoy.json
        - --service-cluster \$(CONTOUR_NAMESPACE)
        - --service-node \$(ENVOY_POD_NAME)
        - --log-level debug
EOF

kubectl -n projectcontour logs -f ds/envoy -c envoy
