


# Start new cluster.
kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour



kubectl apply -f https://projectcontour.io/quickstart/contour.yaml



# Create secrets.
kubectl create secret generic internal-root-ca --from-file=ca.crt=certs/internal-root-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls echoserver-cert --cert=certs/echoserver.pem --key=certs/echoserver-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl apply -f -


# Configure TLS between Envoy and upstream.
kubectl -n projectcontour create configmap contour --from-file=contour.yaml=configs/contour-config-with-upstream-tls.yaml --dry-run=client -o yaml | kubectl -n projectcontour apply -f -

kubectl -n projectcontour delete pod -lapp=envoy --force
kubectl -n projectcontour scale deployment contour --replicas=0
kubectl -n projectcontour scale deployment contour --replicas=1


# Deploy echoserver with active health check.
kubectl apply -f manifests/echoserver-active-healthcheck-yaml



http http://echoserver.127-0-0-101.nip.io


http http://echoserver.127-0-0-101.nip.io | jq -r .tls.peer_certificates_decoded



# Start wireshark to capture and decrypt TLS traffic, using the keylog file from echoserver pod.
sudo nsenter -t $(pgrep -f "\./echoserver") --net wireshark -f "port 8443" -k -o tls.keylog_file:/proc/$(pgrep -f "\./echoserver")/root/tmp/wireshark-keys.log

# 1. wireshark / file / save as:  /tmp/toubleshooting.pcapng
# 2. wireshark / file / export TLS session keys:  /tmp/toubleshooting.keys
# 3. append keys to pcapng:
sudo chown $USER:$USER /tmp/troubleshooting.{keys,pcapng}
editcap --inject-secrets tls,/tmp/troubleshooting.keys /tmp/troubleshooting.pcapng troubleshooting.pcapng



# Envoy statistics
kubectl -n projectcontour port-forward daemonset/envoy 8002:8002
http http://localhost:8002/stats|grep echoserver|grep -E "(health_check|upstream_cx)"


# Turn echoserver unhealthy
http http://echoserver.127-0-0-101.nip.io/status?set=401





printf "\n\n--- envoy /stats before rotation: $(date) ---\n\n\n" > troubleshooting-logs.txt
http http://localhost:8002/stats >> troubleshooting-logs.txt


# Rotate envoy client certificate
date
rm certs/envoy.pem certs/envoy-key.pem
certyaml --destination certs configs/certs.yaml
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl -n projectcontour apply -f -
until http --check-status GET http://echoserver.127-0-0-101.nip.io; do sleep 1; done; date
date


while true; do http --print=h  GET http://echoserver.127-0-0-101.nip.io; sleep 1; done    # ctlrl+C to stop after observing FAILUREs

printf "\n\n--- envoy /stats after rotation: $(date) ---\n\n\n" >> troubleshooting-logs.txt
http http://localhost:8002/stats >> troubleshooting-logs.txt


printf "\n\n--- envoy logs  $(date) ---\n\n\n" >> troubleshooting-logs.txt
kubectl -n projectcontour logs daemonset/envoy  -c envoy >> troubleshooting-logs.txt


kindps contour-worker contour
sudo cp /proc/$(kindps contour-worker contour -o json | jq -r .[].pids[0].pid)/root/capture/grpc-capture.json troubleshooting-grpc-capture-after-rotation.json
sudo chown $USER:$USER troubleshooting-grpc-capture-after-rotation.json

go run github.com/tsaarni/grpc-json-sniffer/cmd/grpc-json-sniffer-viewer@latest troubleshooting-client-secret-rotation/troubleshooting-grpc-capture-after-rotation.json

http http://localhost:8002/stats|grep echoserver|grep -E "(health_check|upstream_cx)"


http http://echoserver.127-0-0-101.nip.io | jq -r .tls.peer_certificates_decoded



### Patch contour to enable grpc-json-sniffer
README-xds-grpc-json-sniffer.txt

