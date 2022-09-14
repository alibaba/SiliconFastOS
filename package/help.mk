# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public License
# as published by the Free Software Foundation; either version 2
# of the License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

dir1 := $($(PKG)_SOURCE_DIR)

ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
STAMP_SUFFIX = anolis
else
STAMP_SUFFIX = opensource
endif

$(PKG)_TARGET_INSTALL = 		$($(PKG)_SOURCE_DIR)/.stamp_installed.$(STAMP_SUFFIX)
$(PKG)_TARGET_BUILD =			$($(PKG)_SOURCE_DIR)/.stamp_built.$(STAMP_SUFFIX)
$(PKG)_TARGET_CONFIGURE =		$($(PKG)_SOURCE_DIR)/.stamp_configured.$(STAMP_SUFFIX)
$(PKG)_TARGET_EXTRACT =			$($(PKG)_SOURCE_DIR)/.stamp_extracted.$(STAMP_SUFFIX)
$(PKG)_TARGET_SOURCE =			$($(PKG)_SOURCE_DL_DIR)/.stamp_downloaded-$($(PKG)_VERSION).$(STAMP_SUFFIX)
$(PKG)_TARGET_DIRCLEAN =		$($(PKG)_SOURCE_DIR)/.stamp_dircleaned.$(STAMP_SUFFIX)

ifdef $(PKG)_DEPENDENCE
	$(PKG)_TARGET_DEPENDENCE = $(TARGET_DIR)/src/$($(PKG)_DEPENDENCE)/.stamp_installed.$(STAMP_SUFFIX)
endif

all:							$(pkg)
$(pkg):							$(pkg)-install
$(pkg)-install:					$($(PKG)_TARGET_INSTALL)
$($(PKG)_TARGET_INSTALL): 		$($(PKG)_TARGET_BUILD) $(TOPDIR)/include/generated/autoconf.h
$($(PKG)_TARGET_BUILD):			$($(PKG)_TARGET_CONFIGURE)
$($(PKG)_TARGET_CONFIGURE):		$($(PKG)_TARGET_DEPENDENCE) $($(PKG)_TARGET_EXTRACT)

$($(PKG)_TARGET_EXTRACT):		$($(PKG)_TARGET_SOURCE)
download:						$($(PKG)_TARGET_SOURCE)

ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
CLEAN_CMD = rm -rf $($(PKG)_SOURCE_DL_DIR)/.stamp_downloaded*.anolis $($(PKG)_SOURCE_DIR)/.stamp_*.anolis $($(PKG)_SOURCE_DIR)/$(BUILD_DIR)/;
else
CLEAN_CMD = make -C $($(PKG)_SOURCE_DIR) clean;\
		rm -rf $($(PKG)_SOURCE_DIR)/.stamp_built* \
			$($(PKG)_SOURCE_DIRi)/.stamp_configured* \
			$($(PKG)_SOURCE_DIRi)/.stamp_installed*;
endif
clean:
	if [ $(pkg) == bee ]; then \
		cd $($(PKG)_SOURCE_DIR); ./build.sh -c ; \
	else \
		$(CLEAN_CMD) \
	fi

$(DOWNLOAD_DIR)/%/.stamp_downloaded-$($(PKG)_VERSION).$(STAMP_SUFFIX):
	$(Q)$(call PRINT, $($(PKG)_SOURCE_NAME),"Downloading")
	$(Q)$($(PKG)_DOWNLOADS_CMDS)
	$(Q)mkdir -p $(@D)
	$(Q)touch $@

$(TARGET_DIR)/src/%/.stamp_extracted.$(STAMP_SUFFIX):
	$(Q)$(call PRINT, $($(PKG)_SOURCE_NAME),"Extracting")
	$(Q)mkdir -p $(@D)
	$($(PKG)_EXTRACT_CMDS) 2> /dev/null
	$(Q)chmod -R +rw $(@D)
	$(Q)touch $@

$(TARGET_DIR)/src/$($(PKG)_DEPENDENCE)/.stamp_installed.$(STAMP_SUFFIX):
	$(Q)$(call PRINT, $($(PKG)_SOURCE_NAME),"Dependence")
	$(Q)make -C $(PACKAGE_DIR)/$($(PKG)_DEPENDENCE)


$(TARGET_DIR)/src/%/.stamp_configured.$(STAMP_SUFFIX):
	$(Q)$(call PRINT, $($(PKG)_SOURCE_NAME),"Configuring")
	$(Q)$($(PKG)_CONFIGURE_CMDS)
	$(Q)touch $@

