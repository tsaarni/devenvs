package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// httpAuthzHandler implements the HTTP ext-authz service.
func httpAuthzHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("HTTP Check request", "method", r.Method, "host", r.Host, "uri", r.URL.RequestURI(), "headers", formatHTTPHeaders(r.Header))

	// Log TLS info and client certificate if mTLS is used.
	if r.TLS != nil {
		log.Info("HTTP TLS info", "sni", r.TLS.ServerName)
		for _, cert := range r.TLS.PeerCertificates {
			log.Info("HTTP client certificate", "subject", cert.Subject)
		}
	}

	// TODO: Add your authorization logic here.

	// Response headers to add/modify on the upstream request:
	// w.Header().Set("x-user", "alice")

	// Set a cookie on the downstream response.
	// NOTE: Envoy only forwards these if allowed_client_headers_on_success is configured
	// in the ext-authz filter.
	//w.Header().Set("Set-Cookie", "session=http-authz-test; Path=/; HttpOnly")

	status := http.StatusOK
	log.Info("HTTP Check result", "status", status)
	w.WriteHeader(status)
}

func formatHTTPHeaders(headers http.Header) string {
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, strings.Join(headers[k], ",")))
	}
	return strings.Join(parts, ", ")
}
