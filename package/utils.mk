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

ARCH = aarch64
ifeq ($(V),1)
	Q =
else
	Q = @
endif

MY_CFLAGS = -Os -ffunction-sections -fdata-sections -s
MY_DFLAGS = -Wl,--gc-sections

TARGET_DIR = $(TOPDIR)/result
DOWNLOAD_DIR = $(TOPDIR)/download
PACKAGE_DIR = $(TOPDIR)/package
SOURCE_DIR = $(TARGET_DIR)/src

ROOTFS_DIR = $(TARGET_DIR)/rootfs
BEE_DIR = $(TOPDIR)/bee

ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
BUILD_DIR = $(pkg)_build/$(pkg)_rpm_build
else
BUILD_DIR = $(pkg)_build/$(pkg)_source_build
endif

PRINT = echo -e "\033[43;30m$(1) $(2)\033[0m"
UPPERCASE = $(shell echo '$(1)' | tr '[:lower:]' '[:upper:]')

# NPROCS = $(shell grep -c ^processor /proc/cpuinfo)
NPROCS = 64

ANOLIS_MIRROR = https://mirrors.openanolis.cn/anolis/8.5/BaseOS/aarch64/os/Packages/

ifeq ($(CONFIG_ANOLIS_PACKAGE_REPO),y)
$(PKG)_SITE = $(ANOLIS_MIRROR)
$(PKG)_SOURCE_NAME = $(pkg)-$($(PKG)_VERSION)
$(PKG)_RPM_SOURCE = $($(PKG)_SOURCE_NAME)*.$(ARCH).rpm
else
$(PKG)_SOURCE_NAME = $(pkg)-$($(PKG)_VERSION)
$(PKG)_TAR_SOURCE = $($(PKG)_SOURCE_NAME).tar.xz
endif

################################################################################
# make-target -- download package and build target
################################################################################
