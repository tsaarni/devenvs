package fi.protonode.listener;

import java.util.HashMap;

import org.jboss.logging.Logger;
import org.keycloak.models.KeycloakSession;
import org.keycloak.models.UserProvider;
import org.keycloak.timer.ScheduledTask;

public class FederatedUserPruner implements ScheduledTask {

    private static final Logger LOG = Logger.getLogger(FederatedUserPruner.class);

    @Override
    public void run(KeycloakSession session) {
        LOG.info("FederatedUserPruner.run() called");

        // Iterate over all realms
        session.realms().getRealmsStream().forEach(realm -> {
            LOG.infov("realm={0}", realm.getName());

            UserProvider provider = session.getProvider(UserProvider.class);

            // Iterate over all local users
            provider.searchForUserStream(realm, new HashMap<>(), null, null).forEach(user -> {

                // Check if the user is federated
                // TODO: filter by federation provider name
                String federationLink = user.getFederationLink();
                LOG.infov("  user={0} federationLink={1}", user.getUsername(), federationLink);

                // If the user is federated and has no active sessions, remove the user
                if (federationLink != null && session.sessions().getUserSessionsStream(realm, user).findAny().isEmpty()) {
                    LOG.infov("    removing user={0}", user.getUsername());
                    provider.removeUser(realm, user);
                }
            });
        });
    }

}
