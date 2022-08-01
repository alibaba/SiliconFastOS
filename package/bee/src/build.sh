#!/bin/sh

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

export bee_srcdir=`pwd`

COMMON_MODULE_SYMVERS="$bee_srcdir/beeconfigs/Module.symvers"
TOOLS_DIR="$bee_srcdir/beetools"
# MAKE_PROCESS="-j$(getconf _NPROCESSORS_ONLN)"
MAKE_PROCESS="-j64"
if [ -z"$2" ]
  then
    echo '1';KERNEL_PATH="KERNEL_PATH=$2"
fi


function bee_make()
{
	echo $MAKE_PROCESS
	make $MAKE_PROCESS beetools $KERNEL_PATH 2> /dev/null

	rm -rf $COMMON_MODULE_SYMVERS
	touch "$COMMON_MODULE_SYMVERS"

	for file in `find $TOOLS_DIR -name Module.symvers`
	do
		cat $file >> $COMMON_MODULE_SYMVERS
	done

	make $MAKE_PROCESS beecases $KERNEL_PATH 2> /dev/null
}

function bee_install()
{
	make install
	cp ./beeconfigs/ko_config $INSTALL_PATH/build/ko
	cp ./lib64/* $INSTALL_PATH/build/lib64
}

function bee_clean()
{
	make clean
	rm -rf $COMMON_MODULE_SYMVERS
	rm -rf ./build
}

function bee_usage()
{
	cat << EOF
Usage:
$0 [-m] [-i INSTALL_PATH] [-c]
$0 -h

Options:
-h	Print this help
-m	Compiler thr bee test include all beetools and beecases
-i	After compile, you can instal all cases and scripts
	eg1: no arg represent the default install path; eg2:/bee_install: where you want to install the bin
-c	clean all the files after compile
EOF
}

while getopts "hmci" opt; do
case $opt in
	h )
		bee_usage
		exit 0
		;;
	m )
		bee_make
		exit 0
		;;
	i )
		INSTALL_PATH=$bee_srcdir
		if [ -n "$2" ]; then
			INSTALL_PATH="$2"
		fi
		export INSTALL_PATH
		bee_install
		exit 0
		;;
	c )
		bee_clean
		exit 0
		;;
	* )
		echo "error cmd!"
		exit 1
		;;
esac
done
