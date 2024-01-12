





1. Deploy on Kubernetes on Kind

2. Create hosts entry for Keycloak

    Check the cluster IP address of Envoy and add coredns hosts entry

kubectl -n projectcontour get service envoy -o jsonpath='{.spec.clusterIP}'
kubectl -n kube-system edit configmap coredns

# append following to config

hosts inline nip.io {
    [ENVOY_CLUSTER_IP] keycloak.127-0-0-121.nip.io
}


3. Create another realm called "secondrealm" which will be used for federated login

4. Create client "federator" in "secondrealm" and take a note on client secret

5. Create user "joe" in "secondrealm"

6. In "master" realm add new Keycloak OpenID identity provider with following discovery endpoint
https://keycloak.127-0-0-121.nip.io/realms/secondrealm/.well-known/openid-configuration

7. Use incognito window to login to account console using federated user account

https://keycloak.127-0-0-121.nip.io/realms/master/account/





*** Decrypt TLS traffic


# Copy the agent to providers subdirectory (mounted as host volume)
cp ~/package/extract-tls-secrets/target/extract-tls-secrets-4.1.0-SNAPSHOT.jar providers

# Add agent to keycloak deployment
    - name: JAVA_OPTS_APPEND
      value: "... -javaagent:/opt/keycloak/providers/extract-tls-secrets-4.1.0-SNAPSHOT.jar=/tmp/wireshark-keys.log"


sudo nsenter --target $(pgrep -f io.quarkus.bootstrap.runner.QuarkusEntryPoint) --net wireshark -i any -k -o tls.keylog_file:/proc/$(pgrep -f io.quarkus.bootstrap.runner.QuarkusEntryPoint)/root/tmp/wireshark-keys.log



*** Swap valid and expired certificates


kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret tls keycloak-external --cert=certs/keycloak-server-expired.pem --key=certs/keycloak-server-expired-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl logs statefulset/keycloak
