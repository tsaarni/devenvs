package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/lmittmann/tint"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var log *slog.Logger

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
		pool.AppendCertsFromPEM(caCert)
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		log.Info("mTLS enabled", "ca", caFile)
	}

	return tlsConfig, nil
}

// mustTLSConfig builds a TLS config from the given cert/key/ca files.
// Returns nil if cert and key are empty (TLS disabled). Exits on error.
func mustTLSConfig(cert, key, ca string) *tls.Config {
	if cert == "" || key == "" {
		return nil
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
	lis, err := tls.Listen("tcp", srv.Addr, tlsCfg)
	if err != nil {
		return err
	}
	return srv.Serve(lis)
}

func main() {
	log = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
	slog.SetDefault(log)

	// gRPC server: implements Envoy ext-authz gRPC "Check" protocol (envoy.service.auth.v3.Authorization).
	// Always uses HTTP/2: h2c (plaintext) when no TLS is configured, h2 (TLS) when TLS is configured.
	grpcPort := flag.Int("grpc-port", 9001, "gRPC ext-authz server port")
	grpcTLSCert := flag.String("grpc-tls-cert", "", "gRPC server TLS certificate file")
	grpcTLSKey := flag.String("grpc-tls-key", "", "gRPC server TLS key file")
	grpcTLSCA := flag.String("grpc-tls-ca", "", "gRPC server TLS CA file for client cert verification (enables mTLS)")

	// HTTP server: implements Envoy ext-authz HTTP protocol (plain HTTP check request/response).
	// Supports HTTP/1.1 and HTTP/2: plaintext (HTTP/h2c) when no TLS, encrypted (HTTPS/h2) when TLS is configured.
	httpPort := flag.Int("http-port", 9002, "HTTP ext-authz server port")
	httpTLSCert := flag.String("http-tls-cert", "", "HTTP server TLS certificate file")
	httpTLSKey := flag.String("http-tls-key", "", "HTTP server TLS key file")
	httpTLSCA := flag.String("http-tls-ca", "", "HTTP server TLS CA file for client cert verification (enables mTLS)")

	flag.Parse()

	log.Info("Starting envoy ext-authz server", "grpc-port", *grpcPort, "http-port", *httpPort)

	grpcTLS := mustTLSConfig(*grpcTLSCert, *grpcTLSKey, *grpcTLSCA)
	httpTLS := mustTLSConfig(*httpTLSCert, *httpTLSKey, *httpTLSCA)

	// Start gRPC ext-authz server.
	// gRPC requires HTTP/2; without TLS this is h2c, with TLS this is h2 negotiated via ALPN.
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			log.Error("gRPC listen failed", "error", err)
			os.Exit(1)
		}
		var opts []grpc.ServerOption
		if grpcTLS != nil {
			opts = append(opts, grpc.Creds(credentials.NewTLS(grpcTLS)))
		}
		srv := grpc.NewServer(opts...)
		authv3.RegisterAuthorizationServer(srv, &grpcAuthzServer{})
		log.Info("gRPC ext-authz server listening", "port", *grpcPort, "tls", grpcTLS != nil)
		if err := srv.Serve(lis); err != nil {
			log.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Start HTTP ext-authz server.
	// h2c.NewHandler enables h2c (HTTP/2 cleartext) in addition to HTTP/1.1 on plaintext connections.
	// http2.ConfigureServer enables HTTP/2 via ALPN on TLS connections.
	h2s := &http2.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpAuthzHandler)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: h2c.NewHandler(mux, h2s),
	}
	if httpTLS != nil {
		httpServer.TLSConfig = httpTLS
		http2.ConfigureServer(httpServer, h2s) // adds "h2" to TLS NextProtos for ALPN
	}
	log.Info("HTTP ext-authz server listening", "port", *httpPort, "tls", httpTLS != nil)
	if err := listenAndServe(httpServer, httpTLS); err != nil {
		log.Error("HTTP server failed", "error", err)
		os.Exit(1)
	}
}
