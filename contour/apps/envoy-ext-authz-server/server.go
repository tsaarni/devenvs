package main

import "fmt"

type authzServer struct {
	name          string
	defaultStatus int
	quiet         bool          // suppress all per-request logging
	metrics       *checkMetrics // nil disables metrics
}

// needHeaders reports whether request headers are used by the log or the metrics.
func (s *authzServer) needHeaders() bool {
	return !s.quiet || s.metrics != nil
}

func (s *authzServer) decisionFor(get func(string) string) decision {
	d, err := parseDecision(get, s.defaultStatus)
	if err != nil && !s.quiet {
		log.Warn("Invalid authorization controls; denying with 400", "name", s.name, "error", err)
	}
	return d
}

func (s *authzServer) deniedBody(status int) string {
	return fmt.Sprintf("denied by %s (%d)\n", s.name, status)
}
