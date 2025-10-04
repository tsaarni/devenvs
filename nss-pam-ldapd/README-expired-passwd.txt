https://github.com/arthurdejong/nss-pam-ldapd/issues/63
https://github.com/arthurdejong/nss-pam-ldapd/pull/64



# SSSD forced password change
# KbdInteractiveAuthentication no
$ ssh mustchange@localhost -p 2223 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no
mustchange@localhost's password:
Password expired. Change your password now.
Password expired. Change your password now.
WARNING: Your password has expired.
You must change your password now and login again!
Current Password:
New password:
New password:
Retype new password:


# KbdInteractiveAuthentication yes

$ ssh mustchange@localhost -p 2223 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no
Warning: Permanently added '[localhost]:2223' (ED25519) to the list of known hosts.
(mustchange@localhost) Password:
(mustchange@localhost) Password expired. Change your password now.
Current Password:
(mustchange@localhost) New password:
(mustchange@localhost) Retype new password:



#
# Setup
#

# 1. configure devcontainer
cd ~/work/nss-pam-ldapd
mkdir -p .devcontainer .vscode
cp ~/work/devenvs/nss-pam-ldapd/configs/devcontainer.json .devcontainer
cp ~/work/devenvs/nss-pam-ldapd/configs/launch.json .vscode
cp ~/work/devenvs/nss-pam-ldapd/configs/extensions.json .vscode


# 2. Launch vscode. It will also automatically build and launch services from docker-compose.yml
#   - openldap
#   - sssd client (for comparison with nss-pam-ldapd)



#
# Build and install
#

./autogen.sh
./configure --prefix=/usr --with-pam-seclib-dir=/usr/lib/x86_64-linux-gnu/security/
./configure --prefix=/usr --with-pam-seclib-dir=/usr/lib/x86_64-linux-gnu/security/ --enable-nls

make
bear -- make        # create compile_commands.json for vscode while compiling

sudo make install


# Finnish locale will be installed to
#    /usr/share/locale/fi/LC_MESSAGES/nss-pam-ldapd.mo


# run nslcd in foreground

# Since nslcd calls clearenv(), setting locale with LANGUAGE does not work
# Also, just setting LC_MESSAGES does not work unless complete working locale is installed on the system
# So, install finnish locales so that LC_ALL works
sudo apt update && sudo apt install -y language-pack-fi

LC_ALL=fi_FI.UTF-8 sudo /usr/sbin/nslcd -d
LC_MESSAGES=fi_FI.UTF-8 sudo /usr/sbin/nslcd -d


# start another terminal and run sshd in foreground
sudo mkdir -p /run/sshd

# For sshd which does not call clearenv() this works:
export LC_MESSAGES=en_US.UTF-8
export LANGUAGE=fi

# kbd interactive: no
while true; do sudo /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_no; done

# kbd interactive: yes
#   Note: translations come from passwd command, not from nss-pam-ldapd
while true; do sudo /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_yes; done


# test successful login
sshpass -p joe ssh joe@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

# test password change
#  (nykyinen) LDAP-salasana:
sshpass -p mustchange ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no


# test expired password
#  Salasana vanhentunut, 4 kirjautumista jäljellä
sshpass -p expired ssh expired@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no


# to reset grace logins for expired user
docker restart nss-pam-ldapd-openldap-1



# change language by adding following to /etc/security/pam_env.conf
#   LANGUAGE=fi
#
# NOTE!!! you might need to login twice to see the change take effect

sudo chown vscode /etc/security/pam_env.conf
code /etc/security/pam_env.conf


#
# Debugging
#

# 1. Install vscode debugger extension
#    via Extensions, seach for @recommended


# 2. Hack: set some insecure permissions to allow nslcd to run as vscode user
sudo chown vscode /etc/nslcd.conf
sudo chmod a+rw /var/run/nslcd/

# 3. Launch debugger


# 4. Hack: view syslog messages from PAM under sshd
sudo rsyslogd -n


#### NOT WORKING for some reason
# Compile syslog-redirector to see PAM debug messages
curl -L -o /tmp/syslog-redirector.tar.gz https://github.com/Nordix/syslog-redirector/archive/refs/tags/1.0.1.tar.gz
tar -xvf /tmp/syslog-redirector.tar.gz -C /tmp
cd /tmp/syslog-redirector-1.0.1
make CFLAGS="-O2"

