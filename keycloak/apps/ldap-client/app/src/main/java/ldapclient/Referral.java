package ldapclient;

import java.util.Hashtable;

import javax.naming.Context;
import javax.naming.NamingEnumeration;
import javax.naming.directory.DirContext;
import javax.naming.directory.InitialDirContext;

public class Referral {
    public static void main(String[] args) throws Exception {


        Hashtable<String, String> env = new Hashtable<>();

        env.put(Context.SECURITY_AUTHENTICATION, "simple");
        env.put(Context.SECURITY_PRINCIPAL,"cn=ldap-admin,ou=users,o=example");
        env.put(Context.SECURITY_CREDENTIALS, "ldap-admin");
        env.put(Context.INITIAL_CONTEXT_FACTORY, "com.sun.jndi.ldap.LdapCtxFactory");
        env.put(Context.PROVIDER_URL, "ldap://localhost/ou=users,ou=nonexisting,o=example");

        // Enable referral following on application level
        //env.put(Context.REFERRAL, "follow");

        DirContext ctx = new InitialDirContext(env);
        NamingEnumeration enm = ctx.list("");

        while (enm.hasMore()) {
            System.out.println(enm.next());
        }

        enm.close();
        ctx.close();
    }
}
