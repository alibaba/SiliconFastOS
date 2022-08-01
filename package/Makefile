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

TOP_DIR := ..
SUBDIRS := $(foreach dir, $(wildcard */), $(subst /, ,$(dir)))
SUBDIR_DOWNLOAD := $(foreach dir, $(SUBDIRS),$(dir)_download)
SUBDIR_CLEAN := $(foreach dir, $(SUBDIRS),$(dir)_clean)
SUBDIR_KCONFIG := $(foreach dir, $(SUBDIRS),$(dir)_kconfig)

include ../include/config/auto.conf
include utils.mk

export CONFIG_ANOLIS_PACKAGE_REPO
.PHONY : all $(SUBDIRS)
all:
	$(Q)$(foreach dir, $(SUBDIRS), \
	if [ "y" = "$(CONFIG_PACKAGE_$(call UPPERCASE,$(dir)))" ]; then\
		make -C $(dir) CONFIG_ANOLIS_PACKAGE_REPO=$(CONFIG_ANOLIS_PACKAGE_REPO);\
	fi;)


$(SUBDIRS):
	make -C $@ CONFIG_ANOLIS_PACKAGE_REPO=$(CONFIG_ANOLIS_PACKAGE_REPO);\

.PHONY : download
download : $(SUBDIR_DOWNLOAD)

.PHONY : $(SUBDIR_DOWNLOAD)

$(SUBDIR_DOWNLOAD):
	$(Q)if [ "y" = "$(CONFIG_PACKAGE_$(call UPPERCASE,$(word 1, $(subst _, ,$@))))" ]; then\
        make -C $(word 1, $(subst _, ,$@)) download CONFIG_ANOLIS_PACKAGE_REPO=$(CONFIG_ANOLIS_PACKAGE_REPO);\
    fi

.PHONY : $(SUBDIR_CLEAN)

$(SUBDIR_CLEAN):
	make -C $(word 1,$(subst _, ,$@)) clean

.PHONY : $(SUBDIR_KCONFIG)

$(SUBDIR_KCONFIG):
	make -C $(word 1,$(subst _, ,$@)) kconfig