export LD_PRELOAD=/tmp/syslog-redirector-1.0.1/syslog-redirector.so
export SYSLOG_PATH=/dev/stdout




#
# Debugging LDAP stuff
#

ldapwhoami -H ldap://openldap -D cn=joe,ou=users,o=example -w joe
ldapsearch -H ldap://openldap -D cn=ldap-admin,ou=users,o=example -w ldap-admin -b ou=users,o=example



docker logs nss-pam-ldapd-openldap-1 -f        # Monitor OpenLDAP logs
docker-compose restart openldap                # Restart OpenLDAP

sudo nsenter --net -t $(pidof slapd) wireshark -k -i any -Y ldap



###########
#
# Problems
#


### OpenLDAP ACL permissions: have to allow anonymous search to allow nslcd to work?


### Login fails with error
nslcd: [495cff] <authc="mustchange"> ldap_result() failed: Insufficient access: Operations are restricted to bind/unbind/abandon/StartTLS/modify password

# fixed by
pam_authc_search NONE



### nss: queries are not cached and each lookup is a separate LDAP query
touch /tmp/foo
sudo chown 10002 /tmp/foo
ls -l /tmp/foo



### messages printed twice?
# https://github.com/linux-pam/linux-pam/issues/710



### SSH keyboard interactive authentication has (username@hostname) as prompt
###  (mustchange@localhost) Password:


### pam_get_authtok() always uses the first password prompt even if the module does not have try_first_pass
# https://github.com/linux-pam/linux-pam/issues/357
#
# maybe add parameter ignore_first_pass which avoids setting
#    pam_set_item(pamh, PAM_OLDAUTHTOK, ctx->oldpassword);
# in pam/pam.c:pam_sm_authenticate()
#
# originally the function pam_get_authtok() seems to be OpenPAM extension where it did obey "try_first_pass"
# - https://man.freebsd.org/cgi/man.cgi?query=pam_get_authtok&sektion=3&format=html
# - https://git.des.dev/OpenPAM/OpenPAM/src/commit/d61017e61587a577237436025f2d25e04393d64f/lib/libpam/pam_get_authtok.c#L107-L116
# The function has been implemented in Linux PAM since v1.1.0 which seems to be released in 2009
# but it seems "try_first_pass" was never part of Linux PAM implementation.
# Instead it behaves like "try_first_pass" was always given.


# NLS
https://www.guyrutenberg.com/2014/11/01/gettext-with-autotools-tutorial/
https://thegreyblog.blogspot.com/2017/11/making-gettext-optional-during-build-of.html
https://www.gnu.org/software/gettext/FAQ.html


# reference wrapper
/usr/share/gettext/gettext.h





### Getting no translations

"My program compiles and links fine, but doesn't output translated strings."
https://www.gnu.org/software/gettext/FAQ.html

sudo -i
apt install ltrace strace


ltrace /usr/sbin/nslcd -d 2>&1 | grep -E "(bindtextdomain|gettext|open\()"

strace -f /usr/sbin/nslcd -d 2>&1 | grep openat


export DEBUGINFOD_URLS="https://debuginfod.ubuntu.com"
gdb --args /usr/sbin/nslcd -d
set debuginfod enabled on
b bindtextdomain
b dcgettext
b open

msgunfmt /usr/share/locale/fi/LC_MESSAGES/nss-pam-ldapd.mo





cat >test.c <<EOF
#include <libintl.h>
#include <locale.h>
#include <stdio.h>

int main() {

  setlocale(LC_ALL, "");
  bindtextdomain("nss-pam-ldapd", "/usr/share/locale");
  textdomain("nss-pam-ldapd");
  printf("Using locale %s\n", setlocale(LC_ALL, NULL));

  printf(gettext("(current) LDAP Password: "));
  printf("\n");

  printf(gettext("Password expired, %d grace logins left"), 1);
  printf("\n");

  return 0;
}
EOF

gcc -o test test.c && LC_MESSAGES=en_US.UTF-8 LANGUAGE=fi ./test
gcc -o test test.c && LC_MESSAGES=en_US.UTF-8 LANGUAGE=fi strace ./test 2>&1 |grep open



