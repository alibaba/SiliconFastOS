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

MAKE := make
CC := gcc

KERNEL_PATH ?= /lib/modules/$(shell uname -r)/build/

bee_srcdir ?= $(error you must support the bee top directory)

abs_srcdir ?=  $(abspath $(bee_srcdir))

BEE_INCLUDE ?= $(abs_srcdir)/beeinc

INSTALL_PATH ?= $(abs_srcdir)
INTSALL_DIR ?= $(INSTALL_PATH)/build
INTSALL_DIR_KO ?= $(INSTALL_PATH)/build/ko
INTSALL_DIR_BIN ?= $(INSTALL_PATH)/build/bin
INTSALL_DIR_LIB ?= $(INSTALL_PATH)/build/lib64
