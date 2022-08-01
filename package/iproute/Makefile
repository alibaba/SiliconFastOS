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
ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
$(PKG)_VERSION = 5.12.0-4
else
$(PKG)_VERSION = 5.10.0
$(PKG)_SITE = https://cdn.kernel.org/pub/linux/utils/net/iproute2
endif

$(PKG)_SOURCE_DL_DIR = $(DOWNLOAD_DIR)/$(pkg)
$(PKG)_SOURCE_DIR = $(SOURCE_DIR)/$($(PKG)_SOURCE_NAME)

################################################################################
# set your own build cmd if you do not want use default cmd
################################################################################
# $(PKG)_DOWNLOADS_CMDS =
# $(PKG)_EXTRACT_CMDS =
# $(PKG)_CONFIGURE_CMDS =
$(PKG)_BUILD_CMDS = cd $($(PKG)_SOURCE_DIR) ;\
						make -j$(NPROCS)
$(PKG)_INSTALL_CMDS = cd $($(PKG)_SOURCE_DIR) ;\
						make DESTDIR=$($(PKG)_SOURCE_DIR)/$(BUILD_DIR) install 1> /dev/null; $($(PKG)_INSTALL_TARGET_CMDS)

include ../help.mk



