package main

import (
	"io"
	"net/http"
	"time"
)

func (s *authzServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := s.decisionFor(r.Header.Get)
	start := time.Now()
	var flat map[string]string // only needed for logging and metrics
	if s.needHeaders() {
		flat = make(map[string]string, len(r.Header))
		for k := range r.Header {
			flat[k] = r.Header.Get(k)
		}
	}
	info := checkRequestInfo{
		protocol: "http", method: r.Method, host: r.Host, path: r.URL.RequestURI(), tls: r.TLS,
		requestID: r.Header.Get("x-request-id"), headers: flat,
	}
	outcome := d.outcome()
	var checkErr error
	defer func() { s.record(info, d, outcome, checkErr, start) }()
	info.bodyBytes, checkErr = io.Copy(io.Discard, r.Body)
	if checkErr != nil {
		outcome = "read-error"
		panic(http.ErrAbortHandler)
	}
	if checkErr = d.wait(r.Context()); checkErr != nil {
		outcome = "canceled"
		return
	}
	if d.fail != "" {
		// HTTP error responses (including 5xx) are denials, not transport failures.
		panic(http.ErrAbortHandler)
	}
	w.Header().Set("x-authz-server", s.name)
	if d.tag != "" {
		w.Header().Set("x-authz-tag", d.tag)
	}
	if d.status == http.StatusOK {
		w.Header().Set("x-authz-verdict", "allow")
	}
	w.WriteHeader(d.status)
	if d.status != http.StatusOK {
		_, checkErr = io.WriteString(w, s.deniedBody(d.status))
		if checkErr != nil && !s.quiet {
			log.Warn("Writing denial failed", "name", s.name, "error", checkErr)
		}
	}
}
