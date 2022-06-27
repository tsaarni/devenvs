

docker run --volume $PWD/configs:/configs:ro --network host --rm envoyproxy/envoy-distroless:v1.22-latest --config-path /configs/envoy-overload-manager.yaml --log-level debug
python3 -m http.server 8081

http http://localhost:8080/bigfile
http http://localhost:8080/smallfile


make containr
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
