package ldapclient;

public class App {

    public static void main(String[] args) throws Exception {

        // Parse subcommand arguments
        if (args.length < 1) {
            System.out.println("Usage: ldap-client <subcommand> [args]");
            System.exit(1);
        }

        System.out.println("Java version: " + System.getProperty("java.version") + "\n");

        String[] subargs = java.util.Arrays.copyOfRange(args, 1, args.length);

        switch (args[0]) {
            case "referral":
                Referral.main(subargs);
                break;

            case "starttls":
                StartTLS.main(subargs);
                break;

            case "reconnect":
                Reconnect.main(subargs);
                break;

            default:
                System.out.println("Unknown subcommand: " + args[0]);
                System.exit(1);
        }

    }
}
