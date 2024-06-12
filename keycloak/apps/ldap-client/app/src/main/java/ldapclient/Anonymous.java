package ldapclient;

import java.io.IOException;
import java.util.Hashtable;

import javax.naming.Context;
import javax.naming.NamingException;
import javax.naming.ldap.InitialLdapContext;
import javax.naming.ldap.LdapContext;


public class Anonymous {

    public static void main(String[] args) throws NamingException, IOException {
        // Make anonymous LDAP connection.
        Hashtable<String, String> env = new Hashtable<String, String>();

        env.put(Context.INITIAL_CONTEXT_FACTORY, "com.sun.jndi.ldap.LdapCtxFactory");
        env.put(Context.PROVIDER_URL, "ldap://ldap.127-0-0-1.nip.io:389");
        env.put(Context.SECURITY_AUTHENTICATION, "simple");
      //  env.put(Context.SECURITY_PRINCIPAL, "");
      //  env.put(Context.SECURITY_CREDENTIALS, "");

        LdapContext ldapContext = new InitialLdapContext(env, null);

        ldapContext.reconnect(null);
        //ldapContext.lookup("");


    }
}
