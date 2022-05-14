#!/bin/bash -ex

rm -rf /data/* || true
mkdir -p /data/config /data/db /data/ldif

LDIF_TEMPLATES_DIR=${LDIF_TEMPLATES_DIR:-/input/templates}
(cd $LDIF_TEMPLATES_DIR && ./init.sh)

#/usr/sbin/slapd -d3 -s trace -h "ldap://0.0.0.0:389/ ldapi://%2Fvar%2Frun%2Fslapd%2Fldapi ldaps://0.0.0.0:636/" -F /data/config
LD_PRELOAD=/libsslkeylog.so SSLKEYLOGFILE=/output/wireshark-keys.log /usr/sbin/slapd -d3 -s trace -h "ldap://0.0.0.0:389/ ldapi://%2Fvar%2Frun%2Fslapd%2Fldapi ldaps://0.0.0.0:636/" -F /data/config
