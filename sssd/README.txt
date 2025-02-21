
Open the devcontainer in this directory in vscode


Following mounts are available

/workspace                  # ~/work/devenvs/sssd
/workspace/source/sssd      # ~/work/sssd
/workspace/source/openssh   # ~/work/openssh-portable


Following ports are exposed

2222                       # sshd in sssd container
12222                      # sshd in vscode container (devcontainer)



# Compile sssd in devcontainer

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



# Run sssd

sudo chmod 600 /etc/sssd/sssd.conf
sudo /usr/sbin/sssd -i


# Run sshd

sudo mkdir -p /run/sshd
while true; do sudo /usr/sbin/sshd -D -d -f /etc/ssh/sshd_config_kdb_interactive_no; done



# Login from host to vscode container

sshpass -p user ssh user@localhost -p 12222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"


# Login from vscode container to sssd container

sshpass -p expired ssh expired@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
