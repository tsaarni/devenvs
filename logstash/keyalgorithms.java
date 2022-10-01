import java.security.Security;
import org.bouncycastle.jce.provider.BouncyCastleProvider;

class KeyAlgorithms {
    public static void main(String[] args) {
        System.out.println("JDK KeyFactory algorithms: " + Security.getAlgorithms("KeyFactory"));

        Security.addProvider(new BouncyCastleProvider());

        System.out.println("BC KeyFactory algorithms: " + Security.getAlgorithms("KeyFactory"));


    }
}
