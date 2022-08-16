package tlsserver;

import java.io.IOException;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.NoSuchAlgorithmException;
import java.security.cert.CertificateException;

public class ReloadingKeyStore extends KeyStore {

    ReloadingKeyStore(String path, char[] password) throws KeyStoreException, NoSuchAlgorithmException, CertificateException, IOException {
        super(new ReloadingKeyStoreSpi(path, password), null, "test");

        // Calling load() is necessary to set the keystore initialized state to true.
        load(null, null);
    }

}
