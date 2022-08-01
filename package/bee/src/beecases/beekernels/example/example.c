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

#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/init.h>
#include <linux/types.h>
#include <linux/err.h>
#include <linux/kallsyms.h>

#include <munit.h>

static int example_0(void)
{
	int a;
	a = 1 + 1;

	LOG_INFO("run example_0 test!\n");
	MUNIT_EXPECT_EQ(a, 2);
}

static struct munit_case test_cases[] = {
	MUNIT_CASE(example_0),
	{}
};

static int munit_init(void)
{
	LOG_INFO("kernel check setup init succ\n");
	return 0;
}

MUNIT_CASE_INIT("example", test_cases, munit_init);
MODULE_LICENSE("GPL");

