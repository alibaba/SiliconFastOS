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

bee_srcdir	?= $(CURDIR)

include $(bee_srcdir)/beeconfigs/bee_env.mk

COMMON_TARGET		+= beetools beecases

CLEAN_TARGET		:= $(addsuffix -clean,$(COMMON_TARGET))

INSTALL_TARGET		:= $(addsuffix -install,$(COMMON_TARGET))

.PHONY: all install $(COMMON_TARGET) $(CLEAN_TARGET) $(INSTALL_TARGET) clean

all: $(COMMON_TARGET)

clean: $(CLEAN_TARGET)

install: $(INSTALL_TARGET)

$(COMMON_TARGET):
	$(MAKE) -C "$@" \
		-f "$(bee_srcdir)/$@/Makefile" KERNEL_PATH=$(KERNEL_PATH) all

$(CLEAN_TARGET):
	$(MAKE) -C "$(subst -clean,,$@)" \
		-f "$(bee_srcdir)/$(subst -clean,,$@)/Makefile" clean

$(INSTALL_TARGET):
	@mkdir -p $(INTSALL_DIR)
	@mkdir -p $(INTSALL_DIR_KO)
	@mkdir -p $(INTSALL_DIR_BIN)
	@mkdir -p $(INTSALL_DIR_LIB)
	@cp ./bt/* $(INTSALL_DIR_BIN)
	$(MAKE) -C "$(subst -install,,$@)" \
		-f "$(bee_srcdir)/$(subst -install,,$@)/Makefile" install
