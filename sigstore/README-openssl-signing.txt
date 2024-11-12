

# Generate certificates for the registry
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml



# Add trusted certificate for Docker client
sudo mkdir -p /etc/docker/certs.d/registry.127-0-10-80.nip.io
sudo chown $USER /etc/docker/certs.d/registry.127-0-10-80.nip.io
cp certs/ca.pem /etc/docker/certs.d/registry.127-0-10-80.nip.io/ca.crt



# setting CA certificate for cosign
# See go containerregistry issue
#    https://github.com/google/go-containerregistry/issues/211
export SSL_CERT_FILE=certs/ca.pem


# generate key-pair
go run github.com/sigstore/cosign/v2/cmd/cosign@v2.4.1 generate-key-pair

# Output:
# Enter password for private key:
# Enter password for private key again:
# Private key written to cosign.key
# Public key written to cosign.pub



# Run two instances of registry
docker run --rm -p 127.0.10.80:443:443 -v $PWD/certs:/certs:ro -e REGISTRY_HTTP_ADDR=0.0.0.0:443 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/registry.pem -e REGISTRY_HTTP_TLS_KEY=/certs/registry-key.pem registry:2
docker run --rm -p 127.0.10.81:443:443 -v $PWD/certs:/certs:ro -e REGISTRY_HTTP_ADDR=0.0.0.0:443 -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/registry.pem -e REGISTRY_HTTP_TLS_KEY=/certs/registry-key.pem registry:2


docker pull alpine:3.20.3
docker tag alpine:3.20.3 registry.127-0-10-80.nip.io/alpine:3.20.3
docker push registry.127-0-10-80.nip.io/alpine:3.20.3


cosign generate registry.127-0-10-80.nip.io/alpine:3.20.3 > payload.json

openssl ecparam -name prime256v1 -genkey -noout -out openssl.key
openssl ec -in openssl.key -pubout -out openssl.pub


openssl dgst -sha256 -sign openssl.key -out payload.sig payload.json
cat payload.sig | base64 > payloadbase64.sig


cosign attach signature --payload payload.json --signature payloadbase64.sig registry.127-0-10-80.nip.io/alpine:3.20.3




docker push registry.127-0-10-80.nip.io/alpine:3.20.3
cosign generate registry.127-0-10-80.nip.io/alpine:3.20.3 > payload.json

openssl dgst -sha256 -sign certs/signer-key.pem -out payload.sig payload.json
base64 < payload.sig > payloadbase64.sig
cosign attach signature --payload payload.json --signature payloadbase64.sig registry.127-0-10-80.nip.io/alpine:3.20.3 --certificate=certs/signer.pem  --certificate-chain=certs/sw-sign-sub-ca.pem

SIGSTORE_ROOT_FILE=certs/sw-sign-root-ca.pem cosign verify --certificate-identity-regexp '.*' --certificate-oidc-issuer-regexp '.*' --insecure-ignore-sct --private-infrastructure registry.127-0-10-80.nip.io/alpine:3.20.3

cosign verify --ca-roots=certs/sw-sign-root-ca.pem --certificate-identity-regexp '.*' --certificate-oidc-issuer-regexp '.*' --insecure-ignore-sct --private-infrastructure registry.127-0-10-80.nip.io/alpine:3.20.3


cosign verify --ca-roots=certs/sw-sign-root-ca.pem --ca-intermediates=certs/sw-sign-sub-ca.pem --certificate-identity-regexp '.*' --certificate-oidc-issuer-regexp '.*' --insecure-ignore-sct --private-infrastructure registry.127-0-10-80.nip.io/alpine:3.20.3


cosign verify \
￼    --cert keys/leaf.crt \
￼    --cert-chain keys/certificate_chain.pem \
￼    --certificate-identity-regexp '.*' \
￼    --certificate-oidc-issuer-regexp '.*' \
￼    --private-infrastructure \
￼    --insecure-ignore-sct \
￼    "127.0.0.1:5003/alpine:3.20.3"


openssl dgst -sha256 -sign SSW_DEV_OK_CONT_202101010000_C1_2_1.key -out payload.sig payload.json

