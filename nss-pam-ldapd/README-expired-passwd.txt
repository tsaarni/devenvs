https://github.com/arthurdejong/nss-pam-ldapd/issues/63


# Configure dev environment
cd ~/work/nss-pam-ldapd
mkdir -p .devcontainer .vscode
cp ~/work/devenvs/nss-pam-ldapd/configs/devcontainer.json .devcontainer
cp ~/work/devenvs/nss-pam-ldapd/configs/launch.json .vscode



# Launch vscode in devcontainer
#
# Build and install
#

./autogen.sh

# install module to /usr/lib/x86_64-linux-gnu/security/ which is the default location for pam modules in Ubuntu
./configure --prefix=/usr --with-pam-seclib-dir=/usr/lib/x86_64-linux-gnu/security/
make
sudo make install


#
# Configure
#

# create config file for nslcd
#   https://arthurdejong.org/nss-pam-ldapd/nslcd.conf.5
#

sudo bash -c 'cat > /etc/nslcd.conf <<EOF
uri ldap://openldap
base o=example
pam_authc_search NONE
EOF'



# run nslcd in foreground
sudo /usr/sbin/nslcd -d


# start another terminal and run sshd
sudo mkdir -p /run/sshd

cat > sshd_config <<EOF
KbdInteractiveAuthentication yes
UsePAM yes
EOF

while true; do sudo /usr/sbin/sshd -D -d -f sshd_config; done



# test successful login
sshpass -p joe ssh joe@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

# test password change
sshpass -p mustchange ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no



#
# Debugging
#

# Run nslcd under gdb
code --install-extension ms-vscode.cpptools
sudo chmod a+r /etc/nslcd.conf
sudo chmod a+rw /var/run/nslcd/



#### NOT WORKING for some reason
# Compile syslog-redirector to see PAM debug messages
curl -L -o /tmp/syslog-redirector.tar.gz https://github.com/Nordix/syslog-redirector/archive/refs/tags/1.0.1.tar.gz
tar -xvf /tmp/syslog-redirector.tar.gz -C /tmp
cd /tmp/syslog-redirector-1.0.1
make CFLAGS="-O2"

export LD_PRELOAD=/tmp/syslog-redirector-1.0.1/syslog-redirector.so
export SYSLOG_PATH=/dev/stdout



#########
#
# Test LDAP connectivity
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


### OpenLDAP ACL permissions: had to allow anonymous search?



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
