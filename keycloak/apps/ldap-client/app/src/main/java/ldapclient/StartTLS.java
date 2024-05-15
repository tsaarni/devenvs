package ldapclient;

import java.io.IOException;
import java.util.Hashtable;

import javax.naming.Context;
import javax.naming.NamingException;
import javax.naming.ldap.InitialLdapContext;
import javax.naming.ldap.LdapContext;
import javax.naming.ldap.StartTlsRequest;
import javax.naming.ldap.StartTlsResponse;
import javax.net.ssl.SSLSession;

public class StartTLS {
    public static void main(String[] args) throws NamingException, IOException {

        Hashtable<String, Object> env = new Hashtable<String, Object>();

        // Connection properties
        env.put(Context.INITIAL_CONTEXT_FACTORY, "com.sun.jndi.ldap.LdapCtxFactory");
        env.put(Context.PROVIDER_URL, "ldap://ldap.127-0-0-1.nip.io");

        LdapContext ldapContext = new InitialLdapContext(env, null);
        StartTlsResponse tls = (StartTlsResponse) ldapContext.extendedOperation(new StartTlsRequest());
        SSLSession session = tls.negotiate();

        // Simple bind
        ldapContext.addToEnvironment(Context.SECURITY_AUTHENTICATION, "simple");
        ldapContext.addToEnvironment(Context.SECURITY_PRINCIPAL, "cn=ldap-admin,ou=users,o=example");
        ldapContext.addToEnvironment(Context.SECURITY_CREDENTIALS, "ldap-admin");

        ldapContext.lookup("");

        ldapContext.reconnect(null);
    }
}
