


mkdir -p ~/work/ldapts/dist/test-data/certs
rm ~/work/ldapts/dist/test-data/certs/*
certyaml -d ~/work/ldapts/dist/test-data/certs configs/certs.yaml



npm test

cp ~/work/devenvs/ldapts/test-starttls-external.ts .
npx ts-node test-starttls-external.ts


# run tests in specific file
npx mocha test/mytest.ts


sudo nsenter --target $(pidof slapd) --net wireshark -i lo -f "port 1389" -d tcp.port==1389,ldap -k





# attempt login and check that "Password must be changed" is returned
ldapwhoami -H ldap://localhost:389 -D cn=mustchange,ou=users,dc=example,dc=org -w mustchange -x -e ppolicy

# change password and check the expiry message

ldappasswd -H ldap://127.0.0.1:389 -D cn=expiring,ou=users,dc=example,dc=org -x -w expiring -s newpass
ldapwhoami -H ldap://127.0.0.1:389 -D cn=expiring,ou=users,dc=example,dc=org -x -e ppolicy -w newpass



docker exec ldapts-openldap-1 ldapsearch -H ldapi:/// -Y EXTERNAL -b cn=expiring,ou=users,dc=example,dc=org +




ldapsearch -H ldap://localhost -D cn=user,ou=users,dc=example,dc=org -x -w password -b "dc=example,dc=org"


# sasl external
export LDAPTLS_CERT=$HOME/work/ldapts/dist/test-data/certs/user.pem
export LDAPTLS_KEY=$HOME/work/ldapts/dist/test-data/certs/user-key.pem
export LDAPTLS_CACERT=$HOME/work/ldapts/dist/test-data/certs/server-ca.pem

export LDAPTLS_REQCERT=never
export LDAPTLS_REQCERT=demand

ldapsearch -H ldap://localhost -b "dc=example,dc=org" -ZZ -Y EXTERNAL
ldapsearch -H ldaps://localhost -b "dc=example,dc=org" -Y EXTERNAL




openssl s_client -connect localhost:389 -starttls ldap -CAfile dist/test-data/certs/server-ca.pem -cert dist/test-data/certs/admin.pem -key dist/test-data/certs/admin-key.pem
