
# Note:
# This is a test of shell logins into single UID container.
# It is only for demonstration purposes and insecure since each user will share the same permissions.


# (re)build the sssd container.
docker-compose build sssd

# Start openldap and single UID sssd + sshd container.
docker-compose up

# Test successful login.
sshpass -p user ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

# Test forced password change.
ssh mustchange@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

# Test expired password.
sshpass -p expired ssh expired@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

# Capture LDAP traffic to troubleshoot.
sudo nsenter -t $(pidof slapd) --net wireshark -i any -Y ldap -k