kubectl -n projectcontour port-forward deployments/contour 12345

# Open with browser
http://localhost:12345





kubectl -n projectcontour port-forward daemonset/envoy 9001

# list all keys in config dump
http http://localhost:9001/config_dump?include_eds | jq -r 'paths | map(if (type == "number") then "[" + tostring + "]" else "[\"" + tostring + "\"]" end) | join("")'

http http://localhost:9001/config_dump?include_eds | jq '.["configs"][1]["dynamic_active_clusters"][0]'

{
  "version_info": "954d4a8d-90b3-48e3-a82c-dd56d417e8f2",
  "cluster": {
    "@type": "type.googleapis.com/envoy.config.cluster.v3.Cluster",
    "name": "default/echoserver/443/2cd1eacd71",
    "type": "EDS",
    "eds_cluster_config": {
      "eds_config": {
        "api_config_source": {
          "api_type": "GRPC",
          "grpc_services": [
            {
              "envoy_grpc": {
                "cluster_name": "contour",
                "authority": "contour"
              }
            }
          ],
          "transport_api_version": "V3"
        },
        "resource_api_version": "V3"
      },
      "service_name": "default/echoserver/https"
    },
    "connect_timeout": "2s",
    "health_checks": [
      {
        "timeout": "2s",
        "interval": "5s",
        "unhealthy_threshold": 3,
        "healthy_threshold": 5,
        "http_health_check": {
          "host": "contour-envoy-healthcheck",
          "path": "/status"
        }
      }
    ],
    "transport_socket": {
      "name": "envoy.transport_sockets.tls",
      "typed_config": {
        "@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext",
        "common_tls_context": {
          "tls_params": {
            "tls_minimum_protocol_version": "TLSv1_2",
            "tls_maximum_protocol_version": "TLSv1_3",
            "cipher_suites": [
              "[ECDHE-ECDSA-AES128-GCM-SHA256|ECDHE-ECDSA-CHACHA20-POLY1305]",
              "[ECDHE-RSA-AES128-GCM-SHA256|ECDHE-RSA-CHACHA20-POLY1305]",
              "ECDHE-ECDSA-AES256-GCM-SHA384",
              "ECDHE-RSA-AES256-GCM-SHA384"
            ]
          },
          "validation_context": {
            "trusted_ca": {
              "inline_bytes": "< CA CERT HERE >"
            },
            "match_typed_subject_alt_names": [
              {
                "san_type": "DNS",
                "matcher": {
                  "exact": "echoserver"
                }
              }
            ]
          },
          "tls_certificate_sds_secret_configs": [
            {
              "name": "projectcontour/envoy-client-cert/4c5d24829f",
              "sds_config": {
                "api_config_source": {
                  "api_type": "GRPC",
                  "grpc_services": [
                    {
                      "envoy_grpc": {
                        "cluster_name": "contour",
                        "authority": "contour"
                      }
                    }
                  ],
                  "transport_api_version": "V3"
                },
                "resource_api_version": "V3"
              }
            }
          ]
        }
      }
    },
    "common_lb_config": {
      "healthy_panic_threshold": {}
    },
    "alt_stat_name": "default_echoserver_443",
    "ignore_health_on_host_removal": true
  },
  "last_updated": "2025-12-10T10:07:59.571Z"
}


https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/health_check.proto#envoy-v3-api-msg-config-core-v3-healthcheck

"timeout": "2s",          # The time to wait for a health check response. If the timeout is reached the health
                          # check attempt will be considered a failure.

"interval": "5s",         # The interval between health checks.

"unhealthy_threshold": 3, # The number of unhealthy health checks required before a host is marked unhealthy.
                          # Note that for http health checking if a host responds with a code not in
                          # expected_statuses or retriable_statuses, this threshold is ignored and the host
                          # is considered immediately unhealthy

