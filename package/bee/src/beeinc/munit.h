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

#ifndef MUNIT_H
#define MUNIT_H
#include <linux/cdev.h>

#define MUNIT_LOG(lev, fmt, args...) \
	printk(lev "%s:%d %s(): " fmt, __FILE__, __LINE__, __func__, ## args)

#define MLOG_INFO(fmt, args...) MUNIT_LOG(KERN_INFO, fmt, ## args)
#define MLOG_ERR(fmt, args...) MUNIT_LOG(KERN_ERR, fmt, ## args)
#define MLOG_DBG(fmt, args...) MUNIT_LOG(KERN_DEBUG, fmt, ## args)

extern void munit_get_module_case_log(const char *module_name, const char *fmt, ...);

#define FUNC_LOG(func, lev, fmt, args...)	\
do {						\
	MUNIT_LOG(lev, fmt, ## args);		\
	munit_get_module_case_log(func, fmt, ## args); \
} while(0)

#define FUNC_LOG_INFO(func, fmt, args...) FUNC_LOG(func, KERN_INFO, fmt, ##args)
#define FUNC_LOG_ERR(func, fmt, args...) FUNC_LOG(func, KERN_ERR, fmt, ##args)
#define FUNC_LOG_DBG(func, fmt, args...) FUNC_LOG(func, KERN_DEBUG, fmt, ##args)

#define LOG(lev, fmt, args...)	FUNC_LOG(__func__, lev, fmt, ##args)

#define LOG_INFO(fmt, args...)		\
	LOG(KERN_INFO, fmt, ##args)

#define LOG_ERR(fmt, args...)		\
	LOG(KERN_ERR, fmt, ##args)

#define LOG_DBG(fmt, args...)		\
	LOG_DBG(KERN_DEBUG, fmt, ##args)

#define MUNIT_BINARY_ASSERTION(left, op, right)		\
do {							\
	typeof(left) __left = (left);			\
	typeof(left) __right = (right);			\
	((void)__typecheck(__left, __right));		\
	if (__left op __right)				\
		return 0;				\
	else						\
		return -1;				\
} while(0)

#define BEE_NOT_AUTO	0

/**
 * MUNIT_EXPECT_EQ() sets an expection that @left and @right are equal
 * @left: an arbitrary expression that evaluates to a primitive C type
 * @right: an arbitrary expression that evaluates to a primitive C type.
 *
 * sets and expectation that the values @left and @right evaluate to are equal
 * This is semantically equivalent to MUNIT_EXPECT_EQ(left, right)
 */
#define MUNIT_EXPECT_EQ(left, right)	\
	MUNIT_BINARY_ASSERTION(left, ==, right)

/**
 * MUNIT_EXPECT_LT() sets an expection that @left  is less than @right
 * @left: an arbitrary expression that evaluates to a primitive C type
 * @right: an arbitrary expression that evaluates to a primitive C type.
 *
 * sets and expectation that the values @left is less than @right
 * This is semantically equivalent to MUNIT_EXPECT_EQ(left, right)
 */
#define MUNIT_EXPECT_LT(left, right)	\
	MUNIT_BINARY_ASSERTION(left, <, right)

/**
 * MUNIT_EXPECT_LE() sets an expection that @left is less than or eqaul @right
 * @left: an arbitrary expression that evaluates to a primitive C type
 * @right: an arbitrary expression that evaluates to a primitive C type.
 *
 * sets and expectation that the values @left is less than or eqaul @right
 * This is semantically equivalent to MUNIT_EXPECT_EQ(left, right)
 */
#define MUNIT_EXPECT_LE(left, right)	\
	MUNIT_BINARY_ASSERTION(left, <=, right)

/**
 * MUNIT_EXPECT_GT() sets an expection that @left is greater than @right
 * @left: an arbitrary expression that evaluates to a primitive C type
 * @right: an arbitrary expression that evaluates to a primitive C type.
 *
 * sets and expectation that the values @left is greater than @right
 * This is semantically equivalent to MUNIT_EXPECT_EQ(left, right)
 */
#define MUNIT_EXPECT_GT(left, right)	\
	MUNIT_BINARY_ASSERTION(left, >, right)

/**
 * MUNIT_EXPECT_GE() sets an expection that @left is greater than or equal @right
 * @left: an arbitrary expression that evaluates to a primitive C type
 * @right: an arbitrary expression that evaluates to a primitive C type.
 *
 * sets and expectation that the values @left and @right evaluate to are greater
 * or equal
 * This is semantically equivalent to MUNIT_EXPECT_EQ(left, right)
 */
#define MUNIT_EXPECT_GE(left, right)	\
	MUNIT_BINARY_ASSERTION(left, >=, right)

struct munit_case {
	int (*func_run)(void); /* 0 succ others fail */
	bool succ;
	char *name;
	char *log;
	char res[20];
	struct dentry *case_dir;
};

struct munit_module {
	char *name;
	int num;
	struct list_head list;
	struct munit_case *mcase;
	struct dentry *module_dir;
};

#define MUNIT_CASE(func) \
		{ .func_run = func, .succ = false, .name = #func }

extern void register_cases(char *name, struct munit_case *mcase);
extern void unregister_cases(char *name);

#define MUNIT_CASE_INIT(name, case, func) 		\
	static int case_module_init(void) 		\
	{						\
		register_cases(name, case);	\
		if (func)				\
			func();			\
		return 0;				\
	}						\
	module_init(case_module_init);			\
	static void case_module_exit(void)		\
	{						\
		unregister_cases(name);			\
	}						\
	module_exit(case_module_exit);			\
	MODULE_LICENSE("GPL");

struct munit_dev {
	dev_t devno;
	struct class *class;
	struct cdev cdev;
};

extern struct list_head case_module;
extern struct dentry *munit_debugfs_rootdir;
extern struct munit_dev mdev;
extern struct mutex case_module_mutex;

#endif
