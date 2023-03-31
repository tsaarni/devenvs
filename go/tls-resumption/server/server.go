package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Request from:", req.RemoteAddr)
		w.Write([]byte("This is an example server.\n"))
	})

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	certs, err := os.ReadFile("ca.pem")
	if err != nil {
		panic(err)
	}

	rootCAs.AppendCertsFromPEM(certs)

	cfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// Default curves:  X25519, prime256v1, secp384r1, secp521r1
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		ClientAuth:       tls.ClientAuthType(tls.RequireAndVerifyClientCert),
		ClientCAs:        rootCAs,
	}

	srv := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS("server.pem", "server-key.pem"))
}
