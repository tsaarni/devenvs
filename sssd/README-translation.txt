
# Localization does not work in SSSD PAM module
# https://github.com/SSSD/sssd/issues/7843


# NOTE!!!
#
# For setlocale() and gettext() to work and lookup additional translations from /usr/share/locale/<LANGUAGE>/LC_MESSAGES/<DOMAIN>,
# the following environment variables need to be set:
#
#   LC_MESSAGES or LC_ALL
#   LANGUAGE
#
# For example
#
#   export LC_MESSAGES=en_US.UTF-8
#   export LANGUAGE=en_CUSTOM
#
# LC_MESSAGES or LC_ALL needs to be set to a locale that exist (run `locale -a` to see available locales) e.g. en_US.UTF-8
#
# LANGUAGE will be used to lookup message catalogs before falling back to LC_MESSAGES or LC_ALL.
#
# The values will be looked up in directory set by bindtextdomain().
#
#
# From https://superuser.com/questions/392439/lang-and-language-environment-variable-in-debian-based-systems
#
#  - LANG contain the setting for all categories that are not directly set by a LC_* variable.
#  - LC_ALL is used to override every LC_* and LANG and LANGUAGE. It should not be set in a normal user environment,
#    but can be useful when you are writing a script that depend on the precise output of an internationalized command.
#  - LANGUAGE is used to set messages languages (as LC_MESSAGES) to a multi-valued value,
#    e.g., setting it to fr:de:en will use French messages where they exist; if not, it will use German messages,
#    and will fall back to English if neither German nor French messages are available.
#





# Allow many grace logins to make testing easier
docker exec -i sssd-openldap-1 ldapmodify -x -H ldap://localhost -D "cn=ldap-admin,ou=users,o=example" -w ldap-admin <<EOF
dn: cn=expires-fast,ou=ppolicy,o=example
changetype: modify
replace: pwdGraceAuthNLimit
pwdGraceAuthNLimit: 5000
EOF



# Convert the po file to mo file
cd /workspace/
msgfmt en_CUSTOM.po -o en_CUSTOM.mo


# Copy the mo file to the locales directory
sudo mkdir -p /usr/share/locale/en_CUSTOM/LC_MESSAGES/
sudo cp /workspace/en_CUSTOM.mo /usr/share/locale/en_CUSTOM/LC_MESSAGES/sssd.mo


# Compile SSSD

cd /workspace/source/sssd
autoreconf -i
./configure \
    --with-sssd-user=sssd \
    --disable-static \
    --disable-rpath \
    --prefix=/usr \
    --sysconfdir=/etc \
    --enable-pammoddir=/usr/lib/x86_64-linux-gnu/security/ \
    --enable-nsslibdir=/usr/lib/x86_64-linux-gnu/ \
    --with-systemdunitdir=/lib/systemd/system \
    --without-python2-bindings \
    --with-smb-idmap-interface-version=6
make -j
sudo make install



# Add setlocale() call at the top of sshd PAM stack.
# This is to set locale for translations in  SSSD PAM module.
# SSHD does not set the locale so we need to do it ourselves.
cd docker/sssd/files/source/pam-setlocale/
make
sudo make install
sudo sed -i '1i auth   required     pam_setlocale.so' /etc/pam.d/sshd
cat /etc/pam.d/sshd


# Run SSSD and SSHD

sudo chmod 600 /etc/sssd/sssd.conf
sudo LANGUAGE=en_CUSTOM LC_MESSAGES=en_US.UTF-8 /usr/sbin/sssd -i

sudo mkdir -p /run/sshd
while true; do sudo LANGUAGE=en_CUSTOM LC_MESSAGES=en_US.UTF-8 /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_no; done




# Run from host
sshpass -p expired ssh expired@localhost -p 12222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"




# Optional: Build own version of OpenSSH
cd /workspace/source/openssh
autoreconf -i
./configure \
    --prefix=/usr \
    --libexecdir=/usr/lib/openssh \
    --with-pid-dir=/run \
    --sysconfdir=/etc/ssh \
    --with-selinux \
    --with-security-key-builtin \
    --with-privsep-path=/run/sshd \
    --with-pam \
    --localedir=/usr/share/locale
make -j

sudo -i
export LANGUAGE=en_CUSTOM
export LC_MESSAGES=en_US.UTF-8
/workspace/source/openssh/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_no





$ git diff
diff --git a/sshd.c b/sshd.c
index 9cbe92293..b7cb33cba 100644
--- a/sshd.c
+++ b/sshd.c
@@ -127,6 +127,7 @@
 #include "sk-api.h"
 #include "srclimit.h"
 #include "dh.h"
+#include <locale.h>

 /* Re-exec fds */
 #define REEXEC_DEVCRYPTO_RESERVED_FD   (STDERR_FILENO + 1)
@@ -2159,6 +2160,8 @@ main(int ac, char **av)
        ssh_signal(SIGCHLD, SIG_DFL);
        ssh_signal(SIGINT, SIG_DFL);

+       setlocale(LC_ALL, "");
+
        /*
         * Register our connection.  This turns encryption off because we do
         * not have a key.







*** Test outside devcontainer

docker compose build sssd
docker compose up --force-recreate sssd openldap

sshpass -p expired ssh expired@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
