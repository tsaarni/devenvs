
mkdir -p certs
certyaml --destination certs configs/certs.yaml   # generate certificates and keys



# 1. configure devcontainer
cd ~/work/openldap
mkdir -p .devcontainer
cp ~/work/devenvs/openldap/configs/devcontainer.json .devcontainer/devcontainer.json


# 2. launch vscode. It will automatically build and launch services from docker-compose.yml
code ~/work/openldap


# 3. inside devcontainer, build openldap
cd /workspace
autoreconf
./configure CFLAGS='-g -O0' CXXFLAGS='-g -O0' --with-argon2 --enable-backends=mod --enable-overlays=mod --enable-modules --enable-dynamic --enable-balancer=mod --enable-crypt --enable-argon2 --disable-sql --disable-wt --disable-perl --prefix=/usr --libexec=/usr/lib/ --with-subdir=ldap --sysconfdir=/etc
make depend -j
make -j

sudo make install

sudo ldconfig       # make sure installed .so's can be loaded

# 4. Run openldap

sudo ~/work/devenvs/openldap/apps/run-slapd.sh



sudo nsenter -n -t $(pidof slapd) wireshark -i any -f "port 389 or port 636" -Y ldap -k






go get -u github.com/tsaarni/certyaml             # install certyaml tool
mkdir -p certs
certyaml --destination certs configs/certs.yaml   # generate certificates and keys

docker-compose up
docker-compose rm -f  # clean previous containers




git clone git@github.com:tsaarni/ldap-test-server-notice-of-disconnect.git
cd ~/work/ldap-test-server-notice-of-disconnect
go run .



sudo apt install python3-ldap3


ipython3

from ldap3 import Server, Connection, ALL
s = Server('localhost', port=9012)


# admin
admin = Connection(s, user='cn=Manager,dc=local,dc=com', password='secret')
admin.bind()
admin.search('ou=People,dc=example,dc=com','(objectclass=*)')




# user
bjensen = Connection(s, user='cn=Barbara Jensen,ou=Information Technology Division,ou=People,dc=example,dc=com', password='bjensen')
bjensen.bind()
bjensen.search('ou=People,dc=example,dc=com','(objectclass=*)')

bjensen.search('cn=Notice of Disconnect,ou=RetCodes,dc=example,dc=com','(objectclass=*)')



# anonymous
anon = Connection(s)
anon.bind()
anon.search('ou=People,dc=example,dc=com','(objectclass=*)')






c.entries

c.search('cn=Notice of Disconnect,ou=RetCodes,dc=example,dc=com','(objectclass=*)')
c.unbind()




import ldap3
from ldap3.operation.search import search_operation
req = search_operation('ou=People,dc=example,dc=com','(objectclass=*)', ldap3.SUBTREE, ldap3.DEREF_ALWAYS, None, 0, 0, False, None, None)
bjensen.send("searchRequest", req)




ldapsearch -H ldap://localhost:9012 -D "cn=Barbara Jensen,ou=Information Technology Division,ou=People,dc=example,dc=com" -w "bjensen" -b "ou=People,dc=example,dc=com"
ldapsearch -H ldap://localhost:9012 -D "cn=Manager,dc=local,dc=com" -w "secret" -b "ou=People,dc=example,dc=com"


ldapsearch -H ldap://localhost:9012 -D "cn=Manager,dc=local,dc=com" -w "secret" -b 'cn=Monitor' '(objectClass=*)' '*' '+'


ldapsearch -H ldap://localhost:9012 -D "cn=Manager,dc=local,dc=com" -w "secret" -b "cn=Connections,cn=database 2,cn=databases,cn=monitor" '(objectClass=*)' +


ldapsearch -H ldap://localhost:9012 -D "cn=Manager,dc=local,dc=com" -w "secret" -b "dc=idle-timeout,dc=example,dc=com" -E \!sync=ro



*** Run CI tests in Docker

docker build -t openldap-ci:latest docker/openldap-ci/

docker run --rm -it --volume $PWD:/src:ro openldap-ci:latest bash
make-test





*** Development and running individual tests

# copy vscode config and test case code over the openldap repo
cp -a files/* ~/work/openldap/
cd ~/work/openldap

# compile slapd with overlays and back-ldap linked in statically
autoreconf
./configure CFLAGS='-g -O0' CXXFLAGS='-g -O0' --enable-overlays --enable-ldap

# generate compile_commands.json for vscode
intercept-build make


make -j

tail -F tests/testrun/slapd.2.log

cd ~/work/openldap/tests
killall slapd; SLAPD_DEBUG=-1 ./run -b mdb itsNNNN

gdb ~/work/openldap/tests/../servers/slapd/slapd --pid=$(pidof slapd)


# loop single test for several times
SLAPD_DEBUG=-1 ./run -l 500 -b mdb test079



# build with dynamically loadable modules
./configure CFLAGS='-g -O0' CXXFLAGS='-g -O0' --enable-backends=mod --enable-overlays=mod --enable-modules --enable-dynamic --disable-asyncmeta --disable-balancer --disable-wt
make -j




*** Decrypt LDAPS/StartTLS with wireshark

see docker/openldap/sslkeylog

wireshark -i lo -d tcp.port==9011,ldap -d tcp.port==9012,ldap -k -o tls.keylog_file:$HOME/wireshark-keys.log



*** Troubleshoot threading problems

gdb ~/work/openldap/tests/../servers/slapd/slapd --pid=$(pidof slapd)
b ldap_back_retry

b bind.c:2105
commands
silent
bt
printf "op->o_tag: 0x%x", op->o_tag
p ndn
cont
end



*** Running SSSD as LDAP client

docker build -t ldap-client:latest docker/ldap-client/
docker run --rm -it --publish=2222:22 --add-host=openldap:host-gateway --volume=$HOME/work/openldap/tests/testrun/certs:/certs:ro ldap-client
sshpass -p user ssh user@localhost -p 2222 -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no "echo Hello world!"
