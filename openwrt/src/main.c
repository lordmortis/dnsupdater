//
// Created by LordM on 16/04/2020.
//

#include "main.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <syslog.h>

#include "uci_utils.h"
#include "ip_utils.h"
#include "rest_utils.h"

struct uci_context *context = NULL;
char *secret = NULL;
char *interface = NULL;
char *hostname = NULL;
char *lastIP = NULL;
char *currentIP = NULL;

bool readConfig() {
    secret = read_uci(context, "secret");
    interface = read_uci(context, "interface");
    hostname = read_uci(context, "hostname");
    lastIP = read_uci(context, "last_ip");

    bool goodConfig = true;

    if (secret == NULL) {
        syslog(LOG_ERR, "unable to read secret from UCI config");
        goodConfig = false;
    }

    if (hostname == NULL) {
        syslog(LOG_ERR, "unable to read hostname from UCI config");
        goodConfig = false;
    }

    if (interface == NULL) {
        syslog(LOG_ERR, "unable to read interface from UCI config");
        goodConfig = false;
    }

    return goodConfig;
}

void cleanUp() {
    if (context) uci_free_context(context);
    if (secret != NULL) free(secret);
    if (hostname != NULL) free(hostname);
    if (interface != NULL) free(interface);
    if (lastIP != NULL) free(lastIP);
    if (currentIP != NULL) free(currentIP);
}


int main(void) {
    openlog("s7dnsupdate", LOG_CONS | LOG_NDELAY, LOG_USER);
    context = uci_alloc_context();

    if (!context) {
        syslog(LOG_ERR, "could not initialize UCI context");
        return -1;
    }

    if (!readConfig()) {
        syslog(LOG_ERR, "config not set correctly");
        cleanUp();
        return -1;
    }

    char *currentIP = get_ip_for_interface(interface);
    if (currentIP == NULL) {
        syslog(LOG_ERR, "could not get ip for specified interface - does it exist?");
        cleanUp();
        return -1;
    }

    if (lastIP != NULL && strcmp(lastIP, currentIP) == 0) {
        cleanUp();
        return 0;
    }

    if (update_entry(hostname, secret, currentIP)) {
        if (!write_uci(context, "last_ip", currentIP)) {
            syslog(LOG_ERR, "could not update last updated ip!");
        } else {
            syslog(LOG_INFO, "ip updated");
        }
    } else {
        syslog(LOG_ERR, "could not update ip - is hostname and secret correct?");
    }

    cleanUp();

    return 0;
}