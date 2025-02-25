
# https://hub.docker.com/r/keyfactor/ejbca-ce
# https://github.com/Keyfactor/ejbca-ce/tree/main/charts/ejbca
# https://docs.keyfactor.com/container/latest/ejbca/ejbca-helm-deployment-parameters
#

rm -rf certs
mkdir certs
certyaml -d certs configs/certs.yaml


openssl x509 -in certs/managementCA.pem -out certs/ManagementCA.crt -outform DER
chmod a+r certs/ManagementCA.crt
openssl pkcs12 -export -passout pass:secret -noiter -nomaciter -in certs/ejbca-server.pem -inkey certs/ejbca-server-key.pem -out certs/ejbca-server.p12
chmod a+r certs/ejbca-server.p12
echo secret > certs/ejbca-server.storepasswd



# Run EJBCA with custom imported ManagementCA
docker run -it --rm -p 8443:8443 -e TLS_SETUP_ENABLED="true" -e HTTPSERVER_HOSTNAME="ejbca.127.0.0.1.nip.io" -e ADMINWEB_ACCESS="true" -e INITIAL_ADMIN="ManagementCA;WITH_COMMONNAME;SuperAdmin" -e APPSERVER_KEYSTORE_SECRET="secret" -v $PWD/certs/ManagementCA.crt:/mnt/external/secrets/tls/cas/ManagementCA.crt -v $PWD/certs/ejbca-server.p12:/mnt/external/secrets/tls/ks/server.jks -v $PWD/certs/ejbca-server.storepasswd:/mnt/external/secrets/tls/ks/server.storepasswd keyfactor/ejbca-ce:9.0.0


https://ejbca.127.0.0.1.nip.io:8443/ejbca/adminweb/


http -v --verify=certs/server-ca.pem --cert certs/SuperAdmin.pem --cert-key certs/SuperAdmin-key.pem https://ejbca.127.0.0.1.nip.io:8443/ejbca/adminweb/
