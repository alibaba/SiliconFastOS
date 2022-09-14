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
$(PKG)_VERSION = 5.10
$(PKG)_SITE = https://cdn.kernel.org/pub/linux/kernel/v5.x/
$(PKG)_SOURCE_NAME = $(pkg)-$($(PKG)_VERSION)
$(PKG)_TAR_SOURCE = $($(PKG)_SOURCE_NAME).tar.xz


$(PKG)_SOURCE_DL_DIR = $(DOWNLOAD_DIR)/$(pkg)
$(PKG)_SOURCE_DIR = $(SOURCE_DIR)/$($(PKG)_SOURCE_NAME)

$(PKG)_CFLAGS := -O2 -ffunction-sections -fdata-sections -s
$(PKG)_LDFLAGS := -Wl,--gc-sections

################################################################################
# set your own build cmd if you do not want use default cmd
################################################################################
ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
$(PKG)_DOWNLOADS_CMDS := git clone -b devel-5.10 https://gitee.com/anolis/cloud-kernel.git $($(PKG)_SOURCE_DIR) || cd $($(PKG)_SOURCE_DIR) || git pull
$(PKG)_EXTRACT_CMDS =
endif
# $(PKG)_CONFIGURE_CMDS = cp -f $(PACKAGE_DIR)/$(pkg)/$(pkg).config $($(PKG)_SOURCE_DIR)/.config
$(PKG)_CONFIGURE_CMDS = cd $(@D)/; make $(MAKE_OPTION) defconfig
# $(PKG)_BUILD_CMDS =
define $(PKG)_BUILD_CMDS
	cd $(@D)/; \
		make -j$(NPROCS) $(MAKE_OPTION) 
endef
# $(PKG)_INSTALL_CMDS =
#	cd $$(@D)/$(pkg); \
#		make install

define $(PKG)_INSTALL_CMDS
	cd $(@D)/; \
		make $(MAKE_OPTION) INSTALL_MOD_PATH=$(SOURCE_DIR) modules_install; \
		cp arch/$(ARCH)/boot/Image $(TARGET_DIR)/qemu/



endef
include ../help.mk

kconfig :

linux-menuconfig : $($(PKG)_TARGET_EXTRACT)
	cd $($(PKG)_SOURCE_DIR); make $(MAKE_OPTION) menuconfig

vmlinux :
	cp -rf $(TARGET_DIR)/rootfs/ $($(PKG)_SOURCE_DIR)/usr/ && \
	sed -i 's/^CONFIG_INITRAMFS_SOURCE=\"\"$$/CONFIG_INITRAMFS_SOURCE=\"usr\/rootfs\"/g' $($(PKG)_SOURCE_DIR)/.config && \
	cd $($(PKG)_SOURCE_DIR) && \
	make -j$(NPROCS) $(MAKE_OPTION) && \
	cp  vmlinux $(TARGET_DIR)/vmlinux_include_rootfs && \
	sed -i 's/^CONFIG_INITRAMFS_SOURCE=\"usr\/rootfs\"$$/CONFIG_INITRAMFS_SOURCE=\"\"/g' $($(PKG)_SOURCE_DIR)/.config



