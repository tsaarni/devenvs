/* This file implements PAM module that sets locale by calling
 * setlocale(LC_ALL, "")
 *
 * Add following to PAM stack
 *    auth required pam_setlocale.so
 */

#include <errno.h>
#include <locale.h>
#include <security/pam_ext.h>
#include <security/pam_modules.h>
#include <string.h>
#include <syslog.h>

PAM_EXTERN int pam_sm_authenticate(pam_handle_t *pamh, int /*flags*/,
								   int /*argc*/, const char ** /*argv*/) {
	char *locale = setlocale(LC_ALL, "");
	if (locale == NULL) {
		pam_syslog(pamh, LOG_ERR, "setlocale() failed: %s", strerror(errno));
	}
	pam_syslog(pamh, LOG_DEBUG, "locale set to %s", locale);
	return PAM_SUCCESS;
}
