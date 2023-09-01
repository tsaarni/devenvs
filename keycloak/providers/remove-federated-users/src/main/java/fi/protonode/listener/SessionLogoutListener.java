package fi.protonode.listener;

import org.keycloak.events.Event;
import org.keycloak.events.EventListenerProvider;
import org.keycloak.events.EventType;
import org.keycloak.events.admin.AdminEvent;
import org.keycloak.models.KeycloakSession;
import org.jboss.logging.Logger;

public class SessionLogoutListener implements EventListenerProvider {

    private static final Logger LOG = Logger.getLogger(SessionLogoutListener.class);

    private KeycloakSession session;

    public SessionLogoutListener(KeycloakSession session) {
        LOG.info("SessionLogoutListener() created");
        this.session = session;

    }

    @Override
    public void close() {
        LOG.info("SessionLogoutListener.close() called");
    }

    @Override
    public void onEvent(Event event) {
        EventType type = event.getType();

        LOG.infov("SessionLogoutListener.onEvent(event={0}) called", event.getType().toString());

        if (EventType.LOGOUT.equals(type)) {
            System.out.println("User " + event.getUserId() + " logged out");
        }


    }

    @Override
    public void onEvent(AdminEvent event, boolean includeRepresentation) {
        LOG.infov("SessionLogoutListener.onEvent(AdminEvent={0}) called", event.getOperationType().toString());
    }



}
