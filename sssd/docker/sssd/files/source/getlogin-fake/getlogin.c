// This file implements LD_PRELOADable version of getlogin() that returns the
// value of the USER environment variable. It can be used in single UID
// containers to convey the username from PAM to the application.

#include <syslog.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <security/pam_appl.h>

char *getlogin(void) {
    char *user = getenv("USER");

    if (user == NULL) {
        syslog(LOG_ERR, "getlogin(): USER environment variable not set");
        errno = ENOENT;
        return NULL;
    }

    syslog(LOG_DEBUG, "getlogin() = %s", user);

    return user;
}
