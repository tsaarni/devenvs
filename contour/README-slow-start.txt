

kubectl -n projectcontour get secret contourcert -o jsonpath='{..ca\.crt}' | base64 -d > ca.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.crt}' | base64 -d > tls.crt
kubectl -n projectcontour get secret contourcert -o jsonpath='{..tls\.key}' | base64 -d > tls.key

go run github.com/projectcontour/contour/cmd/contour serve --xds-address=0.0.0.0 --xds-port=8001 --envoy-service-http-port=8080 --envoy-service-https-port=8443 --contour-cafile=ca.crt --contour-cert-file=tls.crt --contour-key-file=tls.key



kubectl apply -f examples/contour/01-crds.yaml

kubectl apply -f manifests/slowstart.yaml

http http://echoserver.127-0-0-101.nip.io



# check that the slow start settings are visible in envoy config_dump
kubectl -n projectcontour port-forward daemonset/envoy 9001
http http://localhost:9001/config_dump| jq '.configs[].dynamic_active_clusters'



# run traffic tests


git clone https://github.com/grafana/k6
cd k6
docker-compose up influxdb grafana


1. Navigate to http://localhost:3000/.
2. Click "Dashboards" icon in left panel, select + import.
3. Copy content of configs/slowstart-grafana-panel.json into Import via panel json text box and click Load.
4. Click Import.

Run the test

docker run --rm -it --network host -e K6_OUT=influxdb=http://localhost:8086/k6 -v $PWD:/input:ro grafana/k6:latest run --vus 5 --duration 1h /input/k6-slowstart.js

pick the PID of one of the echoserver replicas and kill it

$ ps -ef|grep echoserver
65532    3741902 3741816 10 20:56 ?        00:01:47 /echoserver
65532    3741910 3741792 10 20:56 ?        00:01:47 /echoserver
65532    3742037 3741954 10 20:56 ?        00:01:47 /echoserver
65532    3742229 3742092 10 20:56 ?        00:01:47 /echoserver
65532    3742236 3742137 10 20:56 ?        00:01:47 /echoserver
tsaarni  3746917    9433  0 21:13 pts/1    00:00:00 grep --color echoserver
$ sudo kill 3741902

Stop the test
