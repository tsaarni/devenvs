

# start new cluster
kind delete cluster --name keycloak
kind create cluster --config configs/kind-cluster-config.yaml --name keycloak

# Build mariadb image with vault plugin
docker build docker/mariadb/ -t localhost/mariadb:latest
kind load docker-image --name keycloak localhost/mariadb:latest


kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

kubectl create secret tls keycloak-external --cert=certs/keycloak-server.pem --key=certs/keycloak-server-key.pem --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f manifests/mariadb.yaml
kubectl apply -f manifests/keycloak-26.yaml



kubectl delete -f manifests/keycloak-26.yaml --force
kubectl delete -f manifests/mariadb.yaml --force
kubectl delete pvc -l app=mariadb



https://keycloak.127-0-0-121.nip.io
http://phpmyadmin.127-0-0-121.nip.io



docker run --rm -it --env MARIADB_ROOT_PASSWORD=mariadb mariadb:11.5




https://mariadb.com/kb/en/file-key-management-encryption-plugin/
The File Key Management plugin does not currently support key rotation. See MDEV-20713 for more information.
https://github.com/MariaDB/server/tree/main/plugin/file_key_management



https://mariadb.com/kb/en/hashicorp-key-management-plugin/

https://github.com/faizalrf/documentation/blob/master/MariaDB-Encryption-Hashicorp-Vault.md


https://github.com/MariaDB/server/tree/main/plugin/hashicorp_key_management

https://mariadb.com/resources/blog/mariadb-encryption-tde-using-mariadbs-file-key-management-encryption-plugin/


kubectl exec -it mariadb-0 -- bash

kubectl exec -it mariadb-0 -- cat /var/lib/mysql/keycloak/CLIENT.ibd | strings


get_admin_token() {
  http --form POST http://keycloak.127-0-0-121.nip.io/realms/master/protocol/openid-connect/token username=admin password=admin grant_type=password client_id=admin-cli | jq -r .access_token
}


http POST http://keycloak.127-0-0-121.nip.io/admin/realms/master/clients Authorization:"bearer $(get_admin_token)" clientId=myclient name=myclient redirectUris:='["http://localhost"]' secret=mysecret





kubectl delete -f manifests/openbao.yaml --force
kubectl delete pod -l app=openbao --force --grace-period=0


kubectl apply -f manifests/openbao.yaml
kubectl logs $(kubectl get pod -l app=openbao -o jsonpath="{.items[0].metadata.name}") -c configure


kubectl exec -it $(kubectl get pod -l app=openbao -o jsonpath="{.items[0].metadata.name}") -c configure -- cat /unseal/init.json



kubectl exec -it $(kubectl get pod -l app=openbao -o jsonpath="{.items[0].metadata.name}") -c configure -- ash
export BAO_TOKEN=$(jq -r .root_token /unseal/init.json)

http POST http://openbao:8200/v1/auth/kubernetes/login role=my-role jwt=@/var/run/secrets/kubernetes.io/serviceaccount/token
http POST http://openbao:8200/v1/auth/kubernetes/login role=my-role jwt=@/projected/token lease=1h

http GET http://openbao:8200/v1/secret/mariadb-key X-Vault-Token:$BAO_TOKEN


http GET http://openbao:8200/v1/sys/mounts/secret/mariadb-key/tune X-Vault-Token:$BAO_TOKEN




SELECT CASE WHEN INSTR(NAME, '/') = 0
                   THEN '01-SYSTEM TABLESPACES'
                   ELSE CONCAT('02-', SUBSTR(NAME, 1, INSTR(NAME, '/')-1)) END
                     AS "Schema Name",
         SUM(CASE WHEN ENCRYPTION_SCHEME > 0 THEN 1 ELSE 0 END) "Tables Encrypted",
         SUM(CASE WHEN ENCRYPTION_SCHEME = 0 THEN 1 ELSE 0 END) "Tables Not Encrypted"
FROM information_schema.INNODB_TABLESPACES_ENCRYPTION
GROUP BY CASE WHEN INSTR(NAME, '/') = 0
                   THEN '01-SYSTEM TABLESPACES'
                   ELSE CONCAT('02-', SUBSTR(NAME, 1, INSTR(NAME, '/')-1)) END
ORDER BY 1;
;
