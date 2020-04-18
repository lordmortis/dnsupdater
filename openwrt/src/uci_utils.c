//
// Created by LordM on 18/04/2020.
//

#include "uci_utils.h"

#include <string.h>
#include <stdlib.h>

char *baseUCIPath = "s7dnsupdate.s7dnsupdate.";
char *searchPath = NULL;
int searchPathLength = -1;

void createSearchPath(const char *name) {
    int requiredLength = strlen(baseUCIPath) + strlen(name) + 1;

    if (searchPathLength < requiredLength) {
        searchPath = realloc(searchPath, sizeof(char) * requiredLength);
        searchPathLength = requiredLength;
    }

    strcpy(searchPath, baseUCIPath);
    strcpy(searchPath + strlen(baseUCIPath), name);
}

char* read_uci(struct uci_context *context, const char *name) {
    struct uci_ptr ptr;

    createSearchPath(name);

    if ((uci_lookup_ptr(context, &ptr, searchPath, true) != UCI_OK) || (ptr.o==NULL || ptr.o->v.string==NULL)) {
        return NULL;
    }

    char *value = malloc(strlen(ptr.o->v.string) * sizeof(char));
    strcpy(value, ptr.o->v.string);
    return value;
}

bool write_uci(struct uci_context *context, const char *name, const char *value) {
    struct uci_ptr ptr;

    createSearchPath(name);

    if (uci_lookup_ptr(context, &ptr, searchPath, true) != UCI_OK) {
        return false;
    }

    ptr.value = malloc(sizeof(char) * strlen(value));
    strcpy((char*)ptr.value, value);
    if (uci_set(context, &ptr) != 0) return false;
    if (uci_save(context, ptr.p) != 0) return false;
    if (uci_commit(context, &ptr.p, true) != 0) return false;
    return true;
}