$(TARGET_DIR)/src/%/.stamp_built.$(STAMP_SUFFIX):
	$(Q)$(call PRINT, $($(PKG)_SOURCE_NAME),"Building")
	$($(PKG)_BUILD_CMDS)
	$(Q)touch $@

$(TARGET_DIR)/src/%/.stamp_installed.$(STAMP_SUFFIX):
	$(Q)$(call PRINT, $($(PKG)_SOURCE_NAME),"Installing")
	$(Q)$($(PKG)_INSTALL_CMDS)
	$(Q)touch $@


ifeq (,$(filter busybox dropbear linux,$(pkg)))
PKG_FILTER = y
endif
ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO)$(PKG_FILTER),yy)
$(PKG)_DOWNLOADS_CMDS = [ $(pkg) = bee ] || \
							[ -e $($(PKG)_SOURCE_DL_DIR)/$($(PKG)_RPM_SOURCE) ] || wget -P $($(PKG)_SOURCE_DL_DIR) -r -nd -np -nH --reject=html -A "$($(PKG)_RPM_SOURCE)" $(ANOLIS_MIRROR);\
							for i in $($(PKG)_RPM_SET);do \
								[ -e $($(PKG)_SOURCE_DL_DIR)/$$i ] || wget -P $($(PKG)_SOURCE_DL_DIR) -r -nd -np -nH --reject=html -A "$$i*.$(RPM_ARCH).rpm" $(ANOLIS_MIRROR);\
							done;
$(PKG)_EXTRACT_CMDS =
$(PKG)_CONFIGURE_CMDS =
$(PKG)_BUILD_CMDS =
$(PKG)_INSTALL_CMDS = mkdir -p $($(PKG)_SOURCE_DIR)/$(BUILD_DIR) ;cd $($(PKG)_SOURCE_DIR)/$(BUILD_DIR) ; \
						rpm2cpio $($(PKG)_SOURCE_DL_DIR)/$($(PKG)_RPM_SOURCE) | cpio -idmv; \
						for i in $($(PKG)_RPM_SET);do \
							rpm2cpio $($(PKG)_SOURCE_DL_DIR)/$$i*.$(RPM_ARCH).rpm | cpio -idmv; \
						done; \
						$($(PKG)_INSTALL_TARGET_CMDS)
else
$(PKG)_DOWNLOADS_CMDS ?= [ $(pkg) = bee ] || [ -e $($(PKG)_SOURCE_DL_DIR)/$($(PKG)_TAR_SOURCE) ] || wget --no-check-certificate -P $($(PKG)_SOURCE_DL_DIR) $($(PKG)_SITE)/$($(PKG)_TAR_SOURCE)
$(PKG)_EXTRACT_CMDS ?= tar -xf $($(PKG)_SOURCE_DL_DIR)/$($(PKG)_TAR_SOURCE) -C $(SOURCE_DIR)
$(PKG)_CONFIGURE_CMDS ?= \
	cd $($(PKG)_SOURCE_DIR) && rm -rf config.cache && \
	./configure \
		--prefix=$($(PKG)_SOURCE_DIR)/$(BUILD_DIR) \
		--sysconfdir=$(ROOTFS_DIR)/etc/$(pkg) \
		CFLAGS='$(MY_CFLAGS) $($(PKG)_CFLAGS)' \
		LDFLAGS='$(MY_LDFLAGS)' \
		$(AUTOMAKE_OPTION)	\
		$($(PKG)_CONF)

$(PKG)_BUILD_CMDS ?= cd $($(PKG)_SOURCE_DIR) ;\
						make -j$(NPROCS)
$(PKG)_INSTALL_CMDS ?= cd $($(PKG)_SOURCE_DIR) ;\
						make install 2> /dev/null; $($(PKG)_INSTALL_TARGET_CMDS)
endif

cmds = $(shell awk '{if (NR>6){len=split($$2,array,"_"); if(array[2] == "${PKG}" && array[3] != "BIN" && array[3] != "LIB" && array[3] != "SBIN"){str=array[3];if (len > 3)for(i = 4;i<=len;i++){str=str "." array[i];} print tolower(str)}}}' ../../include/generated/autoconf.h)

INSTALL_DIR = $($(PKG)_SOURCE_DIR)/$(BUILD_DIR)/

install_dir = bin sbin etc lib lib64 usr/bin usr/sbin usr/lib64

