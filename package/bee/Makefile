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
$(PKG)_VERSION = 1.0
$(PKG)_SITE =
$(PKG)_SOURCE_NAME = $(pkg)-$($(PKG)_VERSION)
$(PKG)_TAR_SOURCE = $($(PKG)_SOURCE_NAME).tar.xz


$(PKG)_SOURCE_DL_DIR = $(DOWNLOAD_DIR)/$(pkg)
$(PKG)_SOURCE_DIR = $(SOURCE_DIR)/$($(PKG)_SOURCE_NAME)

$(PKG)_CFLAGS := -O2 -ffunction-sections -fdata-sections -s
$(PKG)_LDFLAGS := -Wl,--gc-sections

$(PKG)_DEPENDENCE = linux
################################################################################
# set your own build cmd if you do not want use default cmd
################################################################################
$(PKG)_DOWNLOADS_CMDS =
$(PKG)_EXTRACT_CMDS = cp -rf $(PACKAGE_DIR)/$(pkg)/src/* $($(PKG)_SOURCE_DIR)
$(PKG)_CONFIGURE_CMDS =
# $(PKG)_BUILD_CMDS =
define $(PKG)_BUILD_CMDS
	cd $($(PKG)_SOURCE_DIR); \
		./build.sh -m $(SOURCE_DIR)/lib/modules/`make -C ../linux-* --no-print-directory kernelrelease`*/build
endef
# $(PKG)_INSTALL_CMDS =
#	cd $$(@D)/$(pkg); \
#		make install

define $(PKG)_INSTALL_CMDS
	cd $($(PKG)_SOURCE_DIR); \
		./build.sh -i "$(ROOTFS_DIR)/usr/local/bee/" 1> /dev/null
endef


include ../help.mk

kconfig:



