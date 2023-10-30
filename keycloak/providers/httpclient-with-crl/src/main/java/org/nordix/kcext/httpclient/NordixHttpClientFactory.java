/*
 * Copyright 2016 Red Hat, Inc. and/or its affiliates
 * and other contributors as indicated by the @author tags.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package org.nordix.kcext.httpclient;

import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.client.entity.EntityBuilder;
import org.apache.http.client.methods.HttpGet;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.conn.ssl.NoopHostnameVerifier;
import org.apache.http.conn.ssl.SSLConnectionSocketFactory;
import org.apache.http.entity.ContentType;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.impl.client.HttpClients;
import org.jboss.logging.Logger;
import org.keycloak.Config;
import org.keycloak.connections.httpclient.HttpClientFactory;
import org.keycloak.connections.httpclient.HttpClientProvider;
import org.keycloak.models.KeycloakSession;
import org.keycloak.models.KeycloakSessionFactory;
import org.keycloak.provider.ProviderConfigProperty;
import org.keycloak.provider.ProviderConfigurationBuilder;
import org.keycloak.truststore.TruststoreProvider;

import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.security.InvalidAlgorithmParameterException;
import java.security.KeyManagementException;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.security.cert.CRL;
import java.security.cert.CRLException;
import java.security.cert.CertPathBuilder;
import java.security.cert.CertStore;
import java.security.cert.CertificateException;
import java.security.cert.CertificateFactory;
import java.security.cert.CollectionCertStoreParameters;
import java.security.cert.PKIXBuilderParameters;
import java.security.cert.PKIXRevocationChecker;
import java.security.cert.X509CertSelector;
import java.util.ArrayList;
import java.util.Collection;
import java.util.EnumSet;
import java.util.HashSet;
import java.util.List;

import javax.net.ssl.CertPathTrustManagerParameters;
import javax.net.ssl.SSLContext;
import javax.net.ssl.TrustManagerFactory;

import org.apache.http.client.methods.CloseableHttpResponse;
import org.apache.http.util.EntityUtils;

public class NordixHttpClientFactory implements HttpClientFactory {

    private static final Logger log = Logger.getLogger(NordixHttpClientFactory.class);
    private String crlFile;

    @Override
    public HttpClientProvider create(KeycloakSession session) {
        try {
            return new ClientProvider(createHttpClient(session));
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void close() {
    }

    @Override
    public String getId() {
        return "nordix";
    }

    @Override
    public void init(Config.Scope config) {
        crlFile = config.get("crl-file");
        log.debugv("CRL file configured crlFile={0}", crlFile == null ? "none" : crlFile);
    }


    @Override
    public void postInit(KeycloakSessionFactory factory) {

    }

    @Override
    public List<ProviderConfigProperty> getConfigMetadata() {
        return ProviderConfigurationBuilder.create()
                .property()
                .name("crl-file")
                .type("string")
                .helpText("The file path of the CRL file to use for certificate revocation checking.")
                .add()
                .build();
    }

    private class ClientProvider implements HttpClientProvider {

        private CloseableHttpClient httpClient;

        public ClientProvider(CloseableHttpClient httpClient) {
            this.httpClient = httpClient;
        }


        @Override
        public CloseableHttpClient getHttpClient() {
            return httpClient;
        }

        @Override
        public void close() {
            try {
                httpClient.close();
            } catch (IOException e) {
                e.printStackTrace();
            }
        }

        @Override
        public int postText(String uri, String text) throws IOException {
            log.debugv("POST uri={0}", uri);
            HttpPost request = new HttpPost(uri);
            request.setEntity(EntityBuilder.create().setText(text).setContentType(ContentType.TEXT_PLAIN).build());
            try (CloseableHttpResponse response = httpClient.execute(request)) {
                try {
                    return response.getStatusLine().getStatusCode();
                } finally {
                    EntityUtils.consumeQuietly(response.getEntity());
                }
            } catch (Exception t) {
                log.warn(t.getMessage(), t);
                throw t;
            }
        }

        @Override
        public InputStream get(String uri) throws IOException {
            log.debugv("GET uri={0}", uri);
            HttpGet request = new HttpGet(uri);
            HttpResponse response = httpClient.execute(request);
            int statusCode = response.getStatusLine().getStatusCode();
            HttpEntity entity = response.getEntity();
            if (statusCode < 200 || statusCode >= 300) {
                EntityUtils.consumeQuietly(entity);
                throw new IOException("Unexpected HTTP status code " + response.getStatusLine().getStatusCode() + " when expecting 2xx");
            }
            if (entity == null) {
                throw new IOException("No content returned from HTTP call");
            }
            return entity.getContent();

        }
    }

    private CloseableHttpClient createHttpClient(KeycloakSession session) throws IOException, KeyStoreException, InvalidAlgorithmParameterException, CRLException, CertificateException, NoSuchAlgorithmException, KeyManagementException {
        log.debugv("Creating new HttpClient");

        TruststoreProvider trustStoreProvider = session.getProvider(TruststoreProvider.class);
        if (trustStoreProvider == null) {
            log.error("TruststoreProvider is null");
            throw new RuntimeException("Truststore SPI is not configured");
        }

        KeyStore trustStore = trustStoreProvider.getTruststore();
        if (trustStore == null) {
            log.error("Truststore is null");
            throw new RuntimeException("Truststore not configured for Truststore SPI");
        }

        PKIXBuilderParameters pkixParams = new PKIXBuilderParameters(trustStore, new X509CertSelector());

        if (crlFile != null) {
            Collection<CRL> crls = new HashSet<>();
            crls.add(CertificateFactory.getInstance("X.509").generateCRL(Files.newInputStream(Paths.get(crlFile))));

            List<CertStore> certStores = new ArrayList<>();
            certStores.add(CertStore.getInstance("Collection", new CollectionCertStoreParameters(crls)));

            PKIXRevocationChecker revocationChecker = (PKIXRevocationChecker) CertPathBuilder.getInstance("PKIX").getRevocationChecker();
            revocationChecker.setOptions(
                    EnumSet.of(PKIXRevocationChecker.Option.PREFER_CRLS, PKIXRevocationChecker.Option.NO_FALLBACK));

            pkixParams.setCertStores(certStores);
            pkixParams.addCertPathChecker(revocationChecker);
        }

        TrustManagerFactory tmf = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
        tmf.init(new CertPathTrustManagerParameters(pkixParams));

        SSLContext context = SSLContext.getInstance("TLS");
        context.init(null, tmf.getTrustManagers(), new SecureRandom());

        SSLConnectionSocketFactory sslConnectionSocketFactory = new SSLConnectionSocketFactory(context,
                new NoopHostnameVerifier());

        return HttpClients.custom().setSSLSocketFactory(sslConnectionSocketFactory).build();
    }

}
