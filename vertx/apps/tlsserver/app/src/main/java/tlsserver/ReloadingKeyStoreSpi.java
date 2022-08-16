package tlsserver;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.nio.file.attribute.FileTime;
import java.security.Key;
import java.security.KeyStore;
import java.security.KeyStoreException;
import java.security.KeyStoreSpi;
import java.security.NoSuchAlgorithmException;
import java.security.UnrecoverableKeyException;
import java.security.cert.Certificate;
import java.security.cert.CertificateException;
import java.time.Duration;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Date;
import java.util.Enumeration;
import java.util.logging.Logger;


public class ReloadingKeyStoreSpi extends KeyStoreSpi {

    private static Logger log = Logger.getLogger(ReloadingKeyStoreSpi.class.toString());

    // Defines how often the keystore file should be checked for changes
    final Duration cacheTtl = Duration.of(1, ChronoUnit.SECONDS);

    private String path;
    private char[] password;
    private FileTime lastModified;
    private Instant cacheExpiredTime = Instant.MIN;
    private KeyStore delegate;

    ReloadingKeyStoreSpi(String path, char[] password) throws KeyStoreException, NoSuchAlgorithmException, CertificateException, IOException {
        this.path = path;
        this.password = password;
        delegate = KeyStore.getInstance("PKCS12");

        // Load first time to catch invalid path early.
        refresh();
    }

    private void refresh() throws NoSuchAlgorithmException, CertificateException, IOException {
        // If keystore has been previously loaded, check the modification timestamp to decide if reload is needed.
        if ((lastModified != null) && (lastModified.compareTo(Files.getLastModifiedTime(Paths.get(path))) > 0)) {
            // File was not modified since last reload.
            return;
        }

        log.info("Loading keystore from disk");

        delegate.load(new FileInputStream(path), password);
        this.lastModified = Files.getLastModifiedTime(Paths.get(path));
    }


    private void refreshNoThrow() {
        // Has enough time passed for the keystore to be refreshed?
        if (Instant.now().isBefore(cacheExpiredTime)) {
            return;
        }

        // Set the next time when refresh should be checked for possible update.
        cacheExpiredTime = Instant.now().plus(cacheTtl);

        try {
            refresh();
        } catch (Exception e) {
            log.info("Failed to refresh: " + e);
        }
    }


    @Override
    public Key engineGetKey(String alias, char[] password) throws NoSuchAlgorithmException, UnrecoverableKeyException {
        log.info("engineGetKey");
        refreshNoThrow();
        try {
            return delegate.getKey(alias, password);
        } catch (KeyStoreException e) {
            return null;
        }
    }

    @Override
    public Certificate[] engineGetCertificateChain(String alias) {
        log.info("engineGetCertificateChain");
        refreshNoThrow();
        try {
            return delegate.getCertificateChain(alias);
        } catch (KeyStoreException e) {
            return new Certificate[0];
        }
    }

    @Override
    public Certificate engineGetCertificate(String alias) {
        log.info("engineGetCertificate");
        refreshNoThrow();
        try {
            return delegate.getCertificate(alias);
        } catch (KeyStoreException e) {
            return null;
        }
    }

    @Override
    public Date engineGetCreationDate(String alias) {
        log.info("engineGetCreationDate");
        refreshNoThrow();
        try {
            return delegate.getCreationDate(alias);
        } catch (KeyStoreException e) {
            return null;
        }
    }

    @Override
    public void engineSetKeyEntry(String alias, Key key, char[] password, Certificate[] chain)
            throws KeyStoreException {
        log.info("engineSetKeyEntry");
        throw new UnsupportedOperationException();
    }

    @Override
    public void engineSetKeyEntry(String alias, byte[] key, Certificate[] chain) throws KeyStoreException {
        log.info("engineSetKeyEntry");
        throw new UnsupportedOperationException();
    }

    @Override
    public void engineSetCertificateEntry(String alias, Certificate cert) throws KeyStoreException {
        log.info("engineSetCertificateEntry");
        throw new UnsupportedOperationException();
    }

    @Override
    public void engineDeleteEntry(String alias) throws KeyStoreException {
        log.info("engineDeleteEntry");
        throw new UnsupportedOperationException();
    }

    @Override
    public Enumeration<String> engineAliases() {
        log.info("engineAliases");
        refreshNoThrow();
        try {
            return delegate.aliases();
        } catch (KeyStoreException e) {
            return null;
        }
    }

    @Override
    public boolean engineContainsAlias(String alias) {
        log.info("engineContainsAlias");
        refreshNoThrow();
        try {
            return delegate.containsAlias(alias);
        } catch (KeyStoreException e) {
            return false;
        }
    }

    @Override
    public int engineSize() {
        log.info("engineSize");
        refreshNoThrow();
        try {
            return delegate.size();
        } catch (KeyStoreException e) {
            return 0;
        }
    }

    @Override
    public boolean engineIsKeyEntry(String alias) {
        log.info("engineIsKeyEntry");
        refreshNoThrow();
        try {
            return delegate.isKeyEntry(alias);
        } catch (KeyStoreException e) {
            return false;
        }
    }

    @Override
    public boolean engineIsCertificateEntry(String alias) {
        log.info("engineIsCertificateEntry");
        refreshNoThrow();
        try {
            return delegate.isCertificateEntry(alias);
        } catch (KeyStoreException e) {
            return false;
        }
    }

    @Override
    public String engineGetCertificateAlias(Certificate cert) {
        log.info("engineGetCertificateAlias");
        refreshNoThrow();
        try {
            return delegate.getCertificateAlias(cert);
        } catch (KeyStoreException e) {
            return null;
        }
    }

    @Override
    public void engineStore(OutputStream stream, char[] password)
            throws IOException, NoSuchAlgorithmException, CertificateException {
        log.info("engineStore");
        throw new UnsupportedOperationException();
    }

    @Override
    public void engineLoad(InputStream stream, char[] password)
            throws IOException, NoSuchAlgorithmException, CertificateException {
        log.info("engineLoad");
    }
}
