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

DIR=$(shell pwd)
CC=gcc
CFLAGS = -I${BEE_INCLUDE}
LDLIBS = -lnuma -lpthread
SRCLIST = $(wildcard ${DIR}/*.c)
OBJLIST = $(patsubst %.c,${DIR}/%.o,$(notdir ${SRCLIST}))

.PHONY: all $(MAKE_TARGETS) clean install

all: $(MAKE_TARGETS)
	@echo "$(MAKE_TARGETS)"

$(MAKE_TARGETS): $(OBJLIST)
	$(CC) $(CFLAGS) $(LDLIBS) $< -o $@

.c.o:
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	rm -f *.o $(MAKE_TARGETS)

install:
	cp $(INSTALL_TARGETS) $(INTSALL_DIR_BIN)