"healthy_threshold": 5,   # The number of healthy health checks required before a host is marked healthy.
                          # Note that during startup, only a single successful health check is required to
                          # mark a host healthy.



10:34:38   - rotate cert
           - requests towards backend still work     for about 5 seconds
10:34:43   - requests stop working                   for about 20 seconds
                503 Service Unavailable
                no healthy upstream
10:35:04   - requests start working again            26 seconds after cert rotation



~/work/devenvs/contour/healthcheck-upstream-cert-rotate-grpc-capture.json

first rotation triggered at 12:23:23   (no problems)

12:23:23.499  StreamSecrets  Response  <empty>
12:23:23.500  StreamClusters Response  upstream sds secret config name: projectcontour/envoy-client-cert/1fce084f86
12:23:23.514  StreamSecrets  Response  projectcontour/envoy-client-cert/1fce084f86


second rotation triggered at 12:34:38  (503 for about 20 seconds)

12:34:38.402  StreamSecrets  Response  <empty>
12:34:38.402  StreamClusters Response  upstream sds secret config name: projectcontour/envoy-client-cert/97af60fef6
12:34:38.412  StreamSecrets  Response  projectcontour/envoy-client-cert/97af60fef6



echoserver logs (from another test)

BEFORE      ROTATION

2025/12/10 12:17:04 DEBUG Handling status request url=/status
2025/12/10 12:17:04 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T10:34:38.000Z peer_certificate_not_after=2026-12-10T10:34:38.000Z peer_certificate_serial_number=1765362878274987535

14:17:08    ROTATION (did not affect the used client cert: still same client cert)

2025/12/10 12:17:09 http: TLS handshake error from 10.244.1.4:39898: EOF
2025/12/10 12:17:09 DEBUG Handling echo request method=GET url=/ remote=10.244.1.4:39904
2025/12/10 12:17:09 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T10:34:38.000Z peer_certificate_not_after=2026-12-10T10:34:38.000Z peer_certificate_serial_number=1765362878274987535
2025/12/10 12:17:09 DEBUG Handling status request url=/status
2025/12/10 12:17:09 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T10:34:38.000Z peer_certificate_not_after=2026-12-10T10:34:38.000Z peer_certificate_serial_number=1765362878274987535
2025/12/10 12:17:14 DEBUG Handling echo request method=GET url=/ remote=10.244.1.4:34414
2025/12/10 12:17:14 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T10:34:38.000Z peer_certificate_not_after=2026-12-10T10:34:38.000Z peer_certificate_serial_number=1765362878274987535

12:17:14    First health check request with new client cert  (6 seconds after rotation)
                errors starting:
                    503 Service Unavailable
                    no healthy upstream

2025/12/10 12:17:14 DEBUG Handling status request url=/status
2025/12/10 12:17:14 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603
2025/12/10 12:17:19 DEBUG Handling status request url=/status
2025/12/10 12:17:19 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603
2025/12/10 12:17:24 DEBUG Handling status request url=/status
2025/12/10 12:17:24 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603
2025/12/10 12:17:29 DEBUG Handling status request url=/status
2025/12/10 12:17:29 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603

12:17:34    First request goes through after rotation  (20 seconds after first error)
                requests start working again

2025/12/10 12:17:34 DEBUG Handling status request url=/status
2025/12/10 12:17:34 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603
2025/12/10 12:17:34 DEBUG Handling echo request method=GET url=/ remote=10.244.1.4:46574
2025/12/10 12:17:34 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603
2025/12/10 12:17:39 DEBUG Handling status request url=/status
2025/12/10 12:17:39 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603
2025/12/10 12:17:44 DEBUG Handling status request url=/status
2025/12/10 12:17:44 DEBUG TLS info version="TLS 1.3" cipher_suite=TLS_AES_128_GCM_SHA256 peer_certificate_subject="CN=envoy" peer_certificate_not_before=2025-12-10T12:17:08.000Z peer_certificate_not_after=2026-12-10T12:17:08.000Z peer_certificate_serial_number=1765369028974114603






