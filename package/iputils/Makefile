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

pkg = $(lastword $(subst /, ,$(CURDIR)))
PKG = $(call UPPERCASE,$(pkg))
dependence = ninja meson
include ../utils.mk
################################################################################
# user define
################################################################################
ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
$(PKG)_VERSION = 20180629-7
else
$(PKG)_VERSION = s20151218
$(PKG)_SITE = http://www.skbuff.net/iputils
$(PKG)_TAR_SOURCE = $($(PKG)_SOURCE_NAME).tar.bz2
endif

$(PKG)_SOURCE_DL_DIR = $(DOWNLOAD_DIR)/$(pkg)
$(PKG)_SOURCE_DIR = $(SOURCE_DIR)/$($(PKG)_SOURCE_NAME)

################################################################################
# set your own build cmd if you do not want use default cmd
################################################################################
# $(PKG)_DOWNLOADS_CMDS = wget -P $($(PKG)_SOURCE_DL_DIR) $($(PKG)_SITE)/$($(PKG)_VERSION).tar.gz -O $($(PKG)_SOURCE_DL_DIR)/$($(PKG)_TAR_SOURCE)
# $(PKG)_EXTRACT_CMDS =
$(PKG)_CONFIGURE_CMDS =
$(PKG)_BUILD_CMDS = cd $($(PKG)_SOURCE_DIR) ;\
						make -j$(NPROCS)
$(PKG)_INSTALL_CMDS = cd $($(PKG)_SOURCE_DIR) ; make install 1> /dev/null; mkdir -p $(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/arping 		$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/clockdiff 	$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/ipg 			$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/ping 		$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/ping6 		$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/rarpd 		$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/rdisc 		$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/tftpd 		$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/tracepath 	$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/tracepath6 	$(BUILD_DIR)/bin ; \
						mv -f $($(PKG)_SOURCE_DIR)/traceroute6 	$(BUILD_DIR)/bin ; \
						$($(PKG)_INSTALL_TARGET_CMDS)
include ../help.mk



