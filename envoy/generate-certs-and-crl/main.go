package main

import (
	"github.com/tsaarni/certyaml"
)

func main() {
	serverCA := certyaml.Certificate{
		Subject: "cn=server-ca",
	}

	clientCA := certyaml.Certificate{
		Subject: "cn=client-ca",
	}

	envoy := certyaml.Certificate{
		Issuer:  &serverCA,
		Subject: "cn=envoy",
		SubjectAltNames: []string{
			"DNS:localhost",
		},
	}

	client := certyaml.Certificate{
		Issuer:  &clientCA,
		Subject: "cn=client",
	}

	serverCA.WritePEM("../certs/server-ca.pem", "../certs/server-ca-key.pem")
	clientCA.WritePEM("../certs/client-ca.pem", "../certs/client-ca-key.pem")
	envoy.WritePEM("../certs/envoy.pem", "../certs/envoy-key.pem")
	client.WritePEM("../certs/client.pem", "../certs/client-key.pem")

	crlClientNotRevoked := certyaml.CRL{
		Issuer:  &clientCA,
		Revoked: []*certyaml.Certificate{},
	}

	crlClientRevoked := certyaml.CRL{
		Issuer: &clientCA,
		Revoked: []*certyaml.Certificate{
			&client,
		},
	}

	crlClientNotRevoked.WritePEM("../certs/crl-client-not-revoked.pem")
	crlClientRevoked.WritePEM("../certs/crl-client-revoked.pem")
}