cosign attach signature --payload payload.json --signature payloadbase64.sig registry.127-0-10-80.nip.io/alpine:3.20.3 --certificate=SSW_DEV_OK_CONT_202101010000_C1_2_1.pem

cosign verify --ca-roots=SSW_DEV_OK_CONT_202101010000_C1_2_1.pem  --certificate-identity-regexp '.*' --certificate-oidc-issuer-regexp '.*' --insecure-ignore-sct --private-infrastructure registry.127-0-10-80.nip.io/alpine:3.20.3


### Tracing
### https://gist.github.com/tsaarni/06b06e18b614468caa5b522c85d0c61b

### Signing with openssl
### https://docs.sigstore.dev/cosign/signing/signing_with_containers/






cd ~/work/cosign
patch -p1 < ~/work/devenvs/sigstore/cosign-wireshark-keys.patch

wireshark -i lo -f "port 443" -Y http -k -o tls.keylog_file:wireshark-keys.log




GET /v2/ HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT

{}GET /v2/alpine/manifests/3.20.3 HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept: application/vnd.docker.distribution.manifest.v1+json,application/vnd.docker.distribution.manifest.v1+prettyjws,application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json,application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 528
Content-Type: application/vnd.docker.distribution.manifest.v2+json
Docker-Content-Digest: sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735"
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT

{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": 1471,
      "digest": "sha256:91ef0af61f39ece4d6710e465df5ed6ca12112358344fd51ae6a3b886634148b"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 3623807,
         "digest": "sha256:43c4264eed91be63b206e17d93e75256a6097070ce643c5e8f0379998b44f170"
      }
   ]
}GET /v2/alpine/manifests/sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735 HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept: application/vnd.docker.distribution.manifest.v1+json,application/vnd.docker.distribution.manifest.v1+prettyjws,application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json,application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 528
Content-Type: application/vnd.docker.distribution.manifest.v2+json
Docker-Content-Digest: sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735"
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT

{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": 1471,
      "digest": "sha256:91ef0af61f39ece4d6710e465df5ed6ca12112358344fd51ae6a3b886634148b"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 3623807,
         "digest": "sha256:43c4264eed91be63b206e17d93e75256a6097070ce643c5e8f0379998b44f170"
      }
   ]
}GET /v2/alpine/manifests/sha256-33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735.sig HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept: application/vnd.docker.distribution.manifest.v1+json,application/vnd.docker.distribution.manifest.v1+prettyjws,application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json,application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json
Accept-Encoding: gzip

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT
Content-Length: 165

{"errors":[{"code":"MANIFEST_UNKNOWN","message":"manifest unknown","detail":{"Tag":"sha256-33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735.sig"}}]}
GET /v2/ HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT

{}HEAD /v2/alpine/manifests/sha256-33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735.sig HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept: application/vnd.docker.distribution.manifest.v1+json,application/vnd.docker.distribution.manifest.v1+prettyjws,application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json,application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT
Content-Length: 165

HEAD /v2/alpine/blobs/sha256:d748d2b76ebf8a245176858c537a245f62f4cd9271dcad170c107bd399384703 HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT
Content-Length: 157

HEAD /v2/alpine/blobs/sha256:4306feade2f8d7aa22ad654dfd61d0ee66a94f8970fba8d825f45641e9087a78 HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2

HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT
Content-Length: 157

POST /v2/alpine/blobs/uploads/ HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Content-Length: 0
Content-Type: application/json
Accept-Encoding: gzip

HTTP/1.1 202 Accepted
Content-Length: 0
Docker-Distribution-Api-Version: registry/2.0
Docker-Upload-Uuid: e0bc6061-1316-45a3-a15f-84515edd8c0f
Location: https://registry.127-0-10-80.nip.io/v2/alpine/blobs/uploads/e0bc6061-1316-45a3-a15f-84515edd8c0f?_state=dFHza43FnPJwRDJZlTYK1Qbgp3Y3Hh9djFUDczgKsqB7Ik5hbWUiOiJhbHBpbmUiLCJVVUlEIjoiZTBiYzYwNjEtMTMxNi00NWEzLWExNWYtODQ1MTVlZGQ4YzBmIiwiT2Zmc2V0IjowLCJTdGFydGVkQXQiOiIyMDI0LTEwLTI4VDA5OjAxOjU4Ljc0NzIxODI2N1oifQ%3D%3D
Range: 0-0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:01:58 GMT






