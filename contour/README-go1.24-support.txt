
https://github.com/projectcontour/contour/pull/6987


~/work/devenvs/contour/apps/update-go.sh 1.24


make generate

07 Apr 25 16:52 EEST INF Starting mockery dry-run=false version=v2.43.2
07 Apr 25 16:52 EEST INF Using config: /home/tsaarni/work/contour/.mockery.yaml dry-run=false version=v2.43.2
2025/04/07 16:52:39 internal error: package "context" without types was imported from "github.com/projectcontour/contour/internal/leadership"



go get github.com/vektra/mockery/v2@v2.53.3
go mod tidy


make generate

panic: interface conversion: types.Type is *types.Alias, not *types.Named


go get sigs.k8s.io/controller-tools@v0.17.3
go mod tidy

make generate

Generating API documentation...
I0407 17:04:45.574901  401089 main.go:129] parsing go packages in directory github.com/projectcontour/contour/apis/projectcontour
W0407 17:04:54.443549  401089 parse.go:863] Making unsupported type entry "internal/abi.mapType" for: &types.Alias{obj:(*types.TypeName)(0xc0008cd540), orig:(*types.Alias)(0xc0009494c0), tparams:(*types.TypeParamList)(nil), targs:(*types.TypeList)(nil), fromRHS:(*types.Named)(0xc0001d9c70), actual:(*types.Named)(0xc0001d9c70)}
W0407 17:04:54.444213  401089 parse.go:863] Making unsupported type entry "any" for: &types.Alias{obj:(*types.TypeName)(0xc0000f49b0), orig:(*types.Alias)(0xc0000b2300), tparams:(*types.TypeParamList)(nil), targs:(*types.TypeList)(nil), fromRHS:(*types.Interface)(0xc0000f4910), actual:(*types.Interface)(0xc0000f4910)}
W0407 17:04:54.444489  401089 parse.go:863] Making unsupported type entry "reflect.uncommonType" for: &types.Alias{obj:(*types.TypeName)(0xc002251130), orig:(*types.Alias)(0xc0029dc440), tparams:(*types.TypeParamList)(nil), targs:(*types.TypeList)(nil), fromRHS:(*types.Named)(0xc0001da2a0), actual:(*types.Named)(0xc0001da2a0)}
W0407 17:04:54.444992  401089 parse.go:863] Making unsupported type entry "github.com/projectcontour/contour/apis/projectcontour/v1.Condition" for: &types.Alias{obj:(*types.TypeName)(0xc00b9908c0), orig:(*types.Alias)(0xc00bc52780), tparams:(*types.TypeParamList)(nil), targs:(*types.TypeList)(nil), fromRHS:(*types.Named)(0xc00b9d7730), actual:(*types.Named)(0xc00b9d7730)}
W0407 17:04:54.445002  401089 parse.go:863] Making unsupported type entry "github.com/projectcontour/contour/apis/projectcontour/v1.ConditionStatus" for: &types.Alias{obj:(*types.TypeName)(0xc00b990870), orig:(*types.Alias)(0xc00bc525c0), tparams:(*types.TypeParamList)(nil), targs:(*types.TypeList)(nil), fromRHS:(*types.Named)(0xc00b9d7500), actual:(*types.Named)(0xc00b9d7500)}
W0407 17:04:54.449912  401089 parse.go:863] Making unsupported type entry "k8s.io/api/core/v1.ServiceExternalTrafficPolicyType" for: &types.Alias{obj:(*types.TypeName)(0xc00e49f3b0), orig:(*types.Alias)(0xc00e6ea780), tparams:(*types.TypeParamList)(nil), targs:(*types.TypeList)(nil), fromRHS:(*types.Named)(0xc00e669110), actual:(*types.Named)(0xc00e669110)}
I0407 17:04:54.453354  401089 main.go:231] using package=github.com/projectcontour/contour/apis/projectcontour/v1
I0407 17:04:54.453362  401089 main.go:231] using package=github.com/projectcontour/contour/apis/projectcontour/v1alpha1
F0407 17:04:54.457098  401089 main.go:503] type github.com/projectcontour/contour/apis/projectcontour/v1.Condition has kind=Unsupported which is unhandled
make: *** [Makefile:267: generate-api-docs] Error 255


go get github.com/ahmetb/gen-crd-api-reference-docs@71fefeed8910
go mod tidy

make generate




make lint

internal/build/version.go:40:14: undefined: yaml (typecheck)
        out, err := yaml.Marshal(buildInfo)
                    ^
internal/httpsvc/http.go:55:8: svc.WithError undefined (type *Service has no field or method WithError) (typecheck)
                        svc.WithError(err).Error("terminated HTTP server with error")
                            ^
internal/httpsvc/http.go:57:8: svc.Info undefined (type *Service has no field or method Info) (typecheck)
                        svc.Info("stopped HTTP server")
                            ^
internal/httpsvc/http.go:74:7: svc.Fatal undefined (type *Service has no field or method Fatal) (typecheck)
                svc.Fatal("you must supply at least server certificate and key TLS parameters or none of them")
                    ^
internal/httpsvc/http.go:99:7: svc.WithField undefined (type *Service has no field or method WithField) (typecheck)
                svc.WithField("address", s.Addr).Info("started HTTPS server")
                    ^
internal/httpsvc/http.go:103:6: svc.WithField undefined (type *Service has no field or method WithField) (typecheck)
        svc.WithField("address", s.Addr).Info("started HTTP server")
            ^
.....



~/work/devenvs/contour/apps/update-golangci-lint.sh 1.64


make lint

tools.go:7:2: import "github.com/ahmetb/gen-crd-api-reference-docs" is a program, not an importable package (typecheck)
        _ "github.com/ahmetb/gen-crd-api-reference-docs"
        ^



