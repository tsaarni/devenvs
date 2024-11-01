
### This changed the behavior of the expired password login.

# [RHEL8] sssd attempts LDAP password modify extended op after BIND failure #6768
# https://github.com/SSSD/sssd/issues/6768
# https://github.com/SSSD/sssd/pull/6769
# https://bugzilla.redhat.com/show_bug.cgi?id=2139467



# https://datatracker.ietf.org/doc/html/draft-behera-ldap-password-policy-11

# bindResponse.resultCode = invalidCredentials (49)
# passwordPolicyResponse.error = passwordExpired (0):


# The code was later changed to support sending LDAP Password Modify extended operation
# which does not require a successful BIND, at least with OpenLDAP
# https://github.com/SSSD/sssd/commit/7184541976608d357a5da48d09a7fa08862477d8#diff-22c1ff608c5207b4d1d188b248a5a788c2da106eaa743f58b6f113f5abafa08cR239


# in sssd.conf LDAP domain section, add:
ldap_pwmodify_mode = exop_force

# compile
docker compose build sssd

# Run the container
docker compose up --force-recreate

# Test successful login.
sshpass -p user ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"

# Test password that has already expired.
sshpass -p expired ssh expired@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
