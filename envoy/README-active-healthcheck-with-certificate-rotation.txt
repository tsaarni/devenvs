


### Rotating certificates for active healthcheck does not reproduce the issue

rm -rf certs
mkdir -p certs
certyaml --destination certs configs/certs.yaml


# rotate certs with atomic rename
function rotate_certs() {
  local timestamp=$(date +%s%3N)
  mkdir -p certs/$timestamp

  rm certs/envoy-client*
  certyaml --destination certs configs/certs.yaml > /dev/null 2>&1
  chmod +r certs/*
  cp certs/envoy-client.pem certs/envoy-client-key.pem certs/$timestamp/

  ln -sf $timestamp new-symlink-$timestamp
  mv -T -f new-symlink-$timestamp certs/effective-certs

  # Print the new serial number (in decimal) for envoy-client.pem
  openssl x509 -in certs/envoy-client.pem -noout -text | grep 'Serial Number' | awk '{print $3}'
}


rotate_certs


docker compose -f docker-compose-upstream-tls-rotation.yaml up


http http://localhost:8080

rotate_certs




sudo nsenter --target $(pidof echoserver) --net wireshark -i any -f "port 8443" -Y tls -o tls.keylog_file:/proc/$(pidof echoserver)/root/wireshark-keys.log -k


source/common/upstream/health_checker_impl.h
source/common/upstream/health_checker_impl.cc

source/extensions/health_checkers/http/health_checker_impl.h
source/extensions/health_checkers/http/health_checker_impl.cc