$(PKG)_INSTALL_TARGET_CMDS ?= 	for d in $(install_dir);do \
									if [ -d $(INSTALL_DIR)/$$d ]; then \
										for i in $(cmds);do \
											LIBPATH=$(INSTALL_DIR)/$$d/$$i ; \
											if [ -h $${LIBPATH} ]; then \
												cp -rf $${LIBPATH} $(ROOTFS_DIR)/$$d/; \
												LIBPATH=`readlink -f $${LIBPATH}`; \
												cp -rf $${LIBPATH} $(ROOTFS_DIR)/$$d/; \
											fi; \
											if [ -f $${LIBPATH} ]; then \
												cp -rf $${LIBPATH} $(ROOTFS_DIR)/$$d/; \
											fi; \
										done; \
									fi; \
								done;



# if [ -d $(INSTALL_DIR)/sbin/ ]; then \
# 	for i in $(cmds);do \
# 		if [ -f $(INSTALL_DIR)/sbin/$$i ]; then \
# 			cp -f $(INSTALL_DIR)/sbin/$$i $(ROOTFS_DIR)/sbin/; \
# 		fi; \
# 	done; \
# fi; \
# if [ -d $(INSTALL_DIR)/lib/ ]; then \
# 	for i in $(cmds);do \
# 		if [ -f $(INSTALL_DIR)/lib/$$i ]; then \
# 			cp -f $(INSTALL_DIR)/lib/$$i $(ROOTFS_DIR)/lib/; \
# 		fi; \
# 	done; \
# fi


DIRS = bin lib
KCONFIG_DIR = $(INSTALL_DIR)
bin_cmd  = $(shell [ -d $(KCONFIG_DIR)/bin ] && ls -p $(KCONFIG_DIR)/bin | grep -v /) $(shell [ -d $(KCONFIG_DIR)/usr/bin ] && ls -p $(KCONFIG_DIR)/usr/bin | grep -v /)
lib_cmd  = $(shell [ -d $(KCONFIG_DIR)/lib64 ] && ls -p $(KCONFIG_DIR)/lib64 | grep -v /) $(shell [ -d $(KCONFIG_DIR)/usr/lib64 ] && ls -p $(KCONFIG_DIR)/usr/lib64 | grep -v /)
sbin_cmd = $(shell [ -d $(KCONFIG_DIR)/sbin ] && ls -p $(KCONFIG_DIR)/sbin | grep -v /) $(shell [ -d $(KCONFIG_DIR)/usr/sbin ] && ls -p $(KCONFIG_DIR)/usr/sbin | grep -v /)
kconfig:
	$(Q)echo "menuconfig PACKAGE_$(PKG)" > Kconfig
	$(Q)echo "	bool \"$(pkg)\"" >> Kconfig

	$(Q)echo "if PACKAGE_$(PKG)" >> Kconfig

	$(Q)echo "menuconfig $(PKG)_BIN" >> Kconfig
	$(Q)echo "	bool \"$(pkg)_bin\"" >> Kconfig

	$(Q)echo "if $(PKG)_BIN" >> Kconfig
	$(Q)$(foreach cmd,$(bin_cmd),echo "config $(PKG)_$(call UPPERCASE,$(subst .,_,$(cmd)))" >> Kconfig; echo "	bool \"$(subst .,_,$(cmd))\"" >> Kconfig;)
	$(Q)echo "endif" >> Kconfig

	$(Q)echo "menuconfig $(PKG)_SBIN" >> Kconfig
	$(Q)echo "	bool \"$(pkg)_sbin\"" >> Kconfig

	$(Q)echo "if $(PKG)_SBIN" >> Kconfig
	$(Q)$(foreach cmd,$(sbin_cmd),echo "config $(PKG)_$(call UPPERCASE,$(subst .,_,$(cmd)))" >> Kconfig; echo "	bool \"$(subst .,_,$(cmd))\"" >> Kconfig;)
	$(Q)echo "endif" >> Kconfig

	$(Q)echo "menuconfig $(PKG)_LIB" >> Kconfig
	$(Q)echo "	bool \"$(pkg)_lib\"" >> Kconfig

	$(Q)echo "if $(PKG)_LIB" >> Kconfig

	$(Q)$(foreach cmd,$(lib_cmd),echo "config $(PKG)_$(call UPPERCASE,$(subst .,_,$(cmd)))" >> Kconfig; echo "	bool \"$(subst .,_,$(cmd))\"" >> Kconfig;)
	$(Q)echo "endif" >> Kconfig

	$(Q)echo "endif" >> Kconfig
