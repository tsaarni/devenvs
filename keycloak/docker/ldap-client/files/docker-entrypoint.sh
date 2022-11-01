#!/bin/bash -ex
chmod 600 /etc/sssd/sssd.conf
rm -f /var/run/sssd.pid
mkdir -p /run/sshd
systemctl init ssh sssd
