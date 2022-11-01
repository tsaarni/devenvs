


make cosign
make lint

# build with pkcs#11 and piv support
##  replace c_pcsclite with go-libpcsclite:  https://github.com/go-piv/piv-go/issues/82
##  WIP: piv/internal/pcsc: add pure go pcsc client implementation:  https://github.com/go-piv/piv-go/pull/85
sudo apt install libpcsclite-dev
make cosign-pivkey-pkcs11key



# setting CA certificate for cosign
# See go container registry issue
#    https://github.com/google/go-containerregistry/issues/211
#    https://github.com/containerd/containerd/pull/4138
export SSL_CERT_FILE=/etc/docker/certs.d/registry.127-0-10-80.nip.io/ca.crt

