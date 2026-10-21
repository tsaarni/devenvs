package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type decision struct {
	status   int
	delay    time.Duration
	fail     string
	tag      string // tester-supplied tag to tell requests apart, echoed in logs, metrics and responses
	controls string
	used     []controlUse // controls found in the request and how they were sent
}

// controlUse records one control sent by the client; source is "header" or "authorization".
type controlUse struct{ name, source string }

var controlNames = map[string]bool{"status": true, "delay": true, "fail": true, "tag": true}

func validTag(tag string) bool {
	if len(tag) > 128 {
		return false
	}
	return strings.IndexFunc(tag, func(r rune) bool { return r <= ' ' || r > '~' }) < 0
}

func validStatus(status int) bool {
	return status >= 200 && status <= 599
}

const authzScheme = "x-authz"

// parseAuthorizationHeader parses `Authorization: x-authz status=403 delay=2s fail=error`.
// Other schemes yield no params so ordinary credentials are ignored.
func parseAuthorizationHeader(value string) (map[string]string, error) {
	scheme, rest, _ := strings.Cut(strings.TrimSpace(value), " ")
	if !strings.EqualFold(scheme, authzScheme) {
		return nil, nil
	}
	params := map[string]string{}
	for _, field := range strings.FieldsFunc(rest, func(r rune) bool { return r == ' ' || r == ',' }) {
		k, v, ok := strings.Cut(field, "=")
		if !ok || !controlNames[k] {
			return nil, fmt.Errorf("invalid Authorization parameter %q: expected status=, delay=, fail= or tag=", field)
		}
		params[k] = v
	}
	return params, nil
}

// parseDecision reads the controls from x-authz-* headers, falling back per control
// to the Authorization header; the dedicated headers take precedence.
func parseDecision(get func(string) string, defaultStatus int) (decision, error) {
	params, authErr := parseAuthorizationHeader(get("authorization"))
	var via []string
	var used []controlUse
	pick := func(control string) string {
		if v := get("x-authz-" + control); v != "" {
			used = append(used, controlUse{control, "header"})
			return v
		}
		if v := params[control]; v != "" {
			via = append(via, control)
			used = append(used, controlUse{control, "authorization"})
			return v
		}
		return ""
	}
	status, delay, fail, tag := pick("status"), pick("delay"), pick("fail"), pick("tag")
	d := decision{status: defaultStatus, fail: fail, tag: tag, controls: "none (default status)", used: used}
	var set []string
	for _, c := range []struct{ name, value string }{{"status", status}, {"delay", delay}, {"fail", fail}, {"tag", tag}} {
		if c.value != "" {
			set = append(set, fmt.Sprintf("%s=%q", c.name, c.value))
		}
	}
	if len(set) > 0 {
		d.controls = strings.Join(set, " ")
		if len(via) > 0 {
			d.controls += fmt.Sprintf(" authorization=%v", via)
		}
	}
	invalid := func(err error) (decision, error) {
		return decision{status: http.StatusBadRequest, tag: d.tag, controls: d.controls, used: d.used}, err
	}
	if authErr != nil {
		return invalid(authErr)
	}
	if !validTag(tag) {
		d.tag = ""
		return invalid(fmt.Errorf("invalid x-authz-tag %q: expected up to 128 printable characters without spaces", tag))
	}
	if status != "" {
		n, err := strconv.Atoi(status)
		if err != nil || !validStatus(n) {
			return invalid(fmt.Errorf("invalid x-authz-status %q: expected 200..599", status))
		}
		d.status = n
	}
	if delay != "" {
		var err error
		d.delay, err = time.ParseDuration(delay)
		if err != nil || d.delay < 0 {
			return invalid(fmt.Errorf("invalid x-authz-delay %q: expected a nonnegative Go duration", delay))
		}
	}
	if fail != "" && fail != "error" && fail != "reset" {
		return invalid(fmt.Errorf("invalid x-authz-fail %q: expected error or reset", fail))
	}
	return d, nil
}

func (d decision) wait(ctx context.Context) error {
	if d.delay == 0 {
		return ctx.Err()
	}
	timer := time.NewTimer(d.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (d decision) outcome() string {
	if d.fail != "" {
		return d.fail
	}
	if d.status == http.StatusOK {
		return "allow"
	}
	return "deny"
}
