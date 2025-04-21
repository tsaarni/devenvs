

kubectl -n projectcontour scale deployment contour --replicas=1


cat > certs.yaml <<EOF
subject: cn=internal-ca
key_type: EC
key_size: 384
---
subject: cn=contour
issuer: cn=internal-ca
sans:
- DNS:contour
key_type: EC
key_size: 384
---
subject: cn=envoy
issuer: cn=internal-ca
key_type: EC
key_size: 384
EOF



cat > certs.yaml <<EOF
subject: cn=internal-ca
key_type: EC
key_size: 521
---
subject: cn=contour
issuer: cn=internal-ca
sans:
- DNS:contour
key_type: EC
key_size: 384
---
subject: cn=envoy
issuer: cn=internal-ca
key_type: EC
key_size: 521
EOF


rm *.pem

certyaml


kubectl -n projectcontour create secret generic envoycert --from-file=tls.crt=envoy.pem --from-file=tls.key=envoy-key.pem --from-file=ca.crt=internal-ca.pem --dry-run=client -o yaml | kubectl apply -f -
kubectl -n projectcontour create secret generic contourcert --from-file=tls.crt=contour.pem --from-file=tls.key=contour-key.pem --from-file=ca.crt=internal-ca.pem --dry-run=client -o yaml | kubectl apply -f -

# restart contour and envoy
kubectl -n projectcontour rollout restart deployment contour
kubectl -n projectcontour rollout restart daemonset envoy



kubectl -n projectcontour get pod

kubectl -n projectcontour logs daemonset/envoy -c envoy -f
kubectl -n projectcontour logs deployment/contour -c contour
