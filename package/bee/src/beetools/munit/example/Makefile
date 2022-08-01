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

obj-m += example.o
example-objs := example_drv.o
MUNIT_SYM ?= Module.symvers

ccflags-y += -Wno-unused-function -Wno-unused-variable -Wno-uninitialized -Wno-maybe-uninitialized
BUILD_DIR ?= /lib/modules/$(shell uname -r)/build/

all:
	cp $(MUNIT_SYM) .
	make -C $(BUILD_DIR) M=$(PWD) modules
clean:
	make -C $(BUILD_DIR) M=$(PWD) clean
