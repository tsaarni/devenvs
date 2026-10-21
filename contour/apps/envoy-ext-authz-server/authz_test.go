package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func headers(kv ...string) func(string) string {
	m := map[string]string{}
	for i := 0; i < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	return func(k string) string { return m[k] }
}

func TestParseDecision(t *testing.T) {
	tests := []struct {
		name    string
		kv      []string
		status  int
		delay   time.Duration
		fail    string
		tag     string
		verdict string
		wantErr bool
	}{
		{name: "default allow", status: 200, verdict: "allow"},
		{name: "deny", kv: []string{"x-authz-status", "403"}, status: 403, verdict: "deny"},
		{name: "delay and deny", kv: []string{"x-authz-delay", "2s", "x-authz-status", "401"}, status: 401, delay: 2 * time.Second, verdict: "deny"},
		{name: "error", kv: []string{"x-authz-fail", "error"}, status: 200, fail: "error", verdict: "error"},
		{name: "reset", kv: []string{"x-authz-fail", "reset"}, status: 200, fail: "reset", verdict: "reset"},
		{name: "bad status", kv: []string{"x-authz-status", "abc"}, status: 400, verdict: "deny", wantErr: true},
		{name: "out of range status", kv: []string{"x-authz-status", "99"}, status: 400, verdict: "deny", wantErr: true},
		{name: "bad delay", kv: []string{"x-authz-delay", "soon"}, status: 400, verdict: "deny", wantErr: true},
		{name: "negative delay", kv: []string{"x-authz-delay", "-1s"}, status: 400, verdict: "deny", wantErr: true},
		{name: "bad fail", kv: []string{"x-authz-fail", "boom"}, status: 400, verdict: "deny", wantErr: true},
		{name: "authorization", kv: []string{"authorization", "x-authz status=403 delay=1s"}, status: 403, delay: time.Second, verdict: "deny"},
		{name: "authorization scheme is case-insensitive", kv: []string{"authorization", "X-Authz fail=error,status=401"}, status: 401, fail: "error", verdict: "error"},
		{name: "headers win over authorization", kv: []string{"authorization", "x-authz status=403", "x-authz-status", "401"}, status: 401, verdict: "deny"},
		{name: "authorization fills missing controls", kv: []string{"authorization", "x-authz delay=1s", "x-authz-status", "401"}, status: 401, delay: time.Second, verdict: "deny"},
		{name: "other schemes ignored", kv: []string{"authorization", "Bearer status=403"}, status: 200, verdict: "allow"},
		{name: "authorization unknown key", kv: []string{"authorization", "x-authz stats=403"}, status: 400, verdict: "deny", wantErr: true},
		{name: "tag header", kv: []string{"x-authz-tag", "t1"}, status: 200, verdict: "allow", tag: "t1"},
		{name: "tag via authorization", kv: []string{"authorization", "x-authz tag=t2 status=403"}, status: 403, verdict: "deny", tag: "t2"},
		{name: "tag header wins", kv: []string{"authorization", "x-authz tag=a", "x-authz-tag", "b"}, status: 200, verdict: "allow", tag: "b"},
		{name: "invalid tag", kv: []string{"x-authz-tag", "has space"}, status: 400, verdict: "deny", wantErr: true},
		{name: "authorization bad value", kv: []string{"authorization", "x-authz status=abc"}, status: 400, verdict: "deny", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := parseDecision(headers(tt.kv...), 200)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if d.status != tt.status || d.delay != tt.delay || d.fail != tt.fail || d.tag != tt.tag || d.outcome() != tt.verdict {
				t.Fatalf("got %+v verdict %s", d, d.outcome())
			}
		})
	}
}

func TestParseDecisionDefaultStatus(t *testing.T) {
	d, err := parseDecision(headers(), 403)
	if err != nil || d.status != 403 || d.controls != "none (default status)" {
		t.Fatalf("got %+v, %v", d, err)
	}
}

func newTestServer() *authzServer {
	return &authzServer{name: "test", defaultStatus: 200}
}

func doHTTP(t *testing.T, client *http.Client, url string, kv ...string) (*http.Response, string, error) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	for i := 0; i < len(kv); i += 2 {
		req.Header.Set(kv[i], kv[i+1])
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp, string(body), nil
}

