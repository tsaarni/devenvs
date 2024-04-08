// This file implements PAM module that sets the USER environment variable to
// the username of the user set in the PAM stack. It can be used in single UID
// containers to convey the username from PAM to the application.
//
// Note:
// Linux PAM has "pam_env" module that is used to set environment variables.
// However, the module does not set the environment variable until PAM handle is
// closed which is too late for some use cases. This module sets the USER
// environment variable at open_session so that it is effective immediately.

#include <syslog.h>
#include <stdlib.h>
#include <security/pam_modules.h>
#include <security/pam_ext.h>

PAM_EXTERN int pam_sm_open_session(pam_handle_t *pamh, int /*flags*/, int /*argc*/,
                                   const char ** /* argv */) {

    pam_syslog(pamh, LOG_DEBUG, "pam_sm_open_session()");

    const char *user;
    if (pam_get_item(pamh, PAM_USER, (const void **)&user) != PAM_SUCCESS) {
        pam_syslog(pamh, LOG_ERR, "pam_sm_open_session(): PAM_USER not set");
        return PAM_SESSION_ERR;
    }

    pam_syslog(pamh, LOG_DEBUG, "pam_sm_open_session(): user=%s", user);
    if (setenv("USER", user, 1) == -1) {
        pam_syslog(pamh, LOG_ERR, "pam_sm_open_session(): setenv() failed");
        return PAM_SESSION_ERR;
    }

    return PAM_SUCCESS;
}

PAM_EXTERN int pam_sm_close_session(pam_handle_t *pamh, int /* flags */, int /* argc */,
                                    const char ** /* argv */) {
    pam_syslog(pamh, LOG_DEBUG, "pam_sm_open_session()");
    return PAM_SUCCESS;
}
