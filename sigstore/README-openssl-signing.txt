
### Reverse engineering cosign detached signatures
### https://gist.github.com/tsaarni/06b06e18b614468caa5b522c85d0c61b

### Signing with openssl
### https://docs.sigstore.dev/cosign/signing/signing_with_containers/


# setting CA certificate for cosign
# See go containerregistry issue
#    https://github.com/google/go-containerregistry/issues/211
export SSL_CERT_FILE=certs/ca.pem




# Generate certificates for the registry
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml



# Add trusted certificate for Docker client
sudo mkdir -p /etc/docker/certs.d/registry.127-0-10-80.nip.io
sudo chown $USER /etc/docker/certs.d/registry.127-0-10-80.nip.io
cp certs/ca.pem /etc/docker/certs.d/registry.127-0-10-80.nip.io/ca.crt



# Create key pair
openssl ecparam -name prime256v1 -genkey -noout -out openssl.key
openssl ec -in openssl.key -pubout -out openssl.pub


# Run container registry
docker run --rm -p 127.0.10.80:443:443 -v $PWD/certs:/certs:ro -e REGISTRY_HTTP_ADDR=0.0.0.0:443 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/registry.pem -e REGISTRY_HTTP_TLS_KEY=/certs/registry-key.pem registry:2


# Pull example image
docker pull alpine:3.20.3
docker tag alpine:3.20.3 registry.127-0-10-80.nip.io/alpine:3.20.3




docker push registry.127-0-10-80.nip.io/alpine:3.20.3

cosign generate registry.127-0-10-80.nip.io/alpine:3.20.3 > payload.json

openssl dgst -sha256 -sign certs/signer-key.pem -out payload.sig payload.json
base64 < payload.sig > payloadbase64.sig

# For some reason, the root certificate needs to be part of the chain stored with the signature, even though it should not be needed for verification.
cat certs/sw-sign-sub-ca.pem certs/sw-sign-root-ca.pem > certs/sw-sign-sub-ca-including-root.pem

cosign attach signature --payload payload.json --signature payloadbase64.sig registry.127-0-10-80.nip.io/alpine:3.20.3 --certificate=certs/signer.pem  --certificate-chain=certs/sw-sign-sub-ca-including-root.pem

# --private-infrastructure  - disables transparency log
# --insecure-ignore-sct     - ignore Signed Certificate Timestamp

cosign verify --ca-roots=certs/sw-sign-root-ca.pem --certificate-identity-regexp '.*' --certificate-oidc-issuer-regexp '.*' --private-infrastructure --insecure-ignore-sct registry.127-0-10-80.nip.io/alpine:3.20.3





## Sniffing cosign traffic

cd ~/work/cosign
patch -p1 < ~/work/devenvs/sigstore/cosign-wireshark-keys.patch

wireshark -i lo -f "port 443" -Y http -k -o tls.keylog_file:wireshark-keys.log


####
## Lack of revocation check, CRL

# FR: BYO PKI Revocation: CA-issued cert
https://github.com/sigstore/cosign/issues/2568

# Add interface for certificate validation
https://github.com/sigstore/sigstore-go/issues/298
