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

TARGETS ?= $(error you must support the compile targets)

CLEAN_TARGETS := $(addsuffix -clean,$(TARGETS))

INSTALL_TARGETS := $(addsuffix -install,$(TARGETS))

.PHONY: all $(TARGETS) clean $(CLEAN_TARGETS) $(INSTALL_TARGETS)

all: $(TARGETS)

clean: $(CLEAN_TARGETS)

install: $(INSTALL_TARGETS)
	@echo $@

$(TARGETS):
	@echo $(TARGETS)
	$(MAKE) -C "$@" -f "$(CURDIR)/$@/Makefile" KERNEL_PATH=$(KERNEL_PATH) all

$(CLEAN_TARGETS):
	$(MAKE) -C "$(subst -clean,,$@)" \
	-f "$(CURDIR)/$(subst -clean,,$@)/Makefile" clean

$(INSTALL_TARGETS):
	@echo $@
	$(MAKE) -C "$(subst -install,,$@)" \
		-f "$(CURDIR)/$(subst -install,,$@)/Makefile" install
