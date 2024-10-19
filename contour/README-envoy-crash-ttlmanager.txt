
# Envoy was reported to be crashing on TTLManager, see file ~/work/devenv/envoy/README-envoy-crash-in-ttlmanager.txt
# This is an attempt to reproduce a crash in context of Contour




### Deploy echoserver
kubectl apply -f manifests/echoserver-with-upstream-tls.yaml



### Add socket option, IPv6 traffic class will fail in kind
kubectl apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
    name: contour
    namespace: projectcontour
data:
    contour.yaml: |
        listener:
            socket-options:
                tos: 64
                traffic-class: 64
        tls:
            envoy-client-certificate:
                name: envoy-client-cert
                namespace: projectcontour
EOF

kubectl -n projectcontour edit configmap contour


# Restart


## Contour after changing the configmap
kubectl -n projectcontour scale deployment contour --replicas=0
kubectl -n projectcontour scale deployment contour --replicas=1

## Restart Envoy
kubectl -n projectcontour rollout restart daemonset envoy

kubectl delete pod -n projectcontour -l app=contour


### Logs
kubectl -n projectcontour logs deployment/contour -f
kubectl -n projectcontour logs daemonsets/envoy -c envoy -f


### Test sending requests
http --verify certs/external-root-ca.pem https://protected.127-0-0-101.nip.io
http --verify certs/external-root-ca.pem https://protected.127-0-0-101.nip.io | jq -r '.tls.peerCertificates[0]' | openssl x509 -text -noout


### rotate internal certificate
rm certs/envoy.pem certs/envoy-key.pem
certyaml --destination certs configs/certs.yaml
kubectl -n projectcontour create secret tls envoy-client-cert --cert=certs/envoy.pem --key=certs/envoy-key.pem --dry-run=client -o yaml | kubectl apply -f -
openssl x509 -in certs/envoy.pem -noout -dates








#######################################





[2024-10-18 10:04:27.640][1][info][lua] [source/extensions/filters/http/lua/lua_filter.cc:223] envoy_on_request() function not found. Lua filter will not hook requests.
[2024-10-18 10:04:27.640][1][info][lua] [source/extensions/filters/http/lua/lua_filter.cc:228] envoy_on_response() function not found. Lua filter will not hook responses.
[2024-10-18 10:04:27.641][1][info][upstream] [source/common/listener_manager/lds_api.cc:106] lds: add/update listener 'ingress_http'
[2024-10-18 10:04:27.645][1][info][lua] [source/extensions/filters/http/lua/lua_filter.cc:228] envoy_on_response() function not found. Lua filter will not hook responses.
[2024-10-18 10:04:27.647][1][info][lua] [source/extensions/filters/http/lua/lua_filter.cc:223] envoy_on_request() function not found. Lua filter will not hook requests.
[2024-10-18 10:04:27.647][1][info][lua] [source/extensions/filters/http/lua/lua_filter.cc:228] envoy_on_response() function not found. Lua filter will not hook responses.
[2024-10-18 10:04:27.649][1][info][upstream] [source/common/listener_manager/lds_api.cc:106] lds: add/update listener 'ingress_https'
[2024-10-18 10:04:27.651][1][warning][connection] [source/common/network/socket_option_impl.cc:33] Setting 41/67 option on socket failed: Protocol not available
[2024-10-18 10:04:27.651][1][error][config] [source/common/listener_manager/listener_manager_impl.cc:741] final pre-worker listener init for listener 'ingress_http' failed: cannot set post-listen socket option on socket: 0.0.0.0:8080
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:127] Caught Segmentation fault, suspect faulting address 0x0
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:111] Backtrace (use tools/stack_decode.py to get line numbers):
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:112] Envoy version: 7b8baff1758f0a584dcc3cb657b5032000bcb3d7/1.31.0/Clean/RELEASE/BoringSSL
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:114] Address mapping: 5c767e15f000-5c7680bb8000 /usr/local/bin/envoy
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #0: [0x78d986b04520]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #1: [0x5c767fa19ec8]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #2: [0x5c767fa1bd71]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #3: [0x5c767fa1f159]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #4: [0x5c767ffe5217]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #5: [0x5c767ffebb2a]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #6: [0x5c76801a66cb]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #7: [0x5c76801c727e]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #8: [0x5c767ff31982]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #9: [0x5c768014caac]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #10: [0x5c76801575bd]
[2024-10-18 10:04:27.653][1][critical][backtrace] [./source/server/backtrace.h:121] #11: [0x5c7680160144]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #12: [0x5c7680450e0e]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #13: [0x5c76804589b9]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #14: [0x5c768049fbe7]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #15: [0x5c768049f241]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #16: [0x5c768049c77f]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #17: [0x5c768049c71e]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #18: [0x5c768044bbf3]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #19: [0x5c768044b8c7]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #20: [0x5c7680154c5e]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #21: [0x5c76801558c5]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #22: [0x5c767ffd658d]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #23: [0x5c767ffd8885]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #24: [0x5c768038caf5]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #25: [0x5c7680333bdb]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #26: [0x5c768032fdf5]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #27: [0x5c7680339b96]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #28: [0x5c76803232f6]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #29: [0x5c7680324805]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #30: [0x5c7680596720]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #31: [0x5c7680595061]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #32: [0x5c767fb4fd11]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #33: [0x5c767faf29ba]
[2024-10-18 10:04:27.654][1][critical][backtrace] [./source/server/backtrace.h:121] #34: [0x5c767faf31de]
[2024-10-18 10:04:27.655][1][critical][backtrace] [./source/server/backtrace.h:121] #35: [0x5c767e15f14c]
[2024-10-18 10:04:27.655][1][critical][backtrace] [./source/server/backtrace.h:121] #36: [0x78d986aebd90]
AsyncClient 0x1660bf751600, stream_id_: 16581625266859155751
&stream_info_:
  StreamInfoImpl 0x1660bf7518c0, protocol_: 1, response_code_: 200, response_code_details_: via_upstream, attempt_count_: 1, health_check_request_: 0, getRouteName():   upstream_info_:
    UpstreamInfoImpl 0x1660bf57c6a8, upstream_connection_id_: 2375
