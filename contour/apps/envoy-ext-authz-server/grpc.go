package main

import (
	"context"
	"net/http"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type grpcAuthzServer struct {
	authv3.UnimplementedAuthorizationServer
	*authzServer
}

func responseHeader(key, value string) *corev3.HeaderValueOption {
	return &corev3.HeaderValueOption{
		Header:       &corev3.HeaderValue{Key: key, Value: value},
		AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
	}
}

func (s *grpcAuthzServer) Check(ctx context.Context, req *authv3.CheckRequest) (resp *authv3.CheckResponse, err error) {
	start := time.Now()
	request := req.GetAttributes().GetRequest().GetHttp()
	d := s.decisionFor(func(key string) string { return request.GetHeaders()[key] })
	info := checkRequestInfo{
		protocol: "grpc", method: request.GetMethod(), host: request.GetHost(),
		path: request.GetPath(), contextExtensions: req.GetAttributes().GetContextExtensions(),
		bodyBytes: int64(len(request.GetBody()) + len(request.GetRawBody())),
		requestID: request.GetId(), headers: request.GetHeaders(),
	}
	if info.requestID == "" {
		info.requestID = request.GetHeaders()["x-request-id"]
	}
	if p, ok := peer.FromContext(ctx); ok {
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			info.tls = &tlsInfo.State
		}
	}
	outcome := d.outcome()
	defer func() { s.record(info, d, outcome, err, start) }()
	if err = d.wait(ctx); err != nil {
		outcome = "canceled"
		return nil, status.FromContextError(err).Err()
	}
	switch d.fail {
	case "error":
		return nil, status.Errorf(codes.Unavailable, "authz server error from %s", s.name)
	case "reset":
		outcome = "unsupported-reset"
		return nil, status.Error(codes.Unimplemented, "gRPC connection reset is unsupported; use x-authz-fail: error or stop the server")
	}
	headers := []*corev3.HeaderValueOption{responseHeader("x-authz-server", s.name)}
	if d.tag != "" {
		headers = append(headers, responseHeader("x-authz-tag", d.tag))
	}
	if d.status == http.StatusOK {
		headers = append(headers, responseHeader("x-authz-verdict", "allow"))
		return &authv3.CheckResponse{
			Status: &rpcstatus.Status{Code: int32(codes.OK)},
			HttpResponse: &authv3.CheckResponse_OkResponse{OkResponse: &authv3.OkHttpResponse{
				Headers: headers,
			}},
		}, nil
	}
	return &authv3.CheckResponse{
		Status: &rpcstatus.Status{Code: int32(codes.PermissionDenied)},
		HttpResponse: &authv3.CheckResponse_DeniedResponse{DeniedResponse: &authv3.DeniedHttpResponse{
			Status: &typev3.HttpStatus{Code: typev3.StatusCode(d.status)},
			Body:   s.deniedBody(d.status), Headers: headers,
		}},
	}, nil
}
