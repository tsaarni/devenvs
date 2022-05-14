
# LDAP proxy configuration example

## Use case

LDAP server configuration:

1. Database with local users
2. LDAP backend for proxying remote LDAP server

SSSD is the LDAP client.

LDAP clients are configured with an administrative account in the local server.
The account is used by clients for operations that are not directly related to users logging in.
Clients use StartTLS and SASL EXTERNAL authentication with a client certificate and key.

LDAP proxy is configured with an administrative account in the remote server.
The account is used for proxied operations that are executed by the clients when bound with the local administrative account.
Proxy uses StartTLS when connecting to the remote server and simple bind with username and password for authentication.

## Proxy configuration

This example uses slapd OLC (online configuration).

### Global configuration section

```ldif
dn: cn=config
objectClass: olcGlobal
cn: config
olcTLSCACertificateFile: /certs/internal-client-ca.pem
olcTLSCertificateKeyFile: /certs/proxy-server-key.pem
olcTLSCertificateFile: /certs/proxy-server.pem
olcTLSVerifyClient: demand
olcSecurity: ssf=128
```

Server certificate `olcTLSCACertificateFile` and key `olcTLSCertificateKeyFile` are when serving the TLS sessions from the local clients.
The CA certificate `olcTLSCACertificateFile` is used by the server to validate client certificate.
When `olcTLSVerifyClient` is set to `demand`, the clients are required to present a client certificate when connecting to the server.
Setting `olcSecurity` to `ssf=128` will prevent client from successfully executing operations before it has upgraded the LDAP connection to TLS by issuing StartTLS.

### Local user database

The administrative account name at the local server is `cn=admin-account,dc=local,dc=com`.
Minimal `applicationProcess` object class is used for the the account, since it does not require password.
The DN of the account needs to be set in the `Subject` field of the client certificate.

```ldif
dn: olcDatabase=mdb,cn=config
objectClass: olcDatabaseConfig
objectClass: olcmdbConfig
olcDatabase:  mdb
olcSuffix: dc=local,dc=com
olcDbDirectory: /var/lib/ldap/proxy-slapd.d

dn: dc=local,dc=com
objectclass: top
objectclass: organization
objectclass: dcobject
dc: local
o: Local

dn: cn=admin-account,dc=local,dc=com
objectClass: applicationProcess
cn: admin-account
```

### Proxy configuration


The administrative account name at the remote server is `cn=admin-account,dc=remote,dc=com`.

```ldif
dn: olcDatabase=ldap,cn=config
objectClass: olcDatabaseConfig
objectClass: olcLDAPConfig
olcDatabase: LDAP
olcSuffix: dc=example,dc=com
olcDbURI: $URI1
#olcDbStartTLS: start tls_cert=$TESTDIR/certs/admin-account.pem tls_key=$TESTDIR/certs/admin-account-key.pem tls_cacert=$TESTDIR/certs/external-server-ca.pem
olcDbStartTLS: start tls_cacert=$TESTDIR/certs/internal-client-ca.pem
olcDbIdleTimeout: 3
#olcDbIDAssertBind: bindmethod=simple mode=legacy binddn="cn=admin-account,dc=example,dc=com" credentials="admin-account" tls_cert=$TESTDIR/certs/admin-account.pem tls_key=$TESTDIR/certs/admin-account-key.pem tls_cacert=$TESTDIR/certs/external-server-ca.pem tls_reqcert=demand starttls=critical flags=override
#olcDbIDAssertBind: bindmethod=simple binddn="cn=admin-account,dc=example,dc=com" credentials="admin-account" tls_cert=$TESTDIR/certs/admin-account.pem tls_key=$TESTDIR/certs/admin-account-key.pem tls_cacert=$TESTDIR/certs/external-server-ca.pem tls_reqcert=demand starttls=critical
#olcDbIDAssertBind: bindmethod=simple binddn="cn=admin-account,dc=example,dc=com" credentials="admin-account"
olcDbIDAssertBind: bindmethod=simple binddn="cn=admin-account,dc=example,dc=com" credentials="admin-account" mode=none
olcDbIDAssertAuthzFrom: dn.exact:dc=com,dc=local,cn=admin-account


olcDbIDAssertBind: bindmethod=simple binddn="cn=admin-account,dc=example,dc=com" credentials="admin-account" mode=none
olcDbIDAssertAuthzFrom: dn.exact:dc=com,dc=local,cn=admin-account



Mode none
```

## Certificate rotation


After server certificate and key is rotated, following modify operation will trigger reloading of the files without need for restarting the slapd process:

```ldif
dn: cn=config
changetype: modify
replace: olcTLSCertificateFile
olcTLSCertificateFile: /certs/proxy-server.pem
-
replace: olcTLSCertificateKeyFile
olcTLSCertificateKeyFile: /certs/proxy-server-key.pem
```

For reloading the client certificate and key for authenticating towards remote server, `delete` and `add` must be used instead of single `replace` in order to trigger reload in `slapd-ldap` backend:


```ldif
dn: olcDatabase={2}ldap,cn=config
changetype: modify
delete: olcDbStartTLS
-
add: olcDbStartTLS
olcDbStartTLS: start tls_cacert=/certs/external-server-ca.pem tls_key=$HOME/work/openldap/tests/testrun/certs/admin-account2-key.pem tls_cert=$HOME/work/openldap/tests/testrun/certs/admin-account2.pem
-
delete: olcDbIDAssertBind
-
add: olcDbIDAssertBind
olcDbIDAssertBind: bindmethod=simple mode=none binddn="cn=admin-account,dc=example,dc=com" credentials="admin-account" tls_cacert=$HOME/work/openldap/tests/testrun/certs/external-server-ca.pem tls_key=$HOME/work/openldap/tests/testrun/certs/admin-account2-key.pem tls_cert=$HOME/work/openldap/tests/testrun/certs/admin-account2.pem
```


[1] man slapd.conf - https://www.openldap.org/software/man.cgi?query=slapd.conf
[2] man slapd-ldap - https://www.openldap.org/software/man.cgi?query=slapd-ldap
[3] olc reference - https://tylersguides.com/guides/openldap-online-configuration-reference/