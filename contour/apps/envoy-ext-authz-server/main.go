package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/lmittmann/tint"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log = slog.Default()

func buildTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("loading TLS cert/key: %w", err)
	}
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("reading CA file: %w", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("CA file %q contains no certificates", caFile)
		}
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		log.Info("mTLS enabled", "ca", caFile)
	}

	return tlsConfig, nil
}

// mustTLSConfig builds a TLS config from the given cert/key/ca files.
// Returns nil if cert and key are empty (TLS disabled). Exits on error.
func mustTLSConfig(cert, key, ca string) *tls.Config {
	if cert == "" && key == "" && ca == "" {
		return nil
	}
	if cert == "" || key == "" {
		log.Error("TLS requires both certificate and key")
		os.Exit(1)
	}
	cfg, err := buildTLSConfig(cert, key, ca)
	if err != nil {
		log.Error("TLS configuration failed", "error", err)
		os.Exit(1)
	}
	return cfg
}

// listenAndServe starts srv with TLS if tlsCfg is non-nil, otherwise plaintext.
func listenAndServe(srv *http.Server, tlsCfg *tls.Config) error {
	if tlsCfg == nil {
		return srv.ListenAndServe()
	}
	return srv.ListenAndServeTLS("", "")
}

// config holds all settings. Each flag has a matching AUTHZ_* environment variable
// (e.g. --grpc-tls-cert / AUTHZ_GRPC_TLS_CERT); command line takes precedence.
type config struct {
	Name          string `help:"Server identity in logs and responses." default:"authz"`
	DefaultStatus int    `help:"Verdict status when no x-authz-status is sent (200 allows, 201..599 denies)." default:"200"`

	GRPCAddr    string `name:"grpc-addr" help:"gRPC listen address (empty disables)." default:":9001"`
	GRPCTLSCert string `name:"grpc-tls-cert" help:"gRPC server TLS certificate file." type:"path"`
	GRPCTLSKey  string `name:"grpc-tls-key" help:"gRPC server TLS key file." type:"path"`
	GRPCTLSCA   string `name:"grpc-tls-ca" help:"gRPC client CA file (enables mTLS)." type:"path"`

	HTTPAddr    string `name:"http-addr" help:"HTTP listen address (empty disables)." default:":9002"`
	HTTPTLSCert string `name:"http-tls-cert" help:"HTTP server TLS certificate file." type:"path"`
	HTTPTLSKey  string `name:"http-tls-key" help:"HTTP server TLS key file." type:"path"`
	HTTPTLSCA   string `name:"http-tls-ca" help:"HTTP client CA file (enables mTLS)." type:"path"`
}

func main() {
	log = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	slog.SetDefault(log)

	var cli config
	kong.Parse(&cli,
		kong.Name("envoy-ext-authz-server"),
		kong.Description("Envoy ext_authz test server. Every flag can also be set with the environment variable shown in its help."),
		kong.DefaultEnvars("AUTHZ"),
	)
	name, defaultStatus := &cli.Name, &cli.DefaultStatus
	grpcAddr, httpAddr := &cli.GRPCAddr, &cli.HTTPAddr
	grpcTLSCert, grpcTLSKey, grpcTLSCA := &cli.GRPCTLSCert, &cli.GRPCTLSKey, &cli.GRPCTLSCA
	httpTLSCert, httpTLSKey, httpTLSCA := &cli.HTTPTLSCert, &cli.HTTPTLSKey, &cli.HTTPTLSCA

	if !validStatus(*defaultStatus) || *name == "" || strings.IndexFunc(*name, func(r rune) bool {
		return r < 32 || r > 126
	}) >= 0 || (*grpcAddr == "" && *httpAddr == "") {
		log.Error("Invalid configuration: require printable nonempty name, status 200..599, and at least one listener")
		os.Exit(1)
	}
	server := &authzServer{name: *name, defaultStatus: *defaultStatus}
	log = log.With("instance", *name)

	grpcTLS := mustTLSConfig(*grpcTLSCert, *grpcTLSKey, *grpcTLSCA)
	httpTLS := mustTLSConfig(*httpTLSCert, *httpTLSKey, *httpTLSCA)

	errors := make(chan error, 2)
	if *grpcAddr != "" {
		go func() {
			lis, err := net.Listen("tcp", *grpcAddr)
			if err != nil {
				errors <- fmt.Errorf("gRPC listen: %w", err)
				return
			}
			var opts []grpc.ServerOption
			if grpcTLS != nil {
				opts = append(opts, grpc.Creds(credentials.NewTLS(grpcTLS)))
			}
			srv := grpc.NewServer(opts...)
			authv3.RegisterAuthorizationServer(srv, &grpcAuthzServer{authzServer: server})
			log.Info("gRPC ext-authz server listening", "addr", lis.Addr(), "tls", grpcTLS != nil)
			errors <- srv.Serve(lis)
		}()
	}

	// Start HTTP ext-authz server.
	// h2c.NewHandler enables h2c (HTTP/2 cleartext) in addition to HTTP/1.1 on plaintext connections.
	// http2.ConfigureServer enables HTTP/2 via ALPN on TLS connections.
	h2s := &http2.Server{}
	httpServer := &http.Server{
		Addr:    *httpAddr,
		Handler: h2c.NewHandler(server, h2s),
	}
	if httpTLS != nil {
		httpServer.TLSConfig = httpTLS
		if err := http2.ConfigureServer(httpServer, h2s); err != nil {
			log.Error("HTTP/2 configuration failed", "error", err)
			os.Exit(1)
		}
	}
	if *httpAddr != "" {
		go func() {
			log.Info("HTTP ext-authz server starting", "addr", *httpAddr, "tls", httpTLS != nil)
			errors <- listenAndServe(httpServer, httpTLS)
		}()
	}
	log.Error("Listener stopped", "error", <-errors)
	os.Exit(1)
}
