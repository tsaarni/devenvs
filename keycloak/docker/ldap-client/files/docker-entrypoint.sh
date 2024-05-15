#!/bin/bash -e
chmod 600 /etc/sssd/sssd.conf
rm -f /var/run/sssd.pid
mkdir -p /run/sshd
echo "NOTE: Wait for a while for sssd to start"
systemctl init ssh sssd
