package main

import (
	"context"
	"fmt"
	"sort"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/encoding/protojson"
)

func formatHeaders(h map[string]string) string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	s := ""
	for _, k := range keys {
		s += fmt.Sprintf("\n  %s: %s", k, h[k])
	}
	return s
}

var protoFmt = protojson.MarshalOptions{Multiline: true, EmitDefaultValues: true}

type grpcAuthzServer struct {
	authv3.UnimplementedAuthorizationServer
}

func (s *grpcAuthzServer) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	log.Info("gRPC Check request")
	fmt.Println(protoFmt.Format(req))

	// Log TLS info and client certificate if mTLS is used.
	if p, ok := peer.FromContext(ctx); ok && p.AuthInfo != nil {
		if tlsInfo, ok := p.AuthInfo.(credentials.TLSInfo); ok {
			log.Info("gRPC TLS info", "sni", tlsInfo.State.ServerName)
			for _, cert := range tlsInfo.State.PeerCertificates {
				log.Info("gRPC client certificate", "subject", cert.Subject)
			}
		}
	}

	// TODO: Add your authorization logic here.

	// Example denied response:
	// resp := &authv3.CheckResponse{
	// 	Status: &status.Status{Code: int32(code.Code_PERMISSION_DENIED)},
	// 	HttpResponse: &authv3.CheckResponse_DeniedResponse{
	// 		DeniedResponse: &authv3.DeniedHttpResponse{
	// 			Status: &corev3.HttpStatus{Code: corev3.StatusCode_Forbidden},
	// 			Body:   "Access denied",
	// 		},
	// 	},
	// }

	resp := &authv3.CheckResponse{
		Status: &status.Status{Code: int32(code.Code_OK)},
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				// Headers: add or overwrite request headers forwarded upstream.
				// Headers: []*corev3.HeaderValueOption{
				// 	{Header: &corev3.HeaderValue{Key: "x-user", Value: "alice"}},
				// },
				// HeadersToRemove: strip request headers before forwarding upstream.
				// HeadersToRemove: []string{"x-internal-secret"},
				// ResponseHeadersToAdd: add headers to the downstream response.
				//ResponseHeadersToAdd: []*corev3.HeaderValueOption{
				//	{Header: &corev3.HeaderValue{Key: "Set-Cookie", Value: "session=grpc-authz-test"}},
				//},
				// DynamicMetadata: pass structured metadata to other Envoy filters.
				// DynamicMetadata: &structpb.Struct{...},
			},
		},
	}

	log.Info("gRPC Check result")
	fmt.Println(protoFmt.Format(resp))
	return resp, nil
}
