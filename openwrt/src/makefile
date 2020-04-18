PROFILE = -O2 -s
CFLAGS = $(PROFILE) -std=gnu99
LDFLAGS = -lcurl -luci -lpolarssl

all: main

# build it
main:
	$(CC) $(CFLAGS) main.c uci_utils.c ip_utils.c rest_utils.c $(LDFLAGS) -o s7dnsupdate

# clean it
clean:
	rm s7dnsupdate