Http2::ConnectionImpl 0x1660bf84f2d0, max_headers_kb_: 60, max_headers_count_: 100, per_stream_buffer_limit_: 268435456, allow_metadata_: 0, stream_error_on_invalid_http_messaging_: 0, is_outbound_flood_monitored_control_frame_: 0, dispatching_: 1, raised_goaway_: 0, pending_deferred_reset_streams_.size(): 0
&protocol_constraints_:
  ProtocolConstraints 0x1660bf84f340, outbound_frames_: 0, max_outbound_frames_: 10000, outbound_control_frames_: 0, max_outbound_control_frames_: 1000, consecutive_inbound_frames_with_empty_payload_: 0, max_consecutive_inbound_frames_with_empty_payload_: 1, opened_streams_: 9, inbound_priority_frames_: 0, max_inbound_priority_frames_per_stream_: 100, inbound_window_update_frames_: 58, outbound_data_frames_: 48, max_inbound_window_update_frames_per_data_frame_sent_: 10
Number of active streams: 9, current_stream_id_: 13 Dumping current stream:
stream:
  ConnectionImpl::StreamImpl 0x1660bf2e8480, stream_id_: 13, unconsumed_bytes_: 0, read_disable_count_: 0, local_end_stream_: 0, local_end_stream_sent_: 0, remote_end_stream_: 0, data_deferred_: 1, received_noninformational_headers_: 1, pending_receive_buffer_high_watermark_called_: 0, pending_send_buffer_high_watermark_called_: 0, reset_due_to_messaging_error_: 0, cookies_:   pending_trailers_to_encode_:   null
  absl::get<ResponseHeaderMapPtr>(headers_or_trailers_):   null
Dumping corresponding downstream request for upstream stream 13:
  UpstreamRequest 0x1660bf791000
  request_headers:
    ':method', 'POST'
    ':path', '/envoy.service.route.v3.RouteDiscoveryService/StreamRoutes'
    ':authority', 'contour'
    ':scheme', 'http'
    'te', 'trailers'
    'content-type', 'application/grpc'
    'x-envoy-internal', 'true'
    'x-forwarded-for', '10.244.1.10'
  FilterManager 0x1660bf57c7e0, state_.has_1xx_headers_: 0
  filter_manager_callbacks_.requestHeaders():
    ':method', 'POST'
    ':path', '/envoy.service.route.v3.RouteDiscoveryService/StreamRoutes'
    ':authority', 'contour'
    ':scheme', 'http'
    'te', 'trailers'
    'content-type', 'application/grpc'
    'x-envoy-internal', 'true'
    'x-forwarded-for', '10.244.1.10'
  filter_manager_callbacks_.requestTrailers():   null
  filter_manager_callbacks_.responseHeaders():   null
  filter_manager_callbacks_.responseTrailers():   null
  &streamInfo():
    StreamInfoImpl 0x1660bf7518c0, protocol_: 1, response_code_: 200, response_code_details_: via_upstream, attempt_count_: 1, health_check_request_: 0, getRouteName():     upstream_info_:
      UpstreamInfoImpl 0x1660bf57c6a8, upstream_connection_id_: 2375
