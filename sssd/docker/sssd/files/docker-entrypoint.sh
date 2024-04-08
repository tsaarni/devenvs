#!/bin/bash -ex
# Note: This script is lacking monitoring of the daemon processes.

export LD_PRELOAD="/source/syslog-redirector/syslog-redirector.so /source/getlogin-fake/libgetlogin.so"
export SYSLOG_PATH=file:/tmp/log

/usr/sbin/sssd
/usr/sbin/sshd

tail -F /tmp/log
