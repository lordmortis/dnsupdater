#
# Copyright (C) 2006-2015 OpenWrt.org
#
# This is free software, licensed under the GNU General Public License v2.
# See /LICENSE for more information.
#

include $(TOPDIR)/rules.mk

PKG_NAME:=s7dnsupdater
PKG_VERSION:=0.0.1
PKG_RELEASE:=1
PKG_MAINTAINER:=Brendan Ragan
PKG_LICENSE:=MIT
PKG_CONFIG_DEPENDS:=libcurl libuci

include $(INCLUDE_DIR)/package.mk

PKG_BUILD_DIR := $(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)

TARGET_LDFLAGS+= \
  -Wl,-rpath-link=$(STAGING_DIR)/usr/lib \
  -Wl,-rpath-link=$(STAGING_DIR)/usr/lib/libcurl/lib
  -Wl,-rpath-link=$(STAGING_DIR)/usr/lib/libuci/lib

define Package/s7dnsupdate
  SECTION:=utils
  CATEGORY:=Utilities
  DEPENDS:=+libcurl +libuci
  TITLE:=Updates the DNS entry for the router via S7 services
  URL:=https://github.com/halmartin/load2sqlite
  MENU:=1
endef

define Package/s7dnsupdate/description
  Monitors the specified interface and attempts to update the ip if present and changed
  from the last one
endef

define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) ./src/* $(PKG_BUILD_DIR)/
endef

define Build/Configure
endef

define Build/Compile
	$(MAKE) -C $(PKG_BUILD_DIR) $(TARGET_CONFIGURE_OPTS)
endef

define Package/s7dnsupdate/install
	$(INSTALL_DIR) $(1)/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/s7dnsupdate $(1)/bin/
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(PKG_BUILD_DIR)/s7dnsupdate.conf $(1)/etc/config/s7dnsupdate
endef

$(eval $(call BuildPackage,s7dnsupdate))