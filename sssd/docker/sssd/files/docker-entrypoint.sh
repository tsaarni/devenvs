#!/bin/bash -ex
# Note: This script is lacking monitoring of the daemon processes.

LD_PRELOAD=/source/syslog-redirector/syslog-redirector.so SYSLOG_PATH=file:/tmp/log /usr/sbin/sssd
LD_PRELOAD=/source/syslog-redirector/syslog-redirector.so SYSLOG_PATH=file:/tmp/log /usr/sbin/sshd

tail -F /tmp/log
