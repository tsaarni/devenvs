

[2024-07-04T10:37:37.011+02:00][14][debug][config] [source/extensions/config_subscription/grpc/grpc_subscription_impl.cc:89] gRPC config for type.googleapis.com/envoy.config.route.v3.RouteConfigura
tion accepted with 1 resources with version 10342e68-5534-400b-96dc-cdcf4111afcb
[2024-07-04T10:37:37.013+02:00][14][critical][backtrace] [./source/server/backtrace.h:104] Caught Segmentation fault, suspect faulting address 0x0
[2024-07-04T10:37:37.013+02:00][14][critical][backtrace] [./source/server/backtrace.h:91] Backtrace (use tools/stack_decode.py to get line numbers):
[2024-07-04T10:37:37.013+02:00][14][critical][backtrace] [./source/server/backtrace.h:92] Envoy version: 4fda4d79d06e1bd59e591be3f348223495083648/1.29.1/Modified/DEBUG/BoringSSL
[2024-07-04T10:37:37.097+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #0: Envoy::SignalAction::sigHandler() [0x556151a693c8]
[2024-07-04T10:37:37.097+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #1: __restore_rt [0x7f4d253f9910]
[2024-07-04T10:37:37.155+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #2: Envoy::Config::TtlManager::ScopedTtlUpdate::~ScopedTtlUpdate() [0x55614fb6cb37]
[2024-07-04T10:37:37.209+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #3: Envoy::Config::GrpcMuxImpl::processDiscoveryResources() [0x55614fb627ba]
[2024-07-04T10:37:37.264+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #4: Envoy::Config::GrpcMuxImpl::onDiscoveryResponse() [0x55614fb6570f]
[2024-07-04T10:37:37.316+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #5: Envoy::Config::GrpcStream<>::onReceiveMessage() [0x55614fb70db3]
[2024-07-04T10:37:37.372+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #6: Envoy::Grpc::AsyncStreamCallbacks<>::onReceiveMessageRaw() [0x55614fb70c5f]
[2024-07-04T10:37:37.428+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #7: Envoy::Grpc::AsyncStreamImpl::onData() [0x5561509dd24d]
[2024-07-04T10:37:37.483+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #8: Envoy::Http::AsyncStreamImpl::encodeData() [0x5561509ea1ff]
[2024-07-04T10:37:37.535+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #9: Envoy::Router::Filter::onUpstreamData() [0x556150f81982]
[2024-07-04T10:37:37.587+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #10: Envoy::Router::UpstreamRequest::decodeData() [0x556150fa3fbf]
[2024-07-04T10:37:37.639+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #11: Envoy::Router::UpstreamRequestFilterManagerCallbacks::encodeData() [0x556150fb1e42]
[2024-07-04T10:37:37.692+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #12: Envoy::Http::FilterManager::encodeData() [0x5561510f9ff3]
[2024-07-04T10:37:37.745+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #13: Envoy::Http::ActiveStreamDecoderFilter::encodeData() [0x5561510f865a]
[2024-07-04T10:37:37.800+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #14: Envoy::Router::UpstreamCodecFilter::CodecBridge::decodeData() [0x556150fba40e]
[2024-07-04T10:37:37.852+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #15: Envoy::Http::ResponseDecoderWrapper::decodeData() [0x5561507dd1cc]
[2024-07-04T10:37:37.905+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #16: Envoy::Http::Http2::ConnectionImpl::StreamImpl::decodeData() [0x556150e3377c]
[2024-07-04T10:37:37.957+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #17: Envoy::Http::Http2::ConnectionImpl::onFrameReceived() [0x556150e4a118]
[2024-07-04T10:37:38.017+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #18: Envoy::Http::Http2::ConnectionImpl::Http2Callbacks::Http2Callbacks()::$_24::operator()() [0x556150e5c72d]
[2024-07-04T10:37:38.086+02:00][14][critical][backtrace] [./source/server/backtrace.h:96] #19: Envoy::Http::Http2::ConnectionImpl::Http2Callbacks::Http2Callbacks()::$_24::__invoke() [0x556150e5c6f5]


# issue with similar backtrace
# https://github.com/envoyproxy/envoy/issues/20131 
#
# It is unresolved, and seems to imply that some actions that the control plane could cause the deletion of TtlManager
# before ScopedTtlUpdate is destroyed, causing the null pointer reference


 

# in our case before crash there was failed attempt to set DSCP socket option for IPv6
# which might have caused it

2024-07-04T10:37:37.006+02:00][14][warning][connection] [source/common/network/socket_option_impl.cc:33] Setting 41/67 option on socket failed: Protocol not available
[2024-07-04T10:37:37.006+02:00][14][error][config] [source/common/listener_manager/listener_manager_impl.cc:737] final pre-worker listener init for listener \'ingress_http\' failed: cannot set post-listen socket option on socket: 127.0.0.1:8080
[2024-07-04T10:37:37.006+02:00][14][debug][config] [source/common/listener_manager/listener_manager_impl.cc:888] begin remove listener: name=ingress_http
[2024-07-04T10:37:37.006+02:00][14][debug][config] [source/common/listener_manager/listener_impl.cc:894] removing warming listener: name=ingress_http, hash=17911480236002576541, tag=6, address=127.0.0.1:8080
[2024-07-04T10:37:37.006+02:00][14][debug][init] [source/common/init/watcher_impl.cc:31] Listener-local-init-watcher ingress_http destroyed
[2024-07-04T10:37:37.006+02:00][14][debug][init] [source/common/init/watcher_impl.cc:31] init manager RDS local-init-manager ingress_http destroyed
[2024-07-04T10:37:37.006+02:00][14][debug][init] [source/common/init/target_impl.cc:34] target RdsRouteConfigSubscription RDS local-init-target ingress_http destroyed
[2024-07-04T10:37:37.006+02:00][14][debug][init] [source/common/init/watcher_impl.cc:31] RDS local-init-watcher ingress_http destroyed
[2024-07-04T10:37:37.006+02:00][14][debug][init] [source/common/init/target_impl.cc:68] shared target RdsRouteConfigSubscription RDS init ingress_http destroyed


# Crash does not happen every time socket option fails

grep -E "(Protocol not available|Caught Segmentation fault)" envoy.log

[2024-07-04T09:59:31.462+02:00][14][warning][connection] [source/common/network/socket_option_impl.cc:33] Setting 41/67 option on socket failed: Protocol not available
[2024-07-04T10:37:37.006+02:00][14][warning][connection] [source/common/network/socket_option_impl.cc:33] Setting 41/67 option on socket failed: Protocol not available
[2024-07-04T10:37:37.013+02:00][14][critical][backtrace] [./source/server/backtrace.h:104] Caught Segmentation fault, suspect faulting address 0x0
[2024-07-04T10:37:43.641+02:00][14][warning][connection] [source/common/network/socket_option_impl.cc:33] Setting 41/67 option on socket failed: Protocol not available



# it seemed some certificates were rotated nearly at the same time
# maybe these combined could trigger the problem?

