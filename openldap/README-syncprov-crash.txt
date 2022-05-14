

apt-get install liblmdb-dev

autoreconf
./configure CFLAGS='-g -O0' CXXFLAGS='-g -O0' --enable-backends=mod --enable-overlays=mod --enable-modules --enable-dynamic --enable-balancer=mod --enable-argon2 --disable-wt --disable-argon2


make -j


cd tests
SLAPD_DEBUG=-1 ./run -b mdb its9776



# dump olc config
ldapsearch -LLL -x -D cn=config -wsecret -b cn=config -H ldap://localhost:9011/



ldapadd -x -D cn=Manager,dc=example,dc=com -wsecret -H ldap://localhost:9012/ <<EOF
dn: cn=user1,dc=example,dc=com
objectClass: inetOrgPerson
cn: user1
sn: user
EOF

# List users
ldapsearch -LLL -x -D cn=config -wsecret -b dc=example,dc=com -H ldap://localhost:9011/
ldapsearch -LLL -x -D cn=config -wsecret -b dc=example,dc=com -H ldap://localhost:9012/

ldapsearch -LLL -x -D cn=config -wsecret -b cn=Monitor -H ldap://localhost:9011/ +
ldapsearch -LLL -x -D cn=config -wsecret -b cn=Tasklist,cn=Threads,cn=Monitor -H ldap://localhost:9011/ monitoredInfo

ldapsearch -LLL -x -D cn=config -wsecret -b cn=Current,cn=Connections,cn=Monitor -H ldap://localhost:9011/ monitorCounter


#
# Create users
#
for i in {2..100}; do
ldapmodify -x -D cn=Manager,dc=example,dc=com -wsecret -H ldap://localhost:9011/ <<EOF
dn: cn=user${i},dc=example,dc=com
changetype: add
objectClass: inetOrgPerson
cn: user${i}
sn: user
EOF
done


#
# Toggle server1 between provider and consumer
#
while true; do

echo "Delete syncrepl consumer config"
ldapmodify -D cn=config -H ldap://localhost:9011/ -wsecret <<EOF
dn: olcDatabase={0}config,cn=config
changetype: modify
delete: olcSyncRepl
-
delete: olcMultiProvider

dn: olcDatabase={1}mdb,cn=config
changetype: modify
delete: olcSyncRepl
-
delete: olcMultiProvider
EOF
test $? != 0 && echo FAILURE && break


echo "Add syncrepl consumer config"
ldapmodify -D cn=config -H ldap://localhost:9011/ -wsecret <<EOF
dn: olcDatabase={0}config,cn=config
changetype: modify
add: olcSyncRepl
olcSyncRepl: rid=001 provider=ldap://localhost:9011/ binddn="cn=config" bindmethod=simple
  credentials=secret searchbase="cn=config" type=refreshAndPersist
  retry="5 5 300 5" timeout=1
olcSyncRepl: rid=002 provider=ldap://localhost:9012/ binddn="cn=config" bindmethod=simple
  credentials=secret searchbase="cn=config" type=refreshAndPersist
  retry="5 5 300 5" timeout=1
-
add: olcMultiProvider
olcMultiProvider: TRUE

dn: olcDatabase={1}mdb,cn=config
changetype: modify
add: olcSyncRepl
olcSyncRepl: rid=1 provider=ldap://localhost:9011/ binddn="cn=Manager,dc=example,dc=com" bindmethod=simple
  credentials=secret searchbase="dc=example,dc=com" type=refreshAndPersist
  retry="5 5 300 5" timeout=1
olcSyncRepl: rid=2 provider=ldap://localhost:9012/ binddn="cn=Manager,dc=example,dc=com" bindmethod=simple
  credentials=secret searchbase="dc=example,dc=com" type=refreshAndPersist
  retry="5 5 300 5" timeout=1
-
add: olcMultiProvider
olcMultiProvider: TRUE
EOF
test $? != 0 && echo FAILURE && break

done

