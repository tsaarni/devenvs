
docker compose -f docker-compose-overload.yaml rm -f
docker compose -f docker-compose-overload.yaml up

# Grafana dashboard
http://localhost:3000/




http http://localhost:9901/stats/prometheus

http http://localhost:9901/memory




while true; do http http://localhost:8080; sleep 1; done



dd if=/dev/zero bs=100M count=1 | http -v POST http://localhost:8080/upload?throttle=1K
