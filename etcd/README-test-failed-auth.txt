

# generate the certs
rm -rf certs
mkdir -p certs
certyaml -d certs configs/certs.yaml

# delete storage from previous run
rm -rf tmp



docker-compose -f docker-compose-releases.yaml up etcd0 etcd-client

docker exec -it etcd-etcd-client-1 ash

# configure auth
# https://etcd.io/docs/v3.5/op-guide/authentication/authentication/
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem role add root
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem user add root --no-password
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem user grant-role root root
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem auth enable

etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem role add foo
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem role grant-permission foo --prefix=true readwrite /foo/

etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem user add client --no-password
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem user grant-role client foo


# write test data
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem put /foo "foo value"
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem put /foo/data "data under foo"
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem put /bar "bar value"
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem put /bar/data "data under bar"


# test that root can read all
etcdctl --cacert=/certs/ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem get --prefix --keys-only ""

# test that client can read only /foo/*
etcdctl --cacert=/certs/ca.pem --cert=/certs/client.pem --key=/certs/client-key.pem get /foo/data
etcdctl --cacert=/certs/ca.pem --cert=/certs/client.pem --key=/certs/client-key.pem get /bar/data




# unsuccessful
etcdctl --cacert=/certs/untrusted-ca.pem --cert=/certs/root.pem --key=/certs/root-key.pem get --prefix --keys-only ""
etcdctl --cacert=/certs/ca.pem --cert=/certs/untrusted-client.pem --key=/certs/untrusted-client-key.pem get /foo/data
etcdctl --cacert=/certs/ca.pem --cert=/certs/expired-client.pem --key=/certs/expired-client-key.pem get /foo/data
etcdctl --cacert=/certs/ca.pem --cert=/certs/client.pem --key=/certs/expired-client-key.pem get /foo/data
