
https://github.com/arthurdejong/nss-pam-ldapd/pull/65





clear_first_pass





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
make
bear -- make        # create compile_commands.json for vscode while compiling
sudo make install




# run nslcd in foreground
sudo /usr/sbin/nslcd -d




# start another terminal and run sshd in foreground
sudo mkdir -p /run/sshd

# kbd interactive: yes
while true; do sudo /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_yes; done

# test password change
ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no




# kbd interactive: no
# Note: THIS WILL NOT REQUIRE CLEAR_FIRST_PASS
# because it executes passwd command as separate process and therefore does not have access
# to the password entered by the user which is in memory of the sshd process
while true; do sudo /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_no; done
