

# Cosign Signature Specifications
https://github.com/sigstore/cosign/blob/main/specs/SIGNATURE_SPEC.md


### containers project signatures
https://github.com/containers/image/blob/main/docs/containers-signature.5.md


# generate key-pair (will ask for a password)
./cosign generate-key-pair


# key will look like
cat cosign.key
-----BEGIN ENCRYPTED COSIGN PRIVATE KEY-----
....
<JSON payload in base64 encoding>
...
-----END ENCRYPTED COSIGN PRIVATE KEY-----

$ grep -v COSIGN cosign.key | base64 -d | jq .
{
  "kdf": {
    "name": "scrypt",
    "params": {
      "N": 32768,
      "r": 8,
      "p": 1
    },
    "salt": "Sv1R23Sv0BXPGvQUIPFJjvbYPpDv6xA2bWGqcQQ2hr8="
  },
  "cipher": {
    "name": "nacl/secretbox",
    "nonce": "4Va6f2C39Rwn/ZDUApk1jr8RU4a9itd9"
  },
  "ciphertext": "<base64 secret key>"
}




# Generate certificates for the registry
mkdir -p certs
certyaml -d certs configs/certs.yaml


# add generated CA cert as trusted roto for docker
sudo mkdir -p /etc/docker/certs.d/registry.127-0-10-80.nip.io
sudo chown $USER /etc/docker/certs.d/registry.127-0-10-80.nip.io
cp certs/ca.pem /etc/docker/certs.d/registry.127-0-10-80.nip.io/ca.crt


# start two registries
#  - one for signing images on
#  - another to copy signed images to
docker run --rm -p 127.0.10.80:443:443 -v $PWD/certs:/certs:ro -e REGISTRY_HTTP_ADDR=0.0.0.0:443 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/registry.pem -e REGISTRY_HTTP_TLS_KEY=/certs/registry-key.pem registry:2
docker run --rm -p 127.0.10.81:443:443 -v $PWD/certs:/certs:ro -e REGISTRY_HTTP_ADDR=0.0.0.0:443 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/registry.pem -e REGISTRY_HTTP_TLS_KEY=/certs/registry-key.pem registry:2


# copy image to registry for signing
docker pull alpine:latest
docker tag alpine:latest registry.127-0-10-80.nip.io/alpine:latest
docker push registry.127-0-10-80.nip.io/alpine:latest

IMGHASH=$(docker inspect --format="{{(index .RepoDigests 1)}}" registry.127-0-10-80.nip.io/alpine:latest)

# sign image
export SSL_CERT_FILE=/etc/docker/certs.d/registry.127-0-10-80.nip.io/ca.crt
./cosign sign --tlog-upload=false --key cosign.key $IMGHASH
./cosign verify --insecure-ignore-tlog --key cosign.pub $IMGHASH | jq .

# OUTPUT:
#
# Verification for registry.127-0-10-80.nip.io/alpine@sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870 --
# The following checks were performed on each of these signatures:
#  - The cosign claims were validated
#  - The signatures were verified against the specified public key
#
# [{"critical":{"identity":{"docker-reference":"registry.127-0-10-80.nip.io/alpine"},"image":{"docker-manifest-digest":"sha256:1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870"},"type":"cosign container image signature"},"optional":null}]

# copy image AND signature to another registry
./cosign copy registry.127-0-10-80.nip.io/alpine:latest registry.127-0-10-81.nip.io/alpine:latest



sudo apt install skopeo

skopeo inspect --raw docker://$IMGHASH | jq .
skopeo inspect --config --raw docker://$IMGHASH | jq .


http --verify no -v https://registry.127-0-10-80.nip.io/v2/alpine/manifests/sha256-1304f174557314a7ed9eddb4eab12fed12cb0cd9809e4c28f29af86979a3c870.sig Accept:application/vnd.oci.image.manifest.v1+json

HTTP/1.1 200 OK
Content-Length: 558
Content-Type: application/vnd.oci.image.manifest.v1+json
Date: Mon, 31 Oct 2022 14:35:37 GMT
Docker-Content-Digest: sha256:f1afaf1e30ace4a6328bd7f438acd8fe09cc9c7b7ca90f4cce9cd3eb0ca0aec6
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:f1afaf1e30ace4a6328bd7f438acd8fe09cc9c7b7ca90f4cce9cd3eb0ca0aec6"
X-Content-Type-Options: nosniff

{
    "config": {
        "digest": "sha256:bce153781da31151e974bf96d6fa64de13ea15e4a032b24995008c8ce3c1cd9f",
        "mediaType": "application/vnd.oci.image.config.v1+json",
        "size": 248
    },
    "layers": [
        {
            "annotations": {
                "dev.cosignproject.cosign/signature": "MEYCIQDac1xwNMOaXZZ1+N4xA5Yv/Bo39kazqyWVD/VqhS31lwIhAImgsPDlKAqz89hKWDlLXhJTdGRwyKD89xYde4OLHGYb"
            },
            "digest": "sha256:5d87db555c2f211555acaef3ce9120bf7a46c050e8b6930085a7b871aa927e50",
            "mediaType": "application/vnd.dev.cosign.simplesigning.v1+json",
            "size": 250
        }
    ],
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "schemaVersion": 2
}








################
#
# inspect what cosign does with registry
#

# make cosign dump keys for wireshark
patch -p1 < ~/work/devenvs/sigstore/cosign-wireshark-keys.patch

wireshark -i lo -f "port 443" -Y http -k -o tls.keylog_file:$HOME/work/cosign/wireshark-keys.log





################################
#
# Working with "remote keys"
#

# implementations
https://github.com/sigstore/sigstore/tree/main/pkg/signature/kms

# Generalized Cosign External Signer/HSM Interface
https://github.com/sigstore/cosign/issues/396

# PKCS11 signing support
https://github.com/sigstore/cosign/pull/985

# KMS Integrations
https://github.com/sigstore/cosign/blob/main/KMS.md


## for reading
##  https://github.com/salrashid123/go_pkcs11
##  https://github.com/salrashid123/mtls_pkcs11









Using external signer?
https://docs.keyfactor.com/signserver/latest/tutorial-signserver-container-signing-with-cosign



