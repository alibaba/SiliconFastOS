// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/module.h>

#include "../munit.h"

int add(int a, int b)
{
	return a + b;
}

int test_add_one_plus_one_0(void) /* case name: test + module name + case spec */
{
	int ret;

	ret = add(1, 1);

	LOG_INFO("this add is %d\n", ret); /* log print for case 0 */

	MUNIT_EXPECT_EQ(2, ret);
}

int test_add_one_plus_nonone_1(void) /* case name: test + module name + case spec */
{
	MUNIT_EXPECT_EQ(0, add(-1, 1));
}

struct munit_case mcase[] = {
	MUNIT_CASE(test_add_one_plus_one_0),
	MUNIT_CASE(test_add_one_plus_nonone_1),
	{}
};

void add_init(void)
{
	pr_info("init succ\n");
}

/*
 * add: your test module name
 * mcase: your test case list
 * add_init: your case ko init func
 */
MUNIT_CASE_INIT("add", mcase, add_init)