current slice length: 1042 contents: "\0\0\0\0\0\r��\0�\0\0\0\0\0\r\0\0\0�\n$f369d119-4560-44c7-9cd9-8ec4b4b4a7b7�\n<type.googleapis.com/envoy.config.route.v3.RouteConfiguration�\n
      ingress_http�\nprotected.127-0-0-101.nip.ioprotected.127-0-0-101.nip.io▒�\n\n/\"�\n�\nenvoy.access_loggers.filex\n%\nio.projectcontour.kind\v▒\tHTTPProxy\n%\nio.projectcontour.name\v▒\tprotected\n(\no.projectcontour.namespace\t▒default▒ 2+\n)\nx-request-startt=%START_TIME(%s.%3f)%p\"<type.googleapis.com/envoy.config.route.v3.RouteConfiguration*1\0\0\0\0\0��\0\0\0\0\0\0\0\0\0\n$f369d119-4560-44c7-9cd9-8ec4b4b4a7b7�\n<type.googleapis.com/envoy.config.route.v3.RouteConfiguration�\n\"https/protected.127-0-0-101.nip.io�\nprotected.127-0-0-101.nip.ioprotected.127-0-0-101.nip.io▒�\n\n/\"�\n�\nenvoy.access_loggers.filex\n%\nio.projectcontour.kind\v▒\tHTTPProxy\n%\nio.projectcontour.name\v▒\tprotected\n(\no.projectcontour.namespace\t▒default0J\v\nrese\n!default/echoserver/443/684534031b2+\n)\nx-request-startt=%START_TIME(%s.%3f)%p\"<type.googleapis.com/envoy.config.route.v3.RouteConfiguration*1"
ConnectionImpl 0x1660bf784880, connecting_: 0, bind_error_: 0, state(): Open, read_buffer_limit_: 1048576
socket_:
  ListenSocketImpl 0x1660bf6b3d00, transport_protocol_:
  connection_info_provider_:
    ConnectionInfoSetterImpl 0x1660bf9ddc18, remote_address_: 10.96.91.10:8001, direct_remote_address_: 10.96.91.10:8001, local_address_: 10.244.1.10:40278, server_name_:







