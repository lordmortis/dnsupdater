//
// Created by LordM on 18/04/2020.
//

#include "rest_utils.h"

#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <curl/curl.h>
#include <polarssl/md5.h>

#define MAX_TIMELENGTH 32

char *servicePath = "https://dnsupdate.sektorseven.net/";
char *hostnameField = "hostname";
char *ipField = "ipv4";
char *timestampField = "timestamp";
char *signatureField = "signature";

int statusCode = -1;

size_t data_write_callback(char *ptr, size_t size, size_t nmemb, void *userdata) {
    return size * nmemb;
}

size_t data_header_callback(char *buffer, size_t size, size_t nitems, void *userdata) {
    if (strncmp(buffer, "HTTP", 4) == 0) {
        int index = 0;
        char* token = strtok(buffer, " ");
        while(token) {
            if (index == 1) {
                statusCode = atoi(token);
                return nitems * size;
            }
            index++;
            token = strtok(NULL, " ");
        }
    }
    return nitems * size;
}

bool update_entry(const char *hostname, const char *secret, const char *ip) {
    bool success = false;
    int hostnameLen = strlen(hostname);
    int ipLen = strlen(ip);

    curl_global_init(CURL_GLOBAL_DEFAULT);
    CURL *curl = curl_easy_init();
    curl = curl_easy_init();
    if (!curl) {
        return false;
    }

    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, data_write_callback);
    curl_easy_setopt(curl, CURLOPT_HEADERFUNCTION, data_header_callback);

    md5_context md5Context;
    md5_starts(&md5Context);

    char *hashBytes = malloc(sizeof(char) * 16);
    char *hashChar = malloc(sizeof(char) * 33);
    char *epochString = malloc(sizeof(char) * MAX_TIMELENGTH);

    time_t tmpTime = time(NULL);
    strftime(epochString, MAX_TIMELENGTH - 1, "%s", localtime(&tmpTime));
    int epochLength = strlen(epochString);

    md5_update(&md5Context, hostname, hostnameLen * sizeof(char));
    md5_update(&md5Context, "\n", sizeof(char));
    md5_update(&md5Context, epochString, epochLength * sizeof(char));
    md5_update(&md5Context, "\n", sizeof(char));
    md5_update(&md5Context, ip, ipLen * sizeof(char));
    md5_update(&md5Context, "\n", sizeof(char));
    md5_update(&md5Context, secret, strlen(secret) * sizeof(char));
    md5_update(&md5Context, "\n", sizeof(char));

    md5_finish(&md5Context, hashBytes);
    for(int i = 0; i < 16; i++) {
        sprintf(hashChar + i * 2, "%.2x", (unsigned char)hashBytes[i]);
    }
    hashChar[32] = '\0';

    int urlLength = strlen(servicePath) + 1;
    urlLength += strlen(hostnameField) + 1 + hostnameLen;
    urlLength += 1 + strlen(ipField) + 1 + ipLen;
    urlLength += 1 + strlen(timestampField) + 1 + epochLength;
    urlLength += 1 + strlen(signatureField) + 1 +strlen(hashChar);
    urlLength += 1;

    char *url = malloc(sizeof(char) * urlLength);
    sprintf(url, "%s?%s=%s&%s=%s&%s=%s&%s=%s", servicePath, hostnameField, hostname, ipField, ip, timestampField, epochString, signatureField, hashChar);

    statusCode = -1;
    curl_easy_setopt(curl, CURLOPT_URL, url);
    CURLcode res = curl_easy_perform(curl);
    success = res == CURLE_OK;
    if (!success) {
        fprintf(stderr, "curl_easy_perform() failed: %s\n", curl_easy_strerror(res));
    }

    if (statusCode != 200) {
        success = false;
    }

    free(hashBytes);
    free(hashChar);
    free(epochString);
    curl_global_cleanup();
    md5_free(&md5Context);

    return success;
}