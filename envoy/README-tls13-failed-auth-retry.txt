
# https://github.com/envoyproxy/envoy/issues/9300


mkdir -p certs
certyaml -d certs configs/certs.yaml

docker-compose -f docker-compose-tls13-failed-auth-retry.yaml rm -f
docker-compose -f docker-compose-tls13-failed-auth-retry.yaml up

http http://localhost:8080


# To check HTTP connection manager and TCP proxy behavior, switch between HTTP and TCP proxy envoy configs in the compose file



echoserver_pid=$(pgrep -f echoserver | head -1)
sudo nsenter -t $echoserver_pid --net wireshark -f "port 8443" -k -o tls.keylog_file:/proc/$echoserver_pid/root/tmp/wireshark-keys.log
