// This file implements PAM module that sets locale by calling setlocale(LC_ALL, "")
//
// Add following to PAM stack
//    auth required pam_setlocale.so

#include <security/pam_modules.h>
#include <locale.h>

PAM_EXTERN int pam_sm_authenticate(pam_handle_t * /*pamh*/, int /*flags*/, int /*argc*/,
                                   const char ** /*argv*/) {
    setlocale(LC_ALL, "");
    return PAM_SUCCESS;
}
