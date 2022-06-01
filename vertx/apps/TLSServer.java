//
// Run by executing
//   java apps/TLSServer.java
//
// Or if using default socket w/ default keystore (SNI DOES NOT WORK WITH THIS)
//   java -Djavax.net.ssl.keyStore=/home/tsaarni/work/vert.x/src/test/resources/tls/sni-keystore.jks -Djavax.net.ssl.keyStorePassword=wibble apps/TLSServer.java
//   java -Djavax.net.ssl.keyStore=/home/tsaarni/work/devenvs/vertx/certs/server.p12 -Djavax.net.ssl.keyStorePassword=secret apps/TLSServer.java
//

import javax.net.*;
import javax.net.ssl.*;
import java.io.*;
import java.security.KeyStore;

public class TLSServer {

    //private static final String keyStorePath = "certs/server.p12";
    //private static final String keyStorePath = "certs/sni2/combined.p12";
    //private static final String keyStorePassword = "secret";

    //private static final String keyStorePath = "certs/sni-keystore.p12";
    private static final String keyStorePath = "/home/tsaarni/work/vert.x/src/test/resources/tls/sni-keystore.jks";
    private static final String keyStorePassword = "wibble";

    private static final String algorithm = "NewSunX509"; // NewSunX509 or SunX509 (-Dssl.KeyManagerFactory.algorithm=NewSunX509)
    private static int port = 8443;

    public static void main(String[] args) throws Exception {

        System.setProperty("javax.net.debug", "keymanager");

        //SSLServerSocket socket = createDefaultSocket();
        SSLServerSocket socket = createSocketWithKeyStore();
        System.out.printf("server started on port %d%n", 8443);

        while (true) {
            try (SSLSocket client = (SSLSocket) socket.accept()) {
                System.out.println("accepted");

                // Allow TLS handshake to happen by blocking the server to read().
                InputStream is = new BufferedInputStream(client.getInputStream());
                OutputStream os = new BufferedOutputStream(client.getOutputStream());
                byte[] data = new byte[2048];
                int len = is.read(data);
            }
        }
    }

    static SSLServerSocket createDefaultSocket() throws Exception {
       return (SSLServerSocket) SSLServerSocketFactory.getDefault().createServerSocket(port);
    }

    static SSLServerSocket createSocketWithKeyStore() throws Exception {
        SSLContext ctx = SSLContext.getInstance("TLS");
        KeyManagerFactory kmf = KeyManagerFactory.getInstance(algorithm);
        KeyStore ks = KeyStore.getInstance("PKCS12");
        ks.load(new FileInputStream(keyStorePath), keyStorePassword.toCharArray());
        kmf.init(ks, keyStorePassword.toCharArray());
        ctx.init(kmf.getKeyManagers(), null, null);

        SSLServerSocketFactory ssf = ctx.getServerSocketFactory();

        return (SSLServerSocket) ssf.createServerSocket(port);
    }

}
