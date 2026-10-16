package main

import "fmt"

type authzServer struct {
	name          string
	defaultStatus int
}

func (s *authzServer) decisionFor(get func(string) string) decision {
	d, err := parseDecision(get, s.defaultStatus)
	if err != nil {
		log.Warn("Invalid authorization controls; denying with 400", "name", s.name, "error", err)
	}
	return d
}

func (s *authzServer) deniedBody(status int) string {
	return fmt.Sprintf("denied by %s (%d)\n", s.name, status)
}
