package ldapclient;

import java.util.Hashtable;

import javax.naming.Context;
import javax.naming.NamingException;
import javax.naming.ldap.InitialLdapContext;
import javax.naming.ldap.LdapContext;

public class Reconnect {
    public static void main(String[] args) throws NamingException {

        Hashtable<String, String> env = new Hashtable<>();

        env.put(Context.INITIAL_CONTEXT_FACTORY, "com.sun.jndi.ldap.LdapCtxFactory");
        env.put(Context.PROVIDER_URL, "ldap://localhost/ou=users,ou=nonexisting,o=example");

        // env.put(Context.SECURITY_AUTHENTICATION, "simple");
        // env.put(Context.SECURITY_PRINCIPAL, "cn=ldap-admin,ou=users,o=example");
        // env.put(Context.SECURITY_CREDENTIALS, "ldap-admin");

        LdapContext ldapContext = new InitialLdapContext(env, null);
        ldapContext.addToEnvironment(Context.SECURITY_AUTHENTICATION, "simple");
        ldapContext.addToEnvironment(Context.SECURITY_PRINCIPAL, "cn=ldap-admin,ou=users,o=example");
        ldapContext.addToEnvironment(Context.SECURITY_CREDENTIALS, "ldap-admin");

        ldapContext.reconnect(null);

    }
}
