


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


# dummy update to trigger cluster update in Envoy
kubectl patch httpproxy echoserver --type json -p '[{"op":"add","path":"/spec/routes/0/timeoutPolicy","value":{"response":"60s"}}]'
kubectl patch httpproxy echoserver --type json -p '[{"op":"replace","path":"/spec/routes/0/timeoutPolicy/response","value":"120s"}]'


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



cd ~/work/devenv/contour

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





# test that health check still will notice unhealthy backend
#  - scale down echoserver to 0 replicas
kubectl scale deployment echoserver --replicas=0

until http --check-status GET http://echoserver.127-0-0-101.nip.io; do sleep 1; done; date



####
### Patch contour to enable grpc-json-sniffer
###
README-xds-grpc-json-sniffer.txt

kubectl -n projectcontour port-forward deployments/contour 12345

# Open with browser
http://localhost:12345





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




bazel build -c fastbuild //source/exe:envoy-static



# Test

bazel test -c fastbuild //test/integration:sds_dynamic_integration_test --test_filter="*ClusterWarmingWhileHealthCheckBlocksOnSds*" --test_output=streamed --cache_test_results=no --test_arg="--" --test_arg="-l trace"


# this started failing after the change
bazel test -c fastbuild //test/integration:leds_integration_test --test_filter='*BasicLeds*' --test_output=streamed --cache_test_results=no --test_arg="--" --test_arg="-l trace"


bazel test -c fastbuild //test/integration:sds_dynamic_integration_test --test_filter="*ClusterWarmingWhileHealthCheckBlocksOnSdsDisabled*IPv4_EnvoyGrpc_UsesTcp*" --test_output=streamed --cache_test_results=no --test_arg="--" --test_arg="-l trace"


# health_checker  source/extensions/health_checkers/common/health_checker_base_impl.cc

HealthCheckerImplBase::addHosts()
  for each host
    # add session
    sessions[host] = makeSession()   HealthCheckerImplBase::ActiveHealthCheckSessionPtr
    ClusterImplBase::setHealthChecker(new HealthCheckHostMonitorImpl(host))
    sessions[host]->start()

HealthCheckerImplBase::onClusterMemberUpdate()
  for each new host
    addHosts()
  for each removed host
    remove session

