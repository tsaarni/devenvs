// NOTE: Using this NSS module can be a security risk.
//
// This is a simple NSS module that makes each user have the same UID and GID as
// the current process.
//
// It can be used in a single PID container environments to override the UID and
// GID of the real user, to match the UID and GID of the container process to
// fool applications that rely on the UID and GID of the user to determine
// permissions.

#include <nss.h>
#include <pwd.h>
#include <string.h>
#include <syslog.h>
#include <unistd.h>

enum nss_status _nss_fake_getpwnam_r(const char *name, struct passwd *pwd,
                                     char *buffer, size_t buflen, int *errnop) {
    syslog(LOG_INFO, "_nss_fake_getpwnam_r(name=%s)", name);

    char *next = buffer;

    pwd->pw_name = stpncpy(next, name, buflen);
    next = pwd->pw_name + strlen(pwd->pw_name) + 1;

    pwd->pw_dir = stpncpy(next, "/", buflen - (next - buffer));
    next = pwd->pw_dir + strlen(pwd->pw_dir) + 1;

    // Note:
    // Use /bin/bash as the shell only for demonstration purposes.
    // In a real-world scenario, the shell should be set to /sbin/nologin or
    // /bin/false to prevent shell login into single PID containers.
    pwd->pw_shell = stpncpy(next, "/bin/bash", buflen - (next - buffer));
    next = pwd->pw_shell + strlen(pwd->pw_shell) + 1;

    pwd->pw_gecos = stpncpy(next, "Fake User", buflen - (next - buffer));
    next = pwd->pw_gecos + strlen(pwd->pw_gecos) + 1;

    // Set UID and GID of current process.
    pwd->pw_uid = getuid();
    pwd->pw_gid = getgid();

    return NSS_STATUS_SUCCESS;
}
