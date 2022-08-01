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

include ../utils.mk
################################################################################
# user define
################################################################################
$(PKG)_VERSION = 1.33.1
$(PKG)_SITE = https://www.busybox.net/downloads
$(PKG)_SOURCE_NAME = $(pkg)-$(BUSYBOX_VERSION)
$(PKG)_TAR_SOURCE = $(BUSYBOX_SOURCE_NAME).tar.bz2
$(PKG)_SOURCE_DL_DIR = $(DOWNLOAD_DIR)/$(pkg)
$(PKG)_SOURCE_DIR = $(SOURCE_DIR)/$($(PKG)_SOURCE_NAME)

################################################################################
# set your own build cmd if you do not want use default cmd
################################################################################
# $(PKG)_DOWNLOADS_CMDS =
# $(PKG)_EXTRACT_CMDS =
$(PKG)_CONFIGURE_CMDS = cp -f $(pkg)_no_bash.config $($(PKG)_SOURCE_DIR)/.config
$(PKG)_BUILD_CMDS = cd $($(PKG)_SOURCE_DIR) ;\
						make -j$(NPROCS)
$(PKG)_INSTALL_CMDS = cd $($(PKG)_SOURCE_DIR) ;\
						make install 2> /dev/null;\
						cp -df $($(PKG)_SOURCE_DIR)/_install/bin/* $(ROOTFS_DIR)/bin ;\
						cp -df $($(PKG)_SOURCE_DIR)/_install/sbin/* $(ROOTFS_DIR)/sbin ;\
						cp -rdf $(PACKAGE_DIR)/$(pkg)/etc/. $(ROOTFS_DIR)/etc ;\

include ../help.mk

kconfig:

busybox-menuconfig : $($(PKG)_TARGET_EXTRACT)
	cd $($(PKG)_SOURCE_DIR); make menuconfig