sed -i '/args: --build-tags=/s/e2e,conformance,tools,gcp,oidc,none/e2e,conformance,gcp,oidc,none/' .github/workflows/prbuild.yaml
sed -i '/golangci-lint run --build-tags=/s/e2e,conformance,tools,gcp,oidc,none/e2e,conformance,gcp,oidc,none/' Makefile


make lint


internal/k8s/status_test.go:24:2: could not import sigs.k8s.io/controller-runtime/pkg/client/fake (-: # sigs.k8s.io/controller-runtime/pkg/client/fake
../../go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.7/pkg/client/fake/client.go:861:29: cannot use c.tracker (variable of struct type versionedTracker) as "k8s.io/client-go/testing".ObjectTracker value in argument to dryPatch: versionedTracker does not implement "k8s.io/client-go/testing".ObjectTracker (wrong type for method Create)
                have Create("k8s.io/apimachinery/pkg/runtime/schema".GroupVersionResource, "k8s.io/apimachinery/pkg/runtime".Object, string) error
                want Create("k8s.io/apimachinery/pkg/runtime/schema".GroupVersionResource, "k8s.io/apimachinery/pkg/runtime".Object, string, ..."k8s.io/apimachinery/pkg/apis/meta/v1".CreateOptions) error
../../go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.7/pkg/client/fake/client.go:875:37: cannot use c.tracker (variable of struct type versionedTracker) as "k8s.io/client-go/testing".ObjectTracker value in argument to testing.ObjectReaction: versionedTracker does not implement "k8s.io/client-go/testing".ObjectTracker (wrong type for method Create)
                have Create("k8s.io/apimachinery/pkg/runtime/schema".GroupVersionResource, "k8s.io/apimachinery/pkg/runtime".Object, string) error
                want Create("k8s.io/apimachinery/pkg/runtime/schema".GroupVersionResource, "k8s.io/apimachinery/pkg/runtime".Object, string, ..."k8s.io/apimachinery/pkg/apis/meta/v1".CreateOptions) error) (typecheck)
        "sigs.k8s.io/controller-runtime/pkg/client/fake"
        ^
...



go get sigs.k8s.io/controller-runtime@v0.20.4
go mod tidy



pkg/certs/certgen.go:98:38: G115: integer overflow conversion uint -> int64 (gosec)
        expiry := now.Add(24 * time.Duration(uint32OrDefault(config.Lifetime, DefaultCertificateLifetime)) * time.Hour)
                                            ^
pkg/certs/certgen_test.go:53:65: G115: integer overflow conversion uint -> int64 (gosec)
                                currentTime = currentTime.Add(24 * time.Hour * time.Duration(tc.config.Lifetime)).Add(-time.Hour)
                                                                                            ^
....


cat >> .golangci.yml << EOF
  - linters: ["gosec"]
    text: "G115"
  - linters: ["revive"]
    text: "redefines-builtin-id"
  - linters: ["gosimple"]
    text: "S1009"
  - linters: ["staticcheck"]
    text: "SA1006"
  - linters: ["testifylint"]
    text: "formatter"
  - linters: ["testifylint"]
    text: "negative-positive"
EOF



make lint



make check


--- FAIL: TestServeContextCertificateHandling (0.02s)
    --- FAIL: TestServeContextCertificateHandling/rotating_server_credentials_returns_new_server_cert (0.00s)
        servecontext_test.go:231: Unexpected result when connecting to the server: EOF
    --- FAIL: TestServeContextCertificateHandling/rotating_server_credentials_again_to_ensure_rotation_can_be_repeated (0.01s)
        servecontext_test.go:231: Unexpected result when connecting to the server: EOF
    --- FAIL: TestServeContextCertificateHandling/successful_TLS_connection_established (0.00s)
        servecontext_test.go:231: Unexpected result when connecting to the server: EOF




diff --git a/cmd/contour/servecontext_test.go b/cmd/contour/servecontext_test.go
index cf205d81..9015b9d0 100644
--- a/cmd/contour/servecontext_test.go
+++ b/cmd/contour/servecontext_test.go
@@ -24,8 +24,10 @@ import (
        "testing"
        "time"

+       "github.com/pkg/errors"
        "github.com/stretchr/testify/assert"
        "github.com/tsaarni/certyaml"
+       "golang.org/x/net/http2"
        "google.golang.org/grpc"
        "k8s.io/utils/ptr"

@@ -286,21 +288,25 @@ func checkFatalErr(t *testing.T, err error) {
 // tryConnect tries to establish TLS connection to the server.
 // If successful, return the server certificate.
 func tryConnect(address string, clientCert tls.Certificate, caCertPool *x509.CertPool) (*x509.Certificate, error) {
+       rawConn, err := net.Dial("tcp", address)
+       if err != nil {
+               rawConn.Close()
+               return nil, errors.Wrapf(err, "error dialing %s", address)
+       }
+
        clientConfig := &tls.Config{
                ServerName:   "localhost",
                MinVersion:   tls.VersionTLS13,
                Certificates: []tls.Certificate{clientCert},
                RootCAs:      caCertPool,
+               NextProtos:   []string{http2.NextProtoTLS},
        }
-       conn, err := tls.Dial("tcp", address, clientConfig)
-       if err != nil {
-               return nil, err
-       }
+
+       conn := tls.Client(rawConn, clientConfig)
        defer conn.Close()

-       err = peekError(conn)
-       if err != nil {
-               return nil, err
+       if err := peekError(conn); err != nil {
+               return nil, errors.Wrap(err, "error peeking TLS alert")
        }

        return conn.ConnectionState().PeerCertificates[0], nil
