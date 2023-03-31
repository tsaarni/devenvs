package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}

	trustedCert, err := os.ReadFile("ca.pem")
	if err != nil {
		panic(err)
	}

	rootCAs.AppendCertsFromPEM(trustedCert)

	clientCert, err := tls.LoadX509KeyPair("client.pem", "client-key.pem")
	if err != nil {
		panic(err)
	}

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			//CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			RootCAs:            rootCAs,
			ClientSessionCache: tls.NewLRUClientSessionCache(1000),
		}}}

	req, err := http.NewRequest(http.MethodGet, "https://localhost:8443", nil)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("Sending request...")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
			//fmt.Println("error:", err)
	                //fmt.Println("Sleeping 5 sec...")
      		        //time.Sleep(time.Second * 5)
                        //continue
		}

		fmt.Println("Response status:", resp.Status)
		resp.Body.Close()

		fmt.Println("Sleeping 5 sec...")
		time.Sleep(time.Second * 5)
	}
}
