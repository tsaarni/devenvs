#!/bin/bash -ex

# check if we are running as root
if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root"
    exit 1
fi

rm -rf /data
mkdir -p /data/ldif /data/config /data/db
chown -R $(id -u):$(id -g) /data

cp /home/tsaarni/work/devenvs/openldap/docker/openldap/files/etc/ldap/schema/* /etc/ldap/schema/

# process templates
eval "cat <<< \"$(</home/tsaarni/work/devenvs/openldap/templates/standard-server/database.ldif)\"" > /data/ldif/database.ldif
eval "cat <<< \"$(</home/tsaarni/work/devenvs/openldap/templates/standard-server/users-and-groups.ldif)\"" > /data/ldif/users-and-groups.ldif

# import processed templates
slapadd -v -n 0 -l /data/ldif/database.ldif -F /data/config
slapadd -v -n 1 -l /data/ldif/users-and-groups.ldif -F /data/config

mkdir -p /var/run/slapd

# /usr/sbin/slapd -d3 -s trace -h "ldap://0.0.0.0:389/ ldapi://%2Fvar%2Frun%2Fslapd%2Fldapi ldaps://0.0.0.0:636/" -F /data/config
./servers/slapd/slapd -d3 -s trace -h "ldap://0.0.0.0:389/ ldapi://%2Fvar%2Frun%2Fslapd%2Fldapi" -F /data/config
