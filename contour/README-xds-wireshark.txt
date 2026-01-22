


sudo nsenter --net --target $(kindps contour envoy -o json | jq -r .[0].pids[0].pid) wireshark -i any -f "port 8001" \
  -o tls.keylog_file:/tmp/sslkeylog.log \
  -o protobuf.preload_protos:true \
  -k


sudo tee /root/.config/wireshark/protobuf_search_paths <<EOF
"/home/tsaarni/package/envoy-xds-proto-deps/envoy/api/","FALSE"
"/home/tsaarni/package/envoy-xds-proto-deps/googleapis/","FALSE"
"/home/tsaarni/package/envoy-xds-proto-deps/udpa","FALSE"
"/home/tsaarni/package/envoy-xds-proto-deps/protoc-gen-validate/","FALSE"
EOF

### THIS DOES NOT WORK
###  -o protobuf.search_paths:/home/tsaarni/package/envoy-xds-proto-deps/envoy/api/,/home/tsaarni/package/envoy-xds-proto-deps/googleapis/,/home/tsaarni/package/envoy-xds-proto-deps/udpa,/home/tsaarni/package/envoy-xds-proto-deps/protoc-gen-validate/ \
###  -o "protobuf.protobuf_search_paths:/home/tsaarni/package/envoy-xds-proto-deps/envoy/api/,/home/tsaarni/package/envoy-xds-proto-deps/googleapis/,/home/tsaarni/package/envoy-xds-proto-deps/udpa,/home/tsaarni/package/envoy-xds-proto-deps/protoc-gen-validate/" \




/home/tsaarni/work/envoy/api/envoy/service/discovery/v3/discovery.proto




/home/tsaarni/work/envoy/api/envoy/service/discovery/v3/discovery.proto


cd /home/tsaarni/work/envoy/api
protoc --proto_path=. --include_imports --include_source_info --descriptor_set_out=envoy.pb envoy/service/discovery/v3/discovery.proto

diff --git a/cmd/contour/servecontext.go b/cmd/contour/servecontext.go
index 9ecaa3b4..9e83d1b9 100644
--- a/cmd/contour/servecontext.go
+++ b/cmd/contour/servecontext.go
@@ -184,6 +184,16 @@ func tlsconfig(log logrus.FieldLogger, contourXDSTLS *contour_v1alpha1.TLS) *tls
                log.WithError(err).Fatal("failed to verify TLS flags")
        }

+       // Set KeyLogWriter if SSLKEYLOGFILE is set.
+       var w *os.File
+       if keyLogFile := os.Getenv("SSLKEYLOGFILE"); keyLogFile != "" {
+               log.WithField("SSLKEYLOGFILE", keyLogFile).Info("SSL key logging enabled")
+               w, err = os.OpenFile(os.Getenv("SSLKEYLOGFILE"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
+               if err != nil {
+                       log.WithError(err).Fatal("failed to open SSL key log file")
+               }
+       }
+
        // Define a closure that lazily loads certificates and key at TLS handshake
        // to ensure that latest certificates are used in case they have been rotated.
        loadConfig := func() (*tls.Config, error) {
@@ -211,6 +221,7 @@ func tlsconfig(log logrus.FieldLogger, contourXDSTLS *contour_v1alpha1.TLS) *tls
                        ClientCAs:    certPool,
                        MinVersion:   tls.VersionTLS13,
                        NextProtos:   []string{http2.NextProtoTLS},
+                       KeyLogWriter: w,
                }, nil
        }


sudo apt install protobuf-compiler

git clone --depth 1 https://github.com/envoyproxy/envoy.git
git clone --depth 1 https://github.com/googleapis/googleapis.git
git clone --depth 1 https://github.com/cncf/udpa.git
git clone --depth 1 https://github.com/envoyproxy/protoc-gen-validate.git
git clone --depth 1 https://github.com/cncf/xds.git
git clone --depth 1 https://github.com/prometheus/client_model.git prometheus
git clone --depth 1 https://github.com/open-telemetry/opentelemetry-proto.git

protoc \
  -Ienvoy/api \
  -Igoogleapis \
  -Iudpa \
  -Iprotoc-gen-validate \
  --include_imports \
  --include_source_info \
  --descriptor_set_out=envoy.pb \
  envoy/service/discovery/v3/discovery.proto

# create python bindings
rm -rf python
mkdir python
find xds envoy/api udpa/udpa protoc-gen-validate -name '*.proto' | xargs \
    protoc \
    -Ienvoy/api \
    -Igoogleapis \
    -Ixds \
    -Iudpa \
    -Iprotoc-gen-validate \
    -Iprometheus \
    -Iopentelemetry-proto \
    --include_imports \
    --include_source_info \
    --descriptor_set_out=envoy.pb \
    --python_out=python



  envoy/service/discovery/v3/discovery.proto \
  envoy/config/core/v3/base.proto
