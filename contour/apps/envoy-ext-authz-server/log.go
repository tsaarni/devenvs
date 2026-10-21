package main

import (
	"crypto/tls"
	"fmt"
	"sort"
	"strings"
	"time"
)

// checkRequestInfo describes the request Envoy sent in an authorization check.
type checkRequestInfo struct {
	protocol, method, host, path string
	contextExtensions            map[string]string // gRPC only; Envoy does not send these to HTTP services
	bodyBytes                    int64
	requestID                    string // Envoy's x-request-id
	headers                      map[string]string
	tls                          *tls.ConnectionState
}

func formatHeaders(h map[string]string) string {
	parts := make([]string, 0, len(h))
	for k, v := range h {
		parts = append(parts, fmt.Sprintf("%s=%q", strings.ToLower(k), v))
	}
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

func formatContext(c map[string]string) string {
	parts := make([]string, 0, len(c))
	for k, v := range c {
		parts = append(parts, k+":"+v)
	}
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

func tlsIdentity(info checkRequestInfo) (sni, clientCert string) {
	if info.tls != nil {
		sni = info.tls.ServerName
		if len(info.tls.PeerCertificates) > 0 {
			clientCert = info.tls.PeerCertificates[0].Subject.String()
		}
	}
	return sni, clientCert
}

// record publishes the result of a check to the metrics and the log.
func (s *authzServer) record(info checkRequestInfo, d decision, outcome string, err error, start time.Time) {
	s.metrics.observe(info, d, outcome, time.Since(start))
	s.logDecision(info, d, outcome, err)
}

func (s *authzServer) logDecision(info checkRequestInfo, d decision, outcome string, err error) {
	if s.quiet {
		return
	}
	sni, clientCert := tlsIdentity(info)
	status := d.status
	if outcome != "allow" && outcome != "deny" {
		status = 0
	}
	attrs := []any{
		"name", s.name, "tag", d.tag, "request_id", info.requestID, "proto", info.protocol, "tls", info.tls != nil,
		"sni", sni, "client_cert", clientCert, "method", info.method,
		"host", info.host, "path", info.path,
		"body_bytes", info.bodyBytes, "controls", d.controls, "verdict", outcome, "status", status,
	}
	if len(info.contextExtensions) > 0 {
		attrs = append(attrs, "ctx", formatContext(info.contextExtensions))
	}
	attrs = append(attrs, "headers", formatHeaders(info.headers))
	if err != nil {
		attrs = append(attrs, "error", err)
	}
	log.Info("CHECK", attrs...)
}
