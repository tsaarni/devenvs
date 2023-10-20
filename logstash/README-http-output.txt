
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml



docker-compose rm -f
docker-compose up http

# test server
http --verify certs/server-ca.pem --cert certs/client.pem --cert-key certs/client-key.pem https://localhost:8443/

# make request with curl
curl --cacert certs/server-ca.pem --cert certs/client.pem --key certs/client-key.pem https://localhost:8443/



docker run --rm -it --user=$(id -u):$(id -g) --volume=$HOME/work/devenvs/logstash:/input:ro --network=host docker.elastic.co/logstash/logstash-oss:8.10.2 bin/logstash -f /input/configs/logstash-source-http.conf
