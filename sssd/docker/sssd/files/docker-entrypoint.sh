#!/bin/bash -ex
# Note: This script is lacking monitoring of the daemon processes.

export LD_PRELOAD="/source/syslog-redirector/syslog-redirector.so /source/getlogin-fake/libgetlogin.so"
export SYSLOG_PATH=file:/tmp/log

# LC_ALL or LC_MESSAGES needs to be set, otherwise gettext will not attempt to look up for any translations.
export LC_MESSAGES=en_US.UTF-8

# LANGUAGE will be used to lookup message catalogs before falling back to LC_MESSAGES or LC_ALL.
# The values will be looked up in the /usr/share/locale/<LANGUAGE>/LC_MESSAGES/<DOMAIN> directory
# when the program has called bindtextdomain() with dirname /usr/share/locale.
export LANGUAGE=en_CUSTOM

/usr/sbin/sssd
/usr/sbin/sshd

tail -F /tmp/log
