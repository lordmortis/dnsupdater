//
// Created by LordM on 18/04/2020.
//

#ifndef OPENWRTUPDATER_UCI_UTILS_H
#define OPENWRTUPDATER_UCI_UTILS_H

#include <uci.h>

char *read_uci(struct uci_context *context, const char *name);
bool write_uci(struct uci_context *context, const char *name, const char *value);

#endif //OPENWRTUPDATER_UCI_UTILS_H
