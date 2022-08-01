#!/bin/bash

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

### declare your case list
declare -a case_gcc=(
case1
case2
.....
case55
)

len=${#case_gcc[@]}

###### function to run all case
gcc_test_all()
{
	local pass=0
	local fail=0
	for((i=0;i<${len};i++))
	{
		exec_case###这里执行每一个case,可以根据脚本灵活添加
		res=$?
		if [ "$res" == "0" ]; then
			pass=`expr $pass + 1`
			echo "${case_gcc[$i]}: PASS"
		else
			fail=`expr $fail + 1`
			echo "${case_gcc[$i]}: FAIL"
		fi
	}

	echo "GCC all testcases PASS: $pass, FAIL: $fail"
}

#### function to run single case
gcc_test_single_case()
{
	case_res=$1

	$case_res
	res=$?
	if [ "$res" == "0" ]; then
		echo "${case_gcc[$i]}: PASS"
	else
		echo "${case_gcc[$i]}: FAIL"
	fi
}

while true; do
	case "$1" in
		-l ) ####### list all case
			echo "${case_gcc[*]}"
			exit 0
			;;
		-a ) ######## run all case
			gcc_test_all
			exit 0
			;;
		-s ) ###### run single case
			gcc_test_single_case $2
			exit 0
			;;
	esac
done
