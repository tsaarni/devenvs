

go install github.com/notaryproject/notation/cmd/notation@v1.1.0



notation cert generate-test --default "wabbit-networks.io"
# certs and keys will be generated to
#   /home/tsaarni/.config/notation/



notation key ls        # list key store
notation cert ls       # list trust store



notation sign --insecure-registry 127.0.0.1:5001/busybox@sha256:d319b0e3e1745e504544e931cde012fc5470eba649acc8a7b3607402942e5db7