date
kubectl patch httpproxy echoserver --type='json' -p='[{"op": "replace", "path": "/spec/routes/0/healthCheckPolicy/path", "value": "/status/200"}]'
until http --check-status GET http://echoserver.127-0-0-101.nip.io; do sleep 1; done; date


date
kubectl patch httpproxy echoserver --type='json' -p='[{"op": "replace", "path": "/spec/routes/0/healthCheckPolicy/path", "value": "/status"}]'
until http --check-status GET http://echoserver.127-0-0-101.nip.io; do sleep 1; done; date



OTHER PROBLEMS

- Sometimes when adding HA to TLS upstream, healthcheck is executed once per minute only, even if it returns success






### Non-TLS test (cluster change without TLS does not cause downtime)

kubectl edit httpproxy echoserver  # remove tls and change port to 80

### Check that echoserver is reachable.
http http://echoserver.127-0-0-101.nip.io

# Start wireshark to capture HTTP traffic between Envoy and echoserver.
sudo nsenter -t $(pgrep -f "\./echoserver") --net wireshark -f "port 8080" -k


rm troubleshooting-logs.txt

printf "\n\n--- envoy /stats before change: $(date) ---\n\n\n" > troubleshooting-logs.txt
http http://localhost:8002/stats >> troubleshooting-logs.txt

# Trigger change in cluster
date
kubectl patch httpproxy echoserver --type merge -p '{"spec":{"routes":[{"healthCheckPolicy":{"healthyThresholdCount":5,"intervalSeconds":5,"path":"/status","timeoutSeconds":2,"unhealthyThresholdCount":3},"loadBalancerPolicy":{"strategy":"WeightedLeastRequest"},"services":[{"name":"echoserver","port":80}]}]}}'

while true; do http --print=h  GET http://echoserver.127-0-0-101.nip.io; sleep 1; done


printf "\n\n--- envoy /stats after change: $(date) ---\n\n\n" > troubleshooting-logs.txt
http http://localhost:8002/stats >> troubleshooting-logs.txt


printf "\n\n--- envoy logs  $(date) ---\n\n\n" >> troubleshooting-logs.txt
kubectl -n projectcontour logs daemonset/envoy  -c envoy >> troubleshooting-logs.txt

















**If you are reporting *any* crash or *any* potential security issue, *do not*
open an issue in this repo. Please report the issue via emailing
envoy-security@googlegroups.com where the issue will be triaged appropriately.**

*Title*: *One line description*

*Description*:
>What issue is being seen? Describe what should be happening instead of
the bug, for example: Envoy should not crash, the expected value isn't
returned, etc.

*Repro steps*:
> Include sample requests, environment, etc. All data and inputs
required to reproduce the bug.

>**Note**: The [Envoy_collect tool](https://github.com/envoyproxy/envoy/blob/main/tools/envoy_collect/README.md)
gathers a tarball with debug logs, config and the following admin
endpoints: /stats, /clusters and /server_info. Please note if there are
privacy concerns, sanitize the data prior to sharing the tarball/pasting.

*Admin and Stats Output*:
>Include the admin output for the following endpoints: /stats,
/clusters, /routes, /server_info. For more information, refer to the
[admin endpoint documentation.](https://www.envoyproxy.io/docs/envoy/latest/operations/admin)

>**Note**: If there are privacy concerns, sanitize the data prior to
sharing.

*Config*:
>Include the config used to configure Envoy.

*Logs*:
>Include the access logs and the Envoy logs.

>**Note**: If there are privacy concerns, sanitize the data prior to
sharing.

*Call Stack*:
> If the Envoy binary is crashing, a call stack is **required**.
Please refer to the [Bazel Stack trace documentation](https://github.com/envoyproxy/envoy/tree/main/bazel#stack-trace-symbol-resolution).