func TestHTTP(t *testing.T) {
	ts := httptest.NewTLSServer(newTestServer())
	defer ts.Close()

	t.Run("allow", func(t *testing.T) {
		resp, _, err := doHTTP(t, ts.Client(), ts.URL+"/check")
		if err != nil || resp.StatusCode != 200 {
			t.Fatalf("resp %v err %v", resp, err)
		}
		if resp.Header.Get("x-authz-server") != "test" || resp.Header.Get("x-authz-verdict") != "allow" {
			t.Fatalf("headers %v", resp.Header)
		}
	})
	t.Run("deny", func(t *testing.T) {
		resp, body, err := doHTTP(t, ts.Client(), ts.URL, "x-authz-status", "403")
		if err != nil || resp.StatusCode != 403 || body != "denied by test (403)\n" {
			t.Fatalf("resp %v body %q err %v", resp, body, err)
		}
		if resp.Header.Get("x-authz-server") != "test" {
			t.Fatalf("headers %v", resp.Header)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		resp, _, err := doHTTP(t, ts.Client(), ts.URL, "x-authz-status", "nope")
		if err != nil || resp.StatusCode != 400 {
			t.Fatalf("resp %v err %v", resp, err)
		}
	})
	t.Run("delay", func(t *testing.T) {
		start := time.Now()
		resp, _, err := doHTTP(t, ts.Client(), ts.URL, "x-authz-delay", "200ms")
		if err != nil || resp.StatusCode != 200 || time.Since(start) < 200*time.Millisecond {
			t.Fatalf("resp %v err %v elapsed %v", resp, err, time.Since(start))
		}
	})
	t.Run("fail aborts the connection", func(t *testing.T) {
		for _, mode := range []string{"error", "reset"} {
			if _, _, err := doHTTP(t, ts.Client(), ts.URL, "x-authz-fail", mode); err == nil {
				t.Fatalf("%s: expected connection error", mode)
			}
		}
	})
}

func startGRPC(t *testing.T, tlsCfg *tls.Config) *grpc.ClientConn {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	var serverOpts []grpc.ServerOption
	var dialCreds grpc.DialOption = grpc.WithTransportCredentials(insecure.NewCredentials())
	if tlsCfg != nil {
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(tlsCfg)))
		pool := x509.NewCertPool()
		leaf, _ := x509.ParseCertificate(tlsCfg.Certificates[0].Certificate[0])
		pool.AddCert(leaf)
		dialCreds = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{RootCAs: pool, ServerName: "example.com"}))
	}
	srv := grpc.NewServer(serverOpts...)
	authv3.RegisterAuthorizationServer(srv, &grpcAuthzServer{authzServer: newTestServer()})
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)
	conn, err := grpc.NewClient(lis.Addr().String(), dialCreds)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func checkRequest(kv ...string) *authv3.CheckRequest {
	h := map[string]string{}
	for i := 0; i < len(kv); i += 2 {
		h[kv[i]] = kv[i+1]
	}
	return &authv3.CheckRequest{Attributes: &authv3.AttributeContext{
		ContextExtensions: map[string]string{"env": "staging"},
		Request: &authv3.AttributeContext_Request{Http: &authv3.AttributeContext_HttpRequest{
			Method: "GET", Host: "example.com", Path: "/check", Headers: h,
		}},
	}}
}

func headerMap(hs []*corev3.HeaderValueOption) map[string]string {
	m := map[string]string{}
	for _, h := range hs {
		m[h.GetHeader().GetKey()] = h.GetHeader().GetValue()
	}
	return m
}

func testGRPC(t *testing.T, conn *grpc.ClientConn) {
	client := authv3.NewAuthorizationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("allow", func(t *testing.T) {
		resp, err := client.Check(ctx, checkRequest())
		if err != nil || resp.GetStatus().GetCode() != int32(codes.OK) {
			t.Fatalf("resp %v err %v", resp, err)
		}
		h := headerMap(resp.GetOkResponse().GetHeaders())
		if h["x-authz-server"] != "test" || h["x-authz-verdict"] != "allow" {
			t.Fatalf("headers %v", h)
		}
	})
	t.Run("deny", func(t *testing.T) {
		resp, err := client.Check(ctx, checkRequest("x-authz-status", "401"))
		if err != nil || resp.GetStatus().GetCode() != int32(codes.PermissionDenied) {
			t.Fatalf("resp %v err %v", resp, err)
		}
		denied := resp.GetDeniedResponse()
		if denied.GetStatus().GetCode() != 401 || denied.GetBody() != "denied by test (401)\n" ||
			headerMap(denied.GetHeaders())["x-authz-server"] != "test" {
			t.Fatalf("denied %v", denied)
		}
	})
	t.Run("invalid denies with 400", func(t *testing.T) {
		resp, err := client.Check(ctx, checkRequest("x-authz-delay", "later"))
		if err != nil || resp.GetDeniedResponse().GetStatus().GetCode() != 400 {
			t.Fatalf("resp %v err %v", resp, err)
		}
	})
	t.Run("error", func(t *testing.T) {
		_, err := client.Check(ctx, checkRequest("x-authz-fail", "error"))
		if status.Code(err) != codes.Unavailable {
			t.Fatalf("err %v", err)
		}
	})
	t.Run("reset is unsupported", func(t *testing.T) {
		_, err := client.Check(ctx, checkRequest("x-authz-fail", "reset"))
		if status.Code(err) != codes.Unimplemented {
			t.Fatalf("err %v", err)
		}
	})
	t.Run("delay honors cancellation", func(t *testing.T) {
		short, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
		defer cancel()
		_, err := client.Check(short, checkRequest("x-authz-delay", "10s"))
		if status.Code(err) != codes.DeadlineExceeded {
			t.Fatalf("err %v", err)
		}
	})
}

func TestGRPCPlaintext(t *testing.T) {
	testGRPC(t, startGRPC(t, nil))
}

func TestGRPCTLS(t *testing.T) {
	// httptest ships a certificate valid for example.com.
	ts := httptest.NewTLSServer(http.NotFoundHandler())
	defer ts.Close()
	testGRPC(t, startGRPC(t, &tls.Config{Certificates: ts.TLS.Certificates}))
}
