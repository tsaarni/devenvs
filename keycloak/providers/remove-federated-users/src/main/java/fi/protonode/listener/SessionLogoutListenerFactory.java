package fi.protonode.listener;

import org.jboss.logging.Logger;
import org.keycloak.Config.Scope;
import org.keycloak.events.EventListenerProvider;
import org.keycloak.events.EventListenerProviderFactory;
import org.keycloak.models.KeycloakSession;
import org.keycloak.models.KeycloakSessionFactory;
import org.keycloak.timer.TimerProvider;
import org.keycloak.services.scheduled.ScheduledTaskRunner;

public class SessionLogoutListenerFactory implements EventListenerProviderFactory {

    private static final Logger LOG = Logger.getLogger(SessionLogoutListenerFactory.class);

    @Override
    public String getId() {
        LOG.info("SessionLogoutListenerFactory.getId() called");
        return "remove-user-at-logout";
    }

    @Override
    public void close() {
        LOG.info("SessionLogoutListenerFactory.close() called");
    }

    @Override
    public EventListenerProvider create(KeycloakSession session) {
        return new SessionLogoutListener(session);
    }

    @Override
    public void init(Scope scope) {
        LOG.info("SessionLogoutListenerFactory.init() called");

    }

    @Override
    public void postInit(KeycloakSessionFactory factory) {
        LOG.info("SessionLogoutListenerFactory.postInit() called");

        try (KeycloakSession session = factory.create()) {
            TimerProvider timer = session.getProvider(TimerProvider.class);
            timer.schedule(new ScheduledTaskRunner(factory, new FederatedUserPruner()),
                    1000L, "pruner");
        }

    }



}
