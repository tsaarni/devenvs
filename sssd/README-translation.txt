
# Localization does not work in SSSD PAM module
# https://github.com/SSSD/sssd/issues/7843


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

sudo mkdir -p /usr/share/locale/en_CUSTOM/LC_MESSAGES/
sudo cp en_CUSTOM.mo /usr/share/locale/en_CUSTOM/LC_MESSAGES/sssd.mo




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
    --with-smb-idmap-interface-version=6 \
make -j
sudo make install


sudo bash -c "echo 'LANGUAGE=en_CUSTOM' >> /etc/security/pam_env.conf"


sudo chmod 600 /etc/sssd/sssd.conf
sudo /usr/sbin/sssd -i

sudo mkdir -p /run/sshd
while true; do sudo /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_no; done



# Run from host
sshpass -p expired ssh expired@localhost -p 12222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
