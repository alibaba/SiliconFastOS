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

TOPDIR := $(CURDIR)
PACKAGE_DIR := $(TOPDIR)/package
PKG_KCONFIG_LIST := $(foreach dir, $(shell cd package/; ls -d */),$(subst /,,$(dir))_kconfig)
PKG_CLEAN_LIST := $(foreach dir, $(shell cd package/; ls -d */),$(subst /,,$(dir))_clean)
PKG_SUBDIR = $(foreach dir, $(wildcard package/*/), $(word 2,$(subst /, ,$(dir))))

SYN_KBUILD = $(TOPDIR)/result/src/kconfig/.stamp_syn_kbuild

export TOPDIR

AS		= $(CROSS_COMPILE)as
LD		= $(CROSS_COMPILE)ld
CC		= $(CROSS_COMPILE)gcc
AR		= $(CROSS_COMPILE)ar
NM		= $(CROSS_COMPILE)nm
STRIP		= $(CROSS_COMPILE)strip
OBJCOPY		= $(CROSS_COMPILE)objcopy
OBJDUMP		= $(CROSS_COMPILE)objdump
READELF		= $(CROSS_COMPILE)readelf

ifeq ($(ARCH),riscv)
QEMU = start-qemu-rv.sh
else
QEMU = start-qemu.sh
endif

export ARCH CROSS_COMPILE LD CC
export AR NM STRIP OBJCOPY OBJDUMP READELF 

all: kbuild targrt-dir $(SYN_KBUILD)

	make -C package; \
	rm -rf $(TOPDIR)/result/rootfs/share \
	 		$(TOPDIR)/result/rootfs/libexec \
			$(TOPDIR)/result/rootfs/man
	echo y | mkfs.ext4 -d result/rootfs -r 1 -N 0 -m 5 -L "rootfs" -O ^64bit result/qemu/rootfs.ext4 "200M"; \
	cp scripts/$(QEMU) result/qemu/; \
	cp scripts/bios/* result/qemu/

clean:
	rm -rf ./result
	rm -rf download/
	rm -rf include/*
	rm -rf .config*

mrproper: clean
	rm -rf include/*
	rm -rf .config*

menuconfig: kbuild
	$(TOPDIR)/result/src/kconfig/kconfig/mconf Kconfig

$(SYN_KBUILD) : .config
	$(TOPDIR)/result/src/kconfig/kconfig/conf --syncconfig Kconfig ; touch $@


kbuild: $(TOPDIR)/result/src/kconfig/.stamp_kbuild

$(TOPDIR)/result/src/kconfig/.stamp_kbuild :
	mkdir -p $(TOPDIR)/result/src/kconfig;make -C scripts/kbuild -f Makefile.sample O=$(TOPDIR)/result/src/kconfig -j64 ; touch $@

.PHONY : download
download: kbuild $(SYN_KBUILD)
	mkdir -p download/linux && touch download/linux/.stamp_downloaded-5.10
	make -C package download
	rm -rf download/linux/.stamp_downloaded-5.10

.PHONY : kconfig
kconfig: $(PKG_KCONFIG_LIST)

.PHONY : $(PKG_KCONFIG_LIST)
$(PKG_KCONFIG_LIST):
	make -C package $@

.PHONY : $(PKG_CLEAN_LIST)
$(PKG_CLEAN_LIST):
	make -C package $@

targrt-dir:
	mkdir -p $(TOPDIR)/result \
				$(TOPDIR)/result/src \
				$(TOPDIR)/result/qemu \
				$(TOPDIR)/result/rootfs \
				$(TOPDIR)/result/rootfs/boot \
				$(TOPDIR)/result/rootfs/dev/pts \
				$(TOPDIR)/result/rootfs/etc \
				$(TOPDIR)/result/rootfs/home \
				$(TOPDIR)/result/rootfs/root \
				$(TOPDIR)/result/rootfs/mnt \
				$(TOPDIR)/result/rootfs/proc \
				$(TOPDIR)/result/rootfs/sys \
				$(TOPDIR)/result/rootfs/tmp \
				$(TOPDIR)/result/rootfs/var \
				$(TOPDIR)/result/rootfs/usr \
				$(TOPDIR)/result/rootfs/usr/bin \
				$(TOPDIR)/result/rootfs/usr/sbin \
				$(TOPDIR)/result/rootfs/usr/lib \
				$(TOPDIR)/result/rootfs/usr/lib64
	cd result/rootfs ; \
	if [ ! -h lib ]; then ln -sf usr/lib lib; fi; \
	if [ ! -h bin ]; then ln -sf usr/bin bin; fi; \
	if [ ! -h sbin ]; then ln -sf usr/sbin sbin; fi; \
	if [ ! -h lib64 ]; then ln -sf usr/lib64 lib64; fi;

$(PKG_SUBDIR) :
	make -C package $@

.PHONY : vmlinux
vmlinux :
	make -C package/linux vmlinux

.PHONY : defconfig
defconfig :
	cp configs/defconfig .config

riscv64_defconfig :
	cp configs/riscv64_defconfig .config

busybox-menuconfig :
	make -C package/busybox/ busybox-menuconfig

linux-menuconfig :
	make -C package/linux/ linux-menuconfig