GET /v2/ HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 2
Content-Type: application/json; charset=utf-8
Docker-Distribution-Api-Version: registry/2.0
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:03:31 GMT

{}GET /v2/alpine/manifests/3.20.3 HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept: application/vnd.docker.distribution.manifest.v1+json,application/vnd.docker.distribution.manifest.v1+prettyjws,application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json,application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 528
Content-Type: application/vnd.docker.distribution.manifest.v2+json
Docker-Content-Digest: sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735"
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:03:31 GMT

{
   "schemaVersion": 2,
   "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
   "config": {
      "mediaType": "application/vnd.docker.container.image.v1+json",
      "size": 1471,
      "digest": "sha256:91ef0af61f39ece4d6710e465df5ed6ca12112358344fd51ae6a3b886634148b"
   },
   "layers": [
      {
         "mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
         "size": 3623807,
         "digest": "sha256:43c4264eed91be63b206e17d93e75256a6097070ce643c5e8f0379998b44f170"
      }
   ]
}GET /v2/alpine/manifests/sha256-33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735.sig HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept: application/vnd.docker.distribution.manifest.v1+json,application/vnd.docker.distribution.manifest.v1+prettyjws,application/vnd.docker.distribution.manifest.v2+json,application/vnd.oci.image.manifest.v1+json,application/vnd.docker.distribution.manifest.list.v2+json,application/vnd.oci.image.index.v1+json
Accept-Encoding: gzip

HTTP/1.1 200 OK
Content-Length: 562
Content-Type: application/vnd.oci.image.manifest.v1+json
Docker-Content-Digest: sha256:36570d59b0a2ce9a2ae9da089e72b45c4b51c11bf9a45050dd8cd42e4cb638f4
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:36570d59b0a2ce9a2ae9da089e72b45c4b51c11bf9a45050dd8cd42e4cb638f4"
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:03:31 GMT

{"schemaVersion":2,"mediaType":"application/vnd.oci.image.manifest.v1+json","config":{"mediaType":"application/vnd.oci.image.config.v1+json","size":233,"digest":"sha256:d748d2b76ebf8a245176858c537a245f62f4cd9271dcad170c107bd399384703"},"layers":[{"mediaType":"application/vnd.dev.cosign.simplesigning.v1+json","size":250,"digest":"sha256:4306feade2f8d7aa22ad654dfd61d0ee66a94f8970fba8d825f45641e9087a78","annotations":{"dev.cosignproject.cosign/signature":"MEUCIQDBRZD5sNATW5T+7IFug6L4I9gfaNu77Hvf/p15V1PDTQIgTTIFoW+H8y/ZSSNZD5SdNA0b\nia5Zhsd26RjgvPb6CXU=\n"}}]}GET /v2/alpine/blobs/sha256:4306feade2f8d7aa22ad654dfd61d0ee66a94f8970fba8d825f45641e9087a78 HTTP/1.1
Host: registry.127-0-10-80.nip.io
User-Agent: cosign/devel (linux; amd64) go-containerregistry/v0.20.2
Accept-Encoding: gzip

HTTP/1.1 200 OK
Accept-Ranges: bytes
Cache-Control: max-age=31536000
Content-Length: 250
Content-Type: application/octet-stream
Docker-Content-Digest: sha256:4306feade2f8d7aa22ad654dfd61d0ee66a94f8970fba8d825f45641e9087a78
Docker-Distribution-Api-Version: registry/2.0
Etag: "sha256:4306feade2f8d7aa22ad654dfd61d0ee66a94f8970fba8d825f45641e9087a78"
X-Content-Type-Options: nosniff
Date: Mon, 28 Oct 2024 09:03:31 GMT

{"critical":{"identity":{"docker-reference":"registry.127-0-10-80.nip.io/alpine"},"image":{"docker-manifest-digest":"sha256:33735bd63cf84d7e388d9f6d297d348c523c044410f553bd878c6d7829612735"},"type":"cosign container image signature"},"optional":null}
