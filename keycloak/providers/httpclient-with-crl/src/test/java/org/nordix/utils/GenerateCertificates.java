package org.nordix.utils;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.KeyStore;
import java.util.List;

import fi.protonode.certy.CertificateRevocationList;
import fi.protonode.certy.Credential;
import fi.protonode.certy.Credential.KeyType;

public class GenerateCertificates {

    public static void main(String[] args) throws Exception {

        Credential externalCa = new Credential().subject("cn=external-ca");

        // Keycloak server certificate
        new Credential().subject("cn=keycloak-server")
                .subjectAltNames(
                        List.of("DNS:localhost", "DNS:keycloak.127-0-0-1.nip.io", "DNS:keycloak.127-0-0-121.nip.io"))
                .issuer(externalCa).writeCertificatesAsPem(Paths.get("keycloak-server.pem"))
                .writePrivateKeyAsPem(Paths.get("keycloak-server-key.pem"));

        // IDP server certificate
        Credential idpServer = new Credential().subject("cn=idp").subjectAltName("DNS:idp").issuer(externalCa).crlDistributionPointUri("http://shell/crl.pem")
                .writeCertificatesAsPem(Paths.get("idp-server.pem"))
                .writePrivateKeyAsPem(Paths.get("idp-server-key.pem"));

        // LDAP server certificate
        Credential ldapServer = new Credential().subject("cn=ldap").subjectAltNames(List.of("DNS:localhost", "DNS:ldap", "DNS:ldap.127-0-0-1.nip.io"))
                        .issuer(externalCa).crlDistributionPointUri("http://shell/crl.pem")
                .keyType(KeyType.RSA)
                .writeCertificatesAsPem(Paths.get("ldap.pem"))
                .writePrivateKeyAsPem(Paths.get("ldap-key.pem"));

        // Truststore for Keycloak HTTPClient and LDAP client
        KeyStore trustStore = KeyStore.getInstance("PKCS12");
        trustStore.load(null, null);
        trustStore.setCertificateEntry("ca", externalCa.getCertificate());
        trustStore.store(Files.newOutputStream(Paths.get("truststore.p12")), "secret".toCharArray());

        // CLRs for Keycloak HTTPClient and LDAP client
        new CertificateRevocationList().issuer(externalCa).writeAsPem(Paths.get("crl-server-not-revoked.pem"));
        new CertificateRevocationList().issuer(externalCa).add(idpServer).add(ldapServer)
                .writeAsPem(Paths.get("crl-server-revoked.pem"));
        Credential unrelatedCa = new Credential().subject("cn=unrelated-ca").generate();
        new CertificateRevocationList().issuer(unrelatedCa).writeAsPem(Paths.get("crl-unrelated.pem"));

        new Credential().subject("cn=client-ca").writeCertificatesAsPem(Paths.get("client-ca.pem"));
    }
}
