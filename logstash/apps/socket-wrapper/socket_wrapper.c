#include <stdio.h>
#include <dlfcn.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <errno.h>

int dscp = 46; // EF explicit forwarding.

extern int __bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
extern int __connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen);

int bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen) {
    if (addr->sa_family == AF_INET) {
        printf("Setting DSCP for IPv4 socket during bind\n");
        int tos = dscp << 2;
        if (setsockopt(sockfd, IPPROTO_IP, IP_TOS, &tos, sizeof(tos)) == -1) {
            fprintf(stderr, "Failed to set DSCP for IPv4: %s\n", strerror(errno));
        }
    } else if (addr->sa_family == AF_INET6) {
        printf("Setting DSCP for IPv6 socket during bind\n");
        int tclass = dscp << 2;
        if (setsockopt(sockfd, IPPROTO_IPV6, IPV6_TCLASS, &tclass, sizeof(tclass)) == -1) {
            fprintf(stderr, "Failed to set DSCP for IPv6: %s\n", strerror(errno));
        }

        // Check if IP_TOS can be set on IPv6 sockets
        int tos = dscp << 2;
        if (setsockopt(sockfd, IPPROTO_IP, IP_TOS, &tos, sizeof(tos)) == -1) {
            if (errno == ENOPROTOOPT) {
                fprintf(stderr, "Setting IP_TOS on IPv6 socket is not supported\n");
            } else {
                fprintf(stderr, "Failed to set DSCP for IPv6 (IP_TOS): %s\n", strerror(errno));
            }
        }
    }
    return __bind(sockfd, addr, addrlen);
}

int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen) {
    if (addr->sa_family == AF_INET) {
        printf("Setting DSCP for IPv4 socket during connect\n");
        int tos = dscp << 2;
        if (setsockopt(sockfd, IPPROTO_IP, IP_TOS, &tos, sizeof(tos)) == -1) {
            fprintf(stderr, "Failed to set DSCP for IPv4: %s\n", strerror(errno));
        }
    } else if (addr->sa_family == AF_INET6) {
        int tclass = dscp << 2;
        if (setsockopt(sockfd, IPPROTO_IPV6, IPV6_TCLASS, &tclass, sizeof(tclass)) == -1) {
            fprintf(stderr, "Failed to set DSCP for IPv6: %s\n", strerror(errno));
        }

        // Check if IP_TOS can be set on IPv6 sockets
        int tos = dscp << 2;
        if (setsockopt(sockfd, IPPROTO_IP, IP_TOS, &tos, sizeof(tos)) == -1) {
            if (errno == ENOPROTOOPT) {
                fprintf(stderr, "Setting IP_TOS on IPv6 socket is not supported\n");
            } else {
                fprintf(stderr, "Failed to set DSCP for IPv6 (IP_TOS): %s\n", strerror(errno));
            }
        }
    }

    return __connect(sockfd, addr, addrlen);
}