[2024-10-18 12:25:20.380][1][error][config] [source/common/listener_manager/listener_manager_impl.cc:741] final pre-worker listener init for listener 'ingress_http' failed: cannot set post-listen socket option on socket: 0.0.0.0:8080
[2024-10-18 12:25:20.380][1][debug][config] [source/common/listener_manager/listener_manager_impl.cc:892] begin remove listener: name=ingress_http
[2024-10-18 12:25:20.380][1][debug][config] [source/common/listener_manager/listener_impl.cc:924] removing warming listener: name=ingress_http, hash=6853495128709662139, tag=17, address=0.0.0.0:8080
[2024-10-18 12:25:20.380][1][debug][init] [source/common/init/watcher_impl.cc:31] Listener-local-init-watcher ingress_http destroyed
[2024-10-18 12:25:20.380][1][debug][init] [source/common/init/watcher_impl.cc:31] init manager RDS local-init-manager ingress_http destroyed
[2024-10-18 12:25:20.380][1][debug][init] [source/common/init/target_impl.cc:34] target RdsRouteConfigSubscription RDS local-init-target ingress_http destroyed
[2024-10-18 12:25:20.380][1][debug][init] [source/common/init/watcher_impl.cc:31] RDS local-init-watcher ingress_http destroyed
[2024-10-18 12:25:20.380][1][debug][init] [source/common/init/target_impl.cc:68] shared target RdsRouteConfigSubscription RDS init ingress_http destroyed
[2024-10-18 12:25:20.381][1][debug][router] [source/common/router/upstream_request.cc:493] [Tags: "ConnectionId":"0","StreamId":"1985046678044021176"] resetting pool request
[2024-10-18 12:25:20.381][1][debug][client] [source/common/http/codec_client.cc:159] [Tags: "ConnectionId":"257"] request reset
[2024-10-18 12:25:20.381][1][debug][pool] [source/common/conn_pool/conn_pool_base.cc:215] [Tags: "ConnectionId":"257"] destroying stream: 4 remaining
[2024-10-18 12:25:20.382][1][debug][init] [source/common/init/watcher_impl.cc:31] init manager Listener-local-init-manager ingress_http 6853495128709662139 destroyed
[2024-10-18 12:25:20.382][1][debug][init] [source/common/init/target_impl.cc:34] target Listener-init-target ingress_http destroyed
[2024-10-18 12:25:20.383][1][debug][config] [source/extensions/config_subscription/grpc/grpc_subscription_impl.cc:89] gRPC config for �(M��gleapis.com/envoy.config.route.v3.RouteConfiguration accepted with 1 resources with version 6678bb8f-4fd2-4d5b-8071-b3480e294930
[2024-10-18 12:25:20.383][1][critical][backtrace] [./source/server/backtrace.h:127] Caught Segmentation fault, suspect faulting address 0x0
[2024-10-18 12:25:20.383][1][critical][backtrace] [./source/server/backtrace.h:111] Backtrace (use tools/stack_decode.py to get line numbers):
[2024-10-18 12:25:20.383][1][critical][backtrace] [./source/server/backtrace.h:112] Envoy version: 7b8baff1758f0a584dcc3cb657b5032000bcb3d7/1.31.0/Clean/RELEASE/BoringSSL
[2024-10-18 12:25:20.383][1][critical][backtrace] [./source/server/backtrace.h:114] Address mapping: 5e6da9a93000-5e6dac4ec000 /usr/local/bin/envoy
[2024-10-18 12:25:20.384][1][critical][backtrace] [./source/server/backtrace.h:121] #0: [0x7769ecbd0520]
[2024-10-18 12:25:20.400][1][critical][backtrace] [./source/server/backtrace.h:119] #1: Envoy::Config::GrpcMuxImpl::processDiscoveryResources() [0x5e6dab34dec8]
[2024-10-18 12:25:20.407][1][critical][backtrace] [./source/server/backtrace.h:119] #2: Envoy::Config::GrpcMuxImpl::onDiscoveryResponse() [0x5e6dab34fd71]
[2024-10-18 12:25:20.412][1][critical][backtrace] [./source/server/backtrace.h:119] #3: Envoy::Grpc::AsyncStreamCallbacks<>::onReceiveMessageRaw() [0x5e6dab353159]
[2024-10-18 12:25:20.416][30][debug][conn_handler] [source/common/listener_manager/active_tcp_listener.cc:160] [Tags: "ConnectionId":"258"] new connection from 10.244.1.1:40476
[2024-10-18 12:25:20.416][30][debug][http] [source/common/http/conn_manager_impl.cc:385] [Tags: "ConnectionId":"258"] new stream
[2024-10-18 12:25:20.416][30][debug][http] [source/common/http/conn_manager_impl.cc:1135] [Tags: "ConnectionId":"258","StreamId":"6210591043362764417"] request headers complete (end_stream=true):
':authority', '10.244.1.25:8002'
':path', '/ready'
':method', 'GET'
'user-agent', 'kube-probe/1.30'
'accept', '*/*'
'connection', 'close'

[2024-10-18 12:25:20.416][30][debug][http] [source/common/http/conn_manager_impl.cc:1118] [Tags: "ConnectionId":"258","StreamId":"6210591043362764417"] request end stream
[2024-10-18 12:25:20.416][30][debug][connection] [./source/common/network/connection_impl.h:98] [Tags: "ConnectionId":"258"] current connecting state: false
[2024-10-18 12:25:20.416][30][debug][router] [source/common/router/router.cc:525] [Tags: "ConnectionId":"258","StreamId":"6210591043362764417"] cluster 'envoy-admin' match for URL '/ready'
[2024-10-18 12:25:20.416][30][debug][router] [source/common/router/router.cc:750] [Tags: "ConnectionId":"258","StreamId":"6210591043362764417"] router decoding headers:
':authority', '10.244.1.25:8002'
':path', '/ready'
':method', 'GET'
':scheme', 'http'
'user-agent', 'kube-probe/1.30'
'accept', '*/*'
'x-forwarded-proto', 'http'
'x-request-id', 'f7bf95d8-032b-4e4a-b57f-3e09b6e3ff2e'
'x-envoy-expected-rq-timeout-ms', '15000'

[2024-10-18 12:25:20.416][30][debug][pool] [source/common/conn_pool/conn_pool_base.cc:265] [Tags: "ConnectionId":"41"] using existing fully connected connection
[2024-10-18 12:25:20.416][30][debug][pool] [source/common/conn_pool/conn_pool_base.cc:182] [Tags: "ConnectionId":"41"] creating stream
[2024-10-18 12:25:20.416][30][debug][router] [source/common/router/upstream_request.cc:595] [Tags: "ConnectionId":"258","StreamId":"6210591043362764417"] pool ready
[2024-10-18 12:25:20.416][30][debug][client] [source/common/http/codec_client.cc:142] [Tags: "ConnectionId":"41"] encode complete
[2024-10-18 12:25:20.417][1][critical][backtrace] [./source/server/backtrace.h:119] #4: Envoy::Grpc::AsyncStreamImpl::onData() [0x5e6dab919217]
[2024-10-18 12:25:20.422][1][critical][backtrace] [./source/server/backtrace.h:119] #5: Envoy::Http::AsyncStreamImpl::encodeData() [0x5e6dab91fb2a]
[2024-10-18 12:25:20.426][1][critical][backtrace] [./source/server/backtrace.h:119] #6: Envoy::Router::UpstreamRequest::decodeData() [0x5e6dabada6cb]
[2024-10-18 12:25:20.431][1][critical][backtrace] [./source/server/backtrace.h:119] #7: Envoy::Http::FilterManager::encodeData() [0x5e6dabafb27e]
[2024-10-18 12:25:20.435][1][critical][backtrace] [./source/server/backtrace.h:119] #8: Envoy::Http::ResponseDecoderWrapper::decodeData() [0x5e6dab865982]
[2024-10-18 12:25:20.439][1][critical][backtrace] [./source/server/backtrace.h:119] #9: Envoy::Http::Http2::ConnectionImpl::StreamImpl::decodeData() [0x5e6daba80aac]
[2024-10-18 12:25:20.444][1][critical][backtrace] [./source/server/backtrace.h:119] #10: Envoy::Http::Http2::ConnectionImpl::onBeginData() [0x5e6daba8b5bd]
[2024-10-18 12:25:20.448][1][critical][backtrace] [./source/server/backtrace.h:119] #11: Envoy::Http::Http2::ConnectionImpl::Http2Visitor::OnDataForStream() [0x5e6daba94144]
[2024-10-18 12:25:20.452][1][critical][backtrace] [./source/server/backtrace.h:119] #12: http2::adapter::OgHttp2Session::OnStreamFrameData() [0x5e6dabd84e0e]
[2024-10-18 12:25:20.457][1][critical][backtrace] [./source/server/backtrace.h:119] #13: http2::Http2TraceLogger::OnStreamFrameData() [0x5e6dabd8c9b9]
[2024-10-18 12:25:20.462][1][critical][backtrace] [./source/server/backtrace.h:119] #14: http2::DataPayloadDecoder::StartDecodingPayload() [0x5e6dabdd3be7]
[2024-10-18 12:25:20.466][1][critical][backtrace] [./source/server/backtrace.h:119] #15: http2::Http2FrameDecoder::StartDecodingPayload() [0x5e6dabdd3241]
[2024-10-18 12:25:20.470][1][critical][backtrace] [./source/server/backtrace.h:119] #16: http2::Http2DecoderAdapter::ProcessInputFrame() [0x5e6dabdd077f]
[2024-10-18 12:25:20.474][1][critical][backtrace] [./source/server/backtrace.h:119] #17: http2::Http2DecoderAdapter::ProcessInput() [0x5e6dabdd071e]
[2024-10-18 12:25:20.479][1][critical][backtrace] [./source/server/backtrace.h:119] #18: http2::adapter::OgHttp2Session::ProcessBytesImpl() [0x5e6dabd7fbf3]
[2024-10-18 12:25:20.483][1][critical][backtrace] [./source/server/backtrace.h:119] #19: http2::adapter::OgHttp2Session::ProcessBytes() [0x5e6dabd7f8c7]
[2024-10-18 12:25:20.487][1][critical][backtrace] [./source/server/backtrace.h:119] #20: Envoy::Http::Http2::ConnectionImpl::dispatch() [0x5e6daba88c5e]
[2024-10-18 12:25:20.491][1][critical][backtrace] [./source/server/backtrace.h:119] #21: Envoy::Http::Http2::ConnectionImpl::dispatch() [0x5e6daba898c5]
[2024-10-18 12:25:20.496][1][critical][backtrace] [./source/server/backtrace.h:119] #22: Envoy::Http::CodecClient::onData() [0x5e6dab90a58d]
[2024-10-18 12:25:20.500][1][critical][backtrace] [./source/server/backtrace.h:119] #23: Envoy::Http::CodecClient::CodecReadFilter::onData() [0x5e6dab90c885]
[2024-10-18 12:25:20.504][1][critical][backtrace] [./source/server/backtrace.h:119] #24: Envoy::Network::FilterManagerImpl::onContinueReading() [0x5e6dabcc0af5]
[2024-10-18 12:25:20.509][1][critical][backtrace] [./source/server/backtrace.h:119] #25: Envoy::Network::ConnectionImpl::onReadReady() [0x5e6dabc67bdb]
[2024-10-18 12:25:20.514][1][critical][backtrace] [./source/server/backtrace.h:119] #26: Envoy::Network::ConnectionImpl::onFileEvent() [0x5e6dabc63df5]
[2024-10-18 12:25:20.518][1][critical][backtrace] [./source/server/backtrace.h:119] #27: std::__1::__function::__func<>::operator()() [0x5e6dabc6db96]
[2024-10-18 12:25:20.523][1][critical][backtrace] [./source/server/backtrace.h:119] #28: std::__1::__function::__func<>::operator()() [0x5e6dabc572f6]
[2024-10-18 12:25:20.527][1][critical][backtrace] [./source/server/backtrace.h:119] #29: Envoy::Event::FileEventImpl::mergeInjectedEventsAndRunCb() [0x5e6dabc58805]
[2024-10-18 12:25:20.531][1][critical][backtrace] [./source/server/backtrace.h:119] #30: event_process_active_single_queue [0x5e6dabeca720]
[2024-10-18 12:25:20.536][1][critical][backtrace] [./source/server/backtrace.h:119] #31: event_base_loop [0x5e6dabec9061]
[2024-10-18 12:25:20.540][1][critical][backtrace] [./source/server/backtrace.h:119] #32: Envoy::Server::InstanceBase::run() [0x5e6dab483d11]
[2024-10-18 12:25:20.544][1][critical][backtrace] [./source/server/backtrace.h:119] #33: Envoy::MainCommonBase::run() [0x5e6dab4269ba]
[2024-10-18 12:25:20.548][1][critical][backtrace] [./source/server/backtrace.h:119] #34: Envoy::MainCommon::main() [0x5e6dab4271de]
[2024-10-18 12:25:20.553][1][critical][backtrace] [./source/server/backtrace.h:119] #35: main [0x5e6da9a9314c]
[2024-10-18 12:25:20.553][1][critical][backtrace] [./source/server/backtrace.h:121] #36: [0x7769ecbb7d90]
AsyncClient 0x17ccbf74e600, stream_id_: 1985046678044021176
&stream_info_:
  StreamInfoImpl 0x17ccbf74e8c0, protocol_: 1, response_code_: 200, response_code_details_: via_upstream, attempt_count_: 1, health_check_request_: 0, getRouteName():   upstream_info_:
    UpstreamInfoImpl 0x17ccbfc43e48, upstream_connection_id_: 257
Http2::ConnectionImpl 0x17ccbf2ee790, max_headers_kb_: 60, max_headers_count_: 100, per_stream_buffer_limit_: 268435456, allow_metadata_: 0, stream_error_on_invalid_http_messaging_: 0, is_outbound_flood_monitored_control_frame_: 0, dispatching_: 1, raised_goaway_: 0, pending_deferred_reset_streams_.size(): 0
&protocol_constraints_:
  ProtocolConstraints 0x17ccbf2ee800, outbound_frames_: 0, max_outbound_frames_: 10000, outbound_control_frames_: 0, max_outbound_control_frames_: 1000, consecutive_inbound_frames_with_empty_payload_: 0, max_consecutive_inbound_frames_with_empty_payload_: 1, opened_streams_: 5, inbound_priority_frames_: 0, max_inbound_priority_frames_per_stream_: 100, inbound_window_update_frames_: 15, outbound_data_frames_: 17, max_inbound_window_update_frames_per_data_frame_sent_: 10
Number of active streams: 5, current_stream_id_: 5 Dumping current stream:
stream:
  ConnectionImpl::StreamImpl 0x17ccbf5546c0, stream_id_: 5, unconsumed_bytes_: 0, read_disable_count_: 0, local_end_stream_: 0, local_end_stream_sent_: 0, remote_end_stream_: 0, data_deferred_: 1, received_noninformational_headers_: 1, pending_receive_buffer_high_watermark_called_: 0, pending_send_buffer_high_watermark_called_: 0, reset_due_to_messaging_error_: 0, cookies_:   pending_trailers_to_encode_:   null
  absl::get<ResponseHeaderMapPtr>(headers_or_trailers_):   null
Dumping corresponding downstream request for upstream stream 5:
  UpstreamRequest 0x17ccbf70f000
  request_headers:
    ':method', 'POST'
    ':path', '/envoy.service.route.v3.RouteDiscoveryService/StreamRoutes'
    ':authority', 'contour'
    ':scheme', 'http'
    'te', 'trailers'
    'content-type', 'application/grpc'
    'x-envoy-internal', 'true'
    'x-forwarded-for', '10.244.1.25'
  FilterManager 0x17ccbf1d1260, state_.has_1xx_headers_: 0
  filter_manager_callbacks_.requestHeaders():
    ':method', 'POST'
    ':path', '/envoy.service.route.v3.RouteDiscoveryService/StreamRoutes'
    ':authority', 'contour'
    ':scheme', 'http'
    'te', 'trailers'
    'content-type', 'application/grpc'
    'x-envoy-internal', 'true'
    'x-forwarded-for', '10.244.1.25'
  filter_manager_callbacks_.requestTrailers():   null
  filter_manager_callbacks_.responseHeaders():   null
  filter_manager_callbacks_.responseTrailers():   null
  &streamInfo():
    StreamInfoImpl 0x17ccbf74e8c0, protocol_: 1, response_code_: 200, response_code_details_: via_upstream, attempt_count_: 1, health_check_request_: 0, getRouteName():     upstream_info_:
      UpstreamInfoImpl 0x17ccbfc43e48, upstream_connection_id_: 257
current slice length: 1921 contents: "\0\\0\0\0\0\0\0\0��\0\\0\0\0\0\0\0�\0\0\\0\0\0\0\0\0\0��\0\\0\0\0\0\0\0�\0\0\\0\0\0\0\0\0\0�.\0\\0\0\0\0\t\0\0�\0\0\\0\0\0\0\0\0�)\0\0\0\0\0��\0+\0\0\0\0\0\0\0\0&\n$6678bb8f-4fd2-4d5b-8071-b3480e294930�\t\nDtype.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret\n▒default/ingress/f07d62545c\n��-----BEGIN CERTIFICATE-----\nMIICRDCCAemgAwIBAgIIF6mYTr1NpM8wCgYIKoZIzj0EAwIwGzEZMBcGA1UEAxMQ\nZXh0ZXJuYWwtcm9vdC1jYTAeFw0yNDAxMTIxMjA5NTdaFw0yNTAxMTExMjA5NTda\nMBIxEDAOBgNVBAMTB2luZ3Jlc3MwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAATu\nEmOWDElgr2meXcn+hcOfApTL+US4GXIlwPzJv2QOZEkQyWGt5KL9a1l3nx8clc/4\nMqABJ27pBhK2jblt5d0zo4IBHjCCARowDgYDVR0PAQH/BAQDAgWgMB8GA1UdIwQY\nMBaAFCOdQ3dqv8WjRd58N1YSNkwwy4YAMIHmBgNVHREEgd4wgduCGGhvc3QxLjEy\nNy0wLTAtMTAxLm5pcC5pb4IYaG9zdDIuMTI3LTAtMC0xMDEubmlwLmlvghxwcm90\nZWN0ZWQuMTI3LTAtMC0xMDEubmlwLmlvgh1wcm90ZWN0ZWQyLjEyNy0wLTAtMTAx\nLm5pcC5pb4IncHJvdGVjdGVkLWJhc2ljLWF1dGguMTI3LTAtMC0xMDEubmlwLmlv\ngiJwcm90ZWN0ZWQtb2F1dGguMTI3LTAtMC0xMDEubmlwLmlvghtrZXljbG9hay4x\nMjctMC0wLTEwMS5uaXAuaW8wCgYIKoZIzj0EAwIDSQAwRgIhAO2IbIfUO3KrZ05a\nG3zh0ROzWltIt7u2LRkr41bhhYtAAiEAqbQma/Al+XSp3Lvv65I1JcZNsx9upt6A\np18UqwMjB7k=\n-----END CERTIFICATE-----\n��-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgVh4gHp/UuTbaDuO7\nIO/4O0PQt8aVCehOQLNYw1GqgCehRANCAATuEmOWDElgr2meXcn+hcOfApTL+US4\nGXIlwPzJv2QOZEkQyWGt5KL9a1l3nx8clc/4MqABJ27pBhK2jblt5d0z\n-----END PRIVATE KEY-----\n\"Dtype.googleapis.com/envoy.extensions.transport_sockets.tls.v3.Secret*1\0\0\0\0\0��\0�\0\0\0\0\0\0\0\0�\n$6678bb8f-4fd2-4d5b-8071-b3480e294930�\n<type.googleapis.com/envoy.config.route.v3.RouteConfiguration�\n
                                                                                                     ingress_http�\nprotected.127-0-0-101.nip.ioprotected.127-0-0-101.nip.io▒�\n\n/\"�\n�\nenvoy.access_loggers.filex\n%\nio.projectcontour.kind\v▒\tHTTPProxy\n%\nio.projectcontour.name\v▒\tprotected\n(\no.projectcontour.namespace\t▒default▒ 2+\n)\nx-request-startt=%START_TIME(%s.%3f)%p\"<type.googleapis.com/envoy.config.route.v3.RouteConfiguration*1"
ConnectionImpl 0x17ccbf6bd600, connecting_: 0, bind_error_: 0, state(): Open, read_buffer_limit_: 1048576
socket_:
  ListenSocketImpl 0x17ccbf146f00, transport_protocol_:
  connection_info_provider_:
    ConnectionInfoSetterImpl 0x17ccbf9ddc18, remote_address_: 10.96.91.10:8001, direct_remote_address_: 10.96.91.10:8001, local_address_: 10.244.1.25:47384, server_name_:












[2024-10-18 12:51:54.612][1][critical][backtrace] [./source/server/backtrace.h:119] #1: Envoy::Config::GrpcMuxImpl::processDiscoveryResources() [0x626b1b304ec8]->[0x18baec8] /opt/llvm/bin/../include/c++/v1/vector:505
[2024-10-18 12:51:54.618][1][critical][backtrace] [./source/server/backtrace.h:119] #2: Envoy::Config::GrpcMuxImpl::onDiscoveryResponse() [0x626b1b306d71]->[0x18bcd71] bazel-out/k8-opt/bin/external/com_google_protobuf/src/google/protobuf/_virtual_includes/protobuf/google/protobuf/map_field.h:588
[2024-10-18 12:51:54.623][1][critical][backtrace] [./source/server/backtrace.h:119] #3: Envoy::Grpc::AsyncStreamCallbacks<>::onReceiveMessageRaw() [0x626b1b30a159]->[0x18c0159] /opt/llvm/bin/../include/c++/v1/string:2231
[2024-10-18 12:51:54.627][1][critical][backtrace] [./source/server/backtrace.h:119] #4: Envoy::Grpc::AsyncStreamImpl::onData() [0x626b1b8d0217]->[0x1e86217] /opt/llvm/bin/../include/c++/v1/new:255
[2024-10-18 12:51:54.632][1][critical][backtrace] [./source/server/backtrace.h:119] #5: Envoy::Http::AsyncStreamImpl::encodeData() [0x626b1b8d6b2a]->[0x1e8cb2a] /opt/llvm/bin/../include/c++/v1/new:255
[2024-10-18 12:51:54.636][1][critical][backtrace] [./source/server/backtrace.h:119] #6: Envoy::Router::UpstreamRequest::decodeData() [0x626b1ba916cb]->[0x20476cb] external/com_envoyproxy_protoc_gen_validate/validate/validate.h:81
[2024-10-18 12:51:54.641][1][critical][backtrace] [./source/server/backtrace.h:119] #7: Envoy::Http::FilterManager::encodeData() [0x626b1bab227e]->[0x206827e] /opt/llvm/bin/../include/c++/v1/__memory/allocator.h:101
[2024-10-18 12:51:54.645][1][critical][backtrace] [./source/server/backtrace.h:119] #8: Envoy::Http::ResponseDecoderWrapper::decodeData() [0x626b1b81c982]->[0x1dd2982] /opt/llvm/bin/../include/c++/v1/new:245
