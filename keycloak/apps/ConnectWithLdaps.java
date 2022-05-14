// java -Djavax.net.ssl.trustStore=truststore.p12 -Djavax.net.ssl.trustStorePassword=secret ConnectWithLdaps.java

import java.util.Hashtable;
import javax.naming.*;
import javax.naming.directory.*;

public class ConnectWithLdaps {
	
    public static void main(String[] args) throws NamingException {
		
        Hashtable env = new Hashtable();
		
        // Simple bind
        env.put(Context.SECURITY_AUTHENTICATION, "simple");
        env.put(Context.SECURITY_PRINCIPAL,
                "uid=admin,ou=system");
        env.put(Context.SECURITY_CREDENTIALS, "secret");
		
        env.put(Context.INITIAL_CONTEXT_FACTORY,
            "com.sun.jndi.ldap.LdapCtxFactory");
        env.put(Context.PROVIDER_URL, "ldaps://ldaps.127-0-0-101.nip.io:443/ou=People,dc=keycloak,dc=org");
        //env.put(Context.PROVIDER_URL, "ldaps://localhost:10636/ou=People,dc=keycloak,dc=org");
		
        DirContext ctx = new InitialDirContext(env);
        NamingEnumeration enm = ctx.list("");
		
        while (enm.hasMore()) {
            System.out.println(enm.next());
        }
		
        enm.close();
        ctx.close();
    }
}