HealthCheckerImplBase::start() {
  for cluster_->prioritySet->hostSet()
    addHosts(host_set->hosts()


# upstreamsource/common/upstream/upstream_impl.cc

ClusterImplBase::setHealthChecker()
  health_checker_->start();

ClusterImplBase::onInitDone()
  # SDS is ready




following sequence describes how the active health check is working currently

1. Cluster initialization:
   ClusterImplBase::setHealthChecker(HealthCheckerSharedPtr)
      └─> health_checker_->start()
           (or deferred to onInitDone() after SDS completes)

2. Health checker starts (after cluster warming completes):
   HealthCheckerImplBase::start()
      └─> for each host_set in cluster_.prioritySet().hostSetsPerPriority():
            addHosts(host_set->hosts())

3. Adding hosts to health checker:
   HealthCheckerImplBase::addHosts(hosts)
      └─> for each host in hosts:
            ├─> active_sessions_[host] = makeSession(host)
            │
            ├─> host->setHealthChecker(new HealthCheckHostMonitorImpl(...))
            │          │
            │          └──> Virtual dispatch through:
            │                Host interface → HostDescription interface
            │                      │
            │                      ▼
            │               HostDescriptionImplBase::setHealthChecker(HealthCheckHostMonitorPtr&&)
            │                      │
            │                      └─> health_checker_ = std::move(health_checker)
            │                          (Stores the monitor in the Host object)
            │
            └─> active_sessions_[host]->start()

4. When cluster membership changes (EDS/LEDS update):
   PrioritySet detects host changes
      └─> PrioritySet::updateHosts()
            └─> Calls registered callbacks with (hosts_added, hosts_removed)
                  └─> HealthCheckerImplBase::onClusterMemberUpdate(added, removed)
                        ├─> addHosts(added)  [goes back to step 3]
                        └─> for each removed: delete session


This does not work well, see failing test case log file in @test-fail.log

i would like to change the code to more elegantly defer the creation of health check sessions,
so that so that creating sessions and starting them is deferred until cluster warming completes.

┌─────────────────────────────────────────────────────────────────────────┐
│ During Warming                        After Warming                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   EDS arrives                         onInitDone()                      │
│       │                                   │                              │
│       ▼                                   ▼                              │
│   addHosts()                          health_checker_->start()          │
│   └─ (health checker not                  │                              │
│       started yet, so                     └─> addHosts() for all hosts  │
│       nothing happens)                        └─> makeSession() + start │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘


instead of maintaining started_ state, could we defer session creation instead






1. Cluster initialization:
   ClusterImplBase::setHealthChecker(HealthCheckerSharedPtr)
      │
      ├─> health_checker_ = health_checker
      │
      └─> health_checker_->start()

2. Health checker starts (iterates all priority host sets):
   HealthCheckerImplBase::start()
      └─> for each host_set in cluster_.prioritySet().hostSetsPerPriority():
            addHosts(host_set->hosts())

3. Adding hosts to health checker:
   HealthCheckerImplBase::addHosts(hosts)
      └─> for each host in hosts:
            │
            ├─> active_sessions_[host] = makeSession(host)
            │         │
            │         └─> Creates HttpActiveHealthCheckSession / TcpActiveHealthCheckSession / etc.
            │
            ├─> host->setHealthChecker(HealthCheckHostMonitorPtr)
            │         │
            │         └─> HostDescriptionImplBase::setHealthChecker(HealthCheckHostMonitorPtr&&)
            │                   └─> health_checker_ = std::move(health_checker)
            │
            └─> active_sessions_[host]->start()
                      │
                      └─> ActiveHealthCheckSession::onInitialInterval()
                                │
                                └─> (after jitter) interval_timer_ triggers onIntervalBase()

4. Health check interval fires - connection creation (HTTP example):
   HealthCheckerImplBase::ActiveHealthCheckSession::onIntervalBase()
      │
      ├─> onInterval()  [virtual, implemented by HttpActiveHealthCheckSession]
      │
      ├─> timeout_timer_->enableTimer(parent_.timeout_)
      │
      └─> parent_.stats_.attempt_.inc()

5. HTTP health check session creates connection:
   HttpHealthCheckerImpl::HttpActiveHealthCheckSession::onInterval()
      │
      └─> if (!client_):
            │
            ├─> Upstream::Host::CreateConnectionData conn =
            │       host_->createHealthCheckConnection(
            │           parent_.dispatcher_,
            │           parent_.transportSocketOptions(),      ← TLS options (ALPN, etc.)
            │           parent_.transportSocketMatchMetadata().get())  ← Metadata for socket matching
            │
            └─> client_.reset(parent_.createCodecClient(conn))

6. Host creates health check connection with transport socket resolution:
   HostImplBase::createHealthCheckConnection(dispatcher, transport_socket_options, metadata)
      │
      ├─> Resolve transport socket factory:
      │     │
      │     └─> if (metadata != nullptr):
      │           │
      │           └─> factory = resolveTransportSocketFactory(healthCheckAddress(), metadata)
      │                           │
      │                           └─> HostDescriptionImplBase::resolveTransportSocketFactory()
      │         else:
      │           └─> factory = transportSocketFactory()  ← Uses default/endpoint-matched factory
      │
      └─> return createConnection(dispatcher, cluster(), healthCheckAddress(),
                                   {}, factory, nullptr, transport_socket_options, host)

7. Transport socket factory resolution (with metadata):
   HostDescriptionImplBase::resolveTransportSocketFactory(dest_address, endpoint_metadata)
      │
      └─> cluster_->transportSocketMatcher().resolve(endpoint_metadata, locality_metadata_.get())
                │
                └─> TransportSocketMatcherImpl::resolve(endpoint_metadata, locality_metadata, ...)
                          │
                          ├─> If matcher_ available:
                          │     └─> resolveUsingMatcher()  ← xDS matcher-based resolution
                          │
                          ├─> Else (legacy metadata matching):
                          │     │
                          │     ├─> for each match in matches_:
                          │     │     if Config::Metadata::metadataLabelMatch(match.label_set,
                          │     │            endpoint_metadata, ENVOY_TRANSPORT_SOCKET_MATCH):
                          │     │       return {match.factory, match.stats, match.name}
                          │     │
                          │     ├─> for each match in matches_:  ← fallback to locality metadata
                          │     │     if metadataLabelMatch(..., locality_metadata, ...):
                          │     │       return {match.factory, ...}
                          │     │
                          │     └─> return default_match_.factory  ← If no match found
                          │
                          └─> Returns MatchData: {factory&, stats, name}

8. Static helper creates the actual connection:
   HostImplBase::createConnection(dispatcher, cluster, address, address_list,
                                   socket_factory, options, transport_socket_options, host)
      │
      ├─> source_address_selector = cluster.getUpstreamLocalAddressSelector()
      │
      ├─> Check for proxy redirect (HTTP/1.1 CONNECT proxy)
      │
      └─> Normal path (no proxy, single address):
            │
            ├─> upstream_local_address = source_address_selector->getUpstreamLocalAddress(...)
            │
            └─> connection = dispatcher.createClientConnection(
                    address,                                          ← healthCheckAddress()
                    upstream_local_address.address_,                  ← Local/source address
                    socket_factory.createTransportSocket(             ← Creates TLS socket
                        transport_socket_options, host),
                    upstream_local_address.socket_options_,
                    transport_socket_options)

9. Transport socket creation (TLS example):
   ClientSslSocketFactory::createTransportSocket(transport_socket_options, host)
      │
      ├─> Acquire read lock on ssl_ctx_mu_
      │
      ├─> ssl_ctx = ssl_ctx_  ← Get current SSL context (may be updated by SDS)
      │
      └─> if (ssl_ctx):
            │
            └─> SslSocket::create(ssl_ctx, InitialState::Client,
                                   transport_socket_options, handshaker, host)
                      │
                      └─> Creates SslSocket with:
                            - ClientContextSharedPtr (contains certs, keys, CA)
                            - Transport socket options (ALPN protocols)
                            - Host description (for SNI, etc.)
          else:
            └─> return NotReadySslSocket()  ← Secrets not ready yet

10. SSL Context contains (populated from SDS or static config):
    ClientContextImpl:
      ├─> tls_certificates_  ← Client certificate + private key (for mTLS)
      ├─> cert_validator_    ← CA certificates for server verification
      ├─> alpn_protocols_    ← Application-Layer Protocol Negotiation
      └─> session_keys_      ← TLS session resumption data

11. When cluster membership changes (EDS/LEDS update):
    PrioritySet::updateHosts()
       └─> Calls registered callbacks with (hosts_added, hosts_removed)
             └─> HealthCheckerImplBase::onClusterMemberUpdate(added, removed)
                   ├─> addHosts(added)  [goes back to step 3]
                   └─> for each removed:
                         └─> session->onDeferredDeleteBase()
                             dispatcher_.deferredDelete(session)



┌─────────────────┐      ┌──────────────────┐      ┌─────────────────┐      ┌──────────────────┐
│ ClusterImplBase │      │HealthCheckerImpl │      │    HostImpl     │      │ TransportSocket  │
│                 │      │      Base        │      │                 │      │    Matcher       │
└────────┬────────┘      └────────┬─────────┘      └────────┬────────┘      └────────┬─────────┘
         │                        │                         │                        │
         │ setHealthChecker()     │                         │                        │
         ├───────────────────────>│                         │                        │
         │                        │                         │                        │
         │                        │ start() → addHosts()    │                        │
         │                        ├────────────────────────>│                        │
         │                        │  setHealthChecker()     │                        │
         │                        │                         │                        │
         │                        │ [timer fires]           │                        │
         │                        │ onInterval()            │                        │
         │                        │                         │                        │
         │                        │ createHealthCheckConnection(metadata)            │
         │                        ├────────────────────────>│                        │
         │                        │                         │                        │
         │                        │                         │ resolve(metadata)      │
         │                        │                         ├───────────────────────>│
         │                        │                         │                        │
         │                        │                         │<── SocketFactory& ─────┤
         │                        │                         │                        │
         │                        │                         │ createTransportSocket()│
         │                        │                         ├───────────────────────>│
         │                        │                         │                        │
         │                        │<── Connection ──────────┤                        │



┌─────────────────────────────────────────────────────────────────────────────┐
│                         CLUSTER WARMING TIMELINE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  warming_state_ = 1                                warming_state_ = 0       │
│        │                                                  │                 │
│        ▼                                                  ▼                 │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────────┐          │
│  │ Cluster  │───>│ PreInit  │───>│  Init    │───>│  onInitDone  │          │
│  │ created  │    │ complete │    │ Manager  │    │   called     │          │
│  └──────────┘    └──────────┘    │ running  │    └──────────────┘          │
│                                  └──────────┘                               │
│                                       │                                     │
│                                       │ All targets call ready()            │
│                                       ▼                                     │
│                              ┌────────────────┐                             │
│                              │ init_watcher_  │                             │
│                              │ callback fires │                             │
│                              └────────────────┘                             │
└─────────────────────────────────────────────────────────────────────────────┘
