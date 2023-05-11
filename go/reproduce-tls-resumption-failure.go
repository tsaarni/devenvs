package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"time"
)

func server() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Server: received request from", req.RemoteAddr)
		w.Write([]byte("OK.\n"))
	})

	cert, err := tls.X509KeyPair(serverCertPem, serverKeyPem)
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{
		// Set curve preference without X25519 to trigger mismatch with the client.
		// https://www.rfc-editor.org/rfc/rfc8446#page-14
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},

		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
	}

	srv := &http.Server{
		Addr:         ":8443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
}

func client() {
	rootCA := x509.NewCertPool()
	rootCA.AppendCertsFromPEM(caCertPem)

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			// Enable session resumption.
			ClientSessionCache: tls.NewLRUClientSessionCache(1000),
			// Client's CurvePreferences are not set.
			RootCAs: rootCA,
		}}}

	req, err := http.NewRequest(http.MethodGet, "https://localhost:8443", nil)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Client: sending request to", req.URL.String())
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		fmt.Println("Client: received response from server with status:", resp.Status)
		resp.Body.Close()

		fmt.Println("Client: sleeping 5 sec...")
		time.Sleep(time.Second * 5)
	}
}

func main() {

	go server()
	go client()

	// Wait forever.
	select {}
}

var caCertPem = []byte(`-----BEGIN CERTIFICATE-----
MIIBVDCB+qADAgECAggXUqw36eYA9jAKBggqhkjOPQQDAjANMQswCQYDVQQDEwJj
YTAgFw05MDAxMDEwOTAwMDBaGA8yMTAwMDEwMTA5MDAwMFowDTELMAkGA1UEAxMC
Y2EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASk+bloFMZZQy6lcCR3wVufQ8g6
jeWAYDo1vtLZyrY1xbLjAtxsQBT/UMK2xBFI/3fNYqMjLP89N4Ri7VhH5A2Io0Iw
QDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQULarJ
BBKTTjf9xtAtASkej7dkoHcwCgYIKoZIzj0EAwIDSQAwRgIhAIAulUBbvvHWYA+C
SdD1TMb0uGSWhNj7iLl2h/P4WcghAiEA13jQhgUtT7ht4MgBSHvdurR0MpZPwfmy
SNJlueDSJwI=
-----END CERTIFICATE-----
`)

var serverCertPem = []byte(`-----BEGIN CERTIFICATE-----
MIIBYDCCAQWgAwIBAgIIF1KsN+nykdEwCgYIKoZIzj0EAwIwDTELMAkGA1UEAxMC
Y2EwIBcNOTAwMTAxMDkwMDAwWhgPMjEwMDAxMDEwOTAwMDBaMBExDzANBgNVBAMT
BnNlcnZlcjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJl0Knv6tzXUcRaB3YmB
mda81SsDRSk0ZYaZaaRcDKMDFUbCM5VL5CgbILoH5OazQzlRkqJdObtW1u8Puzct
06ijSTBHMA4GA1UdDwEB/wQEAwIFoDAfBgNVHSMEGDAWgBQtqskEEpNON/3G0C0B
KR6Pt2SgdzAUBgNVHREEDTALgglsb2NhbGhvc3QwCgYIKoZIzj0EAwIDSQAwRgIh
ANKnrkG3PJHdVF0SlU/07yJeAbZ7XpgGTHqL9AEt6lDOAiEA8GLpDUhEVB11il5R
JYT+qMkOmdLy3NoVOMkjhYVf7AM=
-----END CERTIFICATE-----
`)

var serverKeyPem = []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgiCP9Pzp1esKctEV4
P6hLUbJZEAO1NgTTuQ0iah/G3VGhRANCAASZdCp7+rc11HEWgd2JgZnWvNUrA0Up
NGWGmWmkXAyjAxVGwjOVS+QoGyC6B+Tms0M5UZKiXTm7VtbvD7s3LdOo
-----END PRIVATE KEY-----
`)
