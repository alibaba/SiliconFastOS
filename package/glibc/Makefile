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
$(PKG)_VERSION = 2.28-164
$(PKG)_RPM_SET = glibc-common
else
$(PKG)_VERSION = 2.32
$(PKG)_SITE = https://ftp.gnu.org/gnu/glibc
endif

$(PKG)_SOURCE_DL_DIR = $(DOWNLOAD_DIR)/$(pkg)
$(PKG)_SOURCE_DIR = $(SOURCE_DIR)/$($(PKG)_SOURCE_NAME)

$(PKG)_CFLAGS := -O2 -ffunction-sections -fdata-sections -s -fPIC -w
$(PKG)_LDFLAGS := --gc-sections

GLIBC_LIBS = \
	ld*.so.* libc.so.* libdl.so.* \
	libpthread.so.* libresolv.so.* librt.so.* \
	libutil.so.* libnss_files.so.* \
	# libnss_dns.so.* libmvec.so.* libcrypt.so.* libgcc_s.so.* libm.so.*

################################################################################
# set your own build cmd if you do not want use default cmd
################################################################################
# $(PKG)_DOWNLOADS_CMDS =
# $(PKG)_EXTRACT_CMDS =
# $(PKG)_CONFIGURE_CMDS =
define $(PKG)_CONFIGURE_CMDS
	mkdir -p $(@D)/$(BUILD_DIR)
	cd $(@D)/$(BUILD_DIR); \
		$(@D)/configure \
		CFLAGS="$($(PKG)_CFLAGS)" \
		LDFLAGS="$($(PKG)_LDFLAGS)" \
		--prefix=$(@D)/$(BUILD_DIR)
		--sysconfdir=$(ROOTFS_DIR)/etc/$(pkg)
endef
# $(PKG)_BUILD_CMDS =
define $(PKG)_BUILD_CMDS
	cd $(@D)/$(BUILD_DIR); \
		make -j$(NPROCS)
endef
# $(PKG)_INSTALL_CMDS =
#	cd $$(@D)/$(pkg); \
#		make install
#	cd $$(@D)/$(pkg); \
		make install
define GLIBC_INSTALL_CMDS
	cd $(@D)/$(BUILD_DIR); \
		make install;\
		patchelf --set-interpreter /lib/ld-linux-aarch64.so.1 $(@D)/$(BUILD_DIR)/lib/libc.so.6 ;\
		$($(PKG)_INSTALL_TARGET_CMDS)

endef

include ../help.mk



