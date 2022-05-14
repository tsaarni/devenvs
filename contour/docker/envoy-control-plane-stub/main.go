package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	envoy_api_v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	// envoy_api_v2_auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_api_v2_endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoy_api_v2_listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoy_api_v2_route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoy_config_filter_network_http_connection_manager_v2 "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoy_service_discovery_v2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	log "github.com/sirupsen/logrus"
)

const bindAddress = ":8080"

func main() {
	log.SetLevel(log.TraceLevel)

	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, log.StandardLogger())
	server := xds.NewServer(context.Background(), snapshotCache, nil)

	// info on certificate reload
	// https://github.com/grpc/grpc-go/issues/2167
	// https://stackoverflow.com/questions/37473201/is-there-a-way-to-update-the-tls-certificates-in-a-net-http-server-without-any-d

	// load gRPC server credentials
	serverCreds, err := tls.LoadX509KeyPair("/run/secrets/certs/controlplane.pem", "/run/secrets/certs/controlplane-key.pem")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// load CA certificate to validate gRPC clients
	cert, err := ioutil.ReadFile("/run/secrets/certs/internal-root-ca.pem")
	if err != nil {
		log.Fatalf("Could not load CA certificates: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCreds},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		VerifyPeerCertificate: func(certificates [][]byte, _ [][]*x509.Certificate) error {
			for _, asn1Data := range certificates {
				cert, err := x509.ParseCertificate(asn1Data)
				if err != nil {
					log.Println(err.Error())
				}
				log.Printf("Subject: %v, Issuer: %v, Serial: %v, NotBefore: %v, NotAfter: %v", cert.Subject, cert.Issuer, cert.SerialNumber, cert.NotBefore, cert.NotAfter)
			}
			return nil
		},
	}

	grpcOpts := []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsConfig))}
	grpcServer := grpc.NewServer(grpcOpts...)

	// creds, err := credentials.NewServerTLSFromFile("/input/certs/controlplane.pem", "/input/certs/controlplane-key.pem")
	// if err != nil {
	// 	log.Fatalf("could not load TLS keys: %s", err)
	// }
	// grpcOpts := []grpc.ServerOption{grpc.Creds(creds)}
	// grpcServer := grpc.NewServer(grpcOpts...)

	lis, err := net.Listen("tcp", bindAddress)
	if err != nil {
		log.Fatalf("Could not bind to TCP port: %s", err)
	}

	envoy_service_discovery_v2.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	envoy_api_v2.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	envoy_api_v2.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	envoy_api_v2.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	envoy_api_v2.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	go func() {
		log.Infof("Listening: %v", bindAddress)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Error while listening: %s", err)
		}
	}()

	// upstreamTLS := &envoy_api_v2_auth.UpstreamTlsContext{
	// 	Sni: "httpbin",
	// }

	// upstreamTLSSerialized, err := ptypes.MarshalAny(upstreamTLS)
	// if err != nil {
	// 	log.Fatalf("Cannot marshal: %s", err)
	// }

	clusters := []cache.Resource{
		&envoy_api_v2.Cluster{
			Name:           "service_httpbin",
			ConnectTimeout: ptypes.DurationProto(1 * time.Second),
			ClusterDiscoveryType: &envoy_api_v2.Cluster_Type{
				Type: envoy_api_v2.Cluster_LOGICAL_DNS,
			},
			DnsLookupFamily: envoy_api_v2.Cluster_V4_ONLY,
			LbPolicy:        envoy_api_v2.Cluster_ROUND_ROBIN,
			LoadAssignment: &envoy_api_v2.ClusterLoadAssignment{
				ClusterName: "service_httpbin",
				Endpoints: []*envoy_api_v2_endpoint.LocalityLbEndpoints{{
					LbEndpoints: []*envoy_api_v2_endpoint.LbEndpoint{{
						HostIdentifier: &envoy_api_v2_endpoint.LbEndpoint_Endpoint{
							Endpoint: &envoy_api_v2_endpoint.Endpoint{
								Address: &envoy_api_v2_core.Address{
									Address: &envoy_api_v2_core.Address_SocketAddress{
										SocketAddress: &envoy_api_v2_core.SocketAddress{
											Address: "httpbin",
											PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
												PortValue: 80,
											},
										},
									},
								},
							},
						},
					}},
				}},
			},
			// TransportSocket: &envoy_api_v2_core.TransportSocket{
			// 	Name: "envoy.transport_sockets.tls",
			// 	ConfigType: &envoy_api_v2_core.TransportSocket_TypedConfig{
			// 		TypedConfig: upstreamTLSSerialized,
			// 	},
			// },
		},
	}

	connectionManager := &envoy_config_filter_network_http_connection_manager_v2.HttpConnectionManager{
		StatPrefix: "ingress_http",
		RouteSpecifier: &envoy_config_filter_network_http_connection_manager_v2.HttpConnectionManager_RouteConfig{
			RouteConfig: &envoy_api_v2.RouteConfiguration{
				Name: "local_route",
				VirtualHosts: []*envoy_api_v2_route.VirtualHost{{
					Name:    "local_service",
					Domains: []string{"*"},
					Routes: []*envoy_api_v2_route.Route{{
						Match: &envoy_api_v2_route.RouteMatch{
							PathSpecifier: &envoy_api_v2_route.RouteMatch_Prefix{
								Prefix: "/",
							},
						},
						Action: &envoy_api_v2_route.Route_Route{
							Route: &envoy_api_v2_route.RouteAction{
								HostRewriteSpecifier: &envoy_api_v2_route.RouteAction_HostRewrite{
									HostRewrite: "httpbin",
								},
								ClusterSpecifier: &envoy_api_v2_route.RouteAction_Cluster{
									Cluster: "service_httpbin",
								},
								Timeout: &duration.Duration{
									Seconds: 60 * 60,
								},
							},
						},
					}},
				}},
			},
		},
		HttpFilters: []*envoy_config_filter_network_http_connection_manager_v2.HttpFilter{{
			Name: wellknown.Router,
		}},
	}

	connectionManagerSerialized, err := ptypes.MarshalAny(connectionManager)
	if err != nil {
		log.Fatalf("Cannot marshal: %s", err)
	}

	listeners := []cache.Resource{
		&envoy_api_v2.Listener{
			Name: "listener",
			Address: &envoy_api_v2_core.Address{
				Address: &envoy_api_v2_core.Address_SocketAddress{
					SocketAddress: &envoy_api_v2_core.SocketAddress{
						Protocol: envoy_api_v2_core.SocketAddress_TCP,
						Address:  "0.0.0.0",
						PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
							PortValue: 80,
						},
					},
				},
			},
			FilterChains: []*envoy_api_v2_listener.FilterChain{{
				Filters: []*envoy_api_v2_listener.Filter{{
					Name: wellknown.HTTPConnectionManager,
					ConfigType: &envoy_api_v2_listener.Filter_TypedConfig{
						TypedConfig: connectionManagerSerialized,
					},
				}},
			}},
		},
	}

	snapshot := cache.NewSnapshot("1.0", nil, clusters, nil, listeners, nil)
	snapshotCache.SetSnapshot("envoy", snapshot)

	select {} // block forever
}
