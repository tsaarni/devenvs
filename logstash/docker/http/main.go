package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request", "proto", r.Proto, "method", r.Method, "url", r.URL, "remote", r.RemoteAddr)

		// Print client certificate information.
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			slog.Info("Client certificate", "subject", r.TLS.PeerCertificates[0].Subject, "issuer", r.TLS.PeerCertificates[0].Issuer)
		}

		// write request body to stdout
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading request body", "error", err)
		}
		slog.Info("Content", "content-type", r.Header["Content-Type"], "body", string(reqBody))

		// write response
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		//		w.Write([]byte("Hello, TLS!"))
	})

	// Read server certificate and private key.
	cert, err := tls.LoadX509KeyPair("/certs/server.pem", "/certs/server-key.pem")
	if err != nil {
		slog.Error("Error loading server certificate", "error", err)
		return
	}

	// Read CA certificate and create cert pool for client authentication.
	caCert, err := os.ReadFile("/certs/client-ca.pem")
	if err != nil {
		slog.Error("Error loading client CA certificate", "error", err)
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a key log file for wireshark.
	//f, _ := os.OpenFile("/tmp/wireshark-keys.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		// KeyLogWriter: f,
		// Force TLSv1.2 for better visibility of the TLS handshake.
		MaxVersion: tls.VersionTLS12,
	}

	address := "0.0.0.0:8443"

	server := &http.Server{
		Addr:      address,
		TLSConfig: config,
	}

	slog.Info("Server started", "address", address)
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}
