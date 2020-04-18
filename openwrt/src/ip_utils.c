//
// Created by LordM on 18/04/2020.
//

#include "ip_utils.h"

#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <net/if.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <arpa/inet.h>

char *get_ip_for_interface(const char *interface) {
    int fd;
    struct ifreq ifr;

    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd == -1) return NULL;
    ifr.ifr_addr.sa_family = AF_INET;

    strncpy(ifr.ifr_name, interface, IFNAMSIZ-1);
    int error = ioctl(fd, SIOCGIFADDR, &ifr);
    close(fd);
    if (error == -1) return NULL;

    char *tmp = inet_ntoa(((struct sockaddr_in *)&ifr.ifr_addr)->sin_addr);
    char *value = malloc(strlen(tmp) * sizeof(char));
    strcpy(value, tmp);
    return value;
}