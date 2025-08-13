

make

bin/bao server --dev -dev-root-token-id=root

export BAO_ADDR='http://127.0.0.1:8200'
bin/bao login token

# check mounts
bin/bao read sys/mounts

bin/bao kv put secret/foo bar=baz
bin/bao kv get secret/foo

# show the curl command that does the same thing
bin/bao --
