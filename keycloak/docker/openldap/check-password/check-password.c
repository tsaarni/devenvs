#include <string.h>

#define LDAP_SUCCESS 0
#define LDAP_ERROR 1

int check_password(char *pPasswd, char **ppErrStr, void *pEntry) {
    // check that password has at least 6 characters
    if (strlen(pPasswd) < 6) {
        *ppErrStr = strdup("Password must have at least 6 characters");
        return LDAP_ERROR;
    }

    // check that password has at least one digit
    char *digits = "0123456789";
    if (!strpbrk(pPasswd, digits)) {
        *ppErrStr = strdup("Password must have at least one digit");
        return LDAP_ERROR;
    }

    return LDAP_SUCCESS;
}
