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

#include <linux/err.h>
#include <linux/slab.h>
#include <linux/cdev.h>
#include <linux/errno.h>
#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/debugfs.h>
#include <linux/uaccess.h>

#include <munit.h>

#define DEVICE_NAME "munit"

struct dentry *munit_debugfs_rootdir;
EXPORT_SYMBOL_GPL(munit_debugfs_rootdir);

struct munit_dev mdev;
EXPORT_SYMBOL_GPL(mdev);

struct mutex case_module_mutex = __MUTEX_INITIALIZER(case_module_mutex);
EXPORT_SYMBOL_GPL(case_module_mutex);

struct list_head case_module = LIST_HEAD_INIT(case_module);
EXPORT_SYMBOL_GPL(case_module);
extern struct file_operations munit_fops;
extern struct file_operations munit_res_fops;

static struct munit_case *munit_get_module_case(const char *string, struct munit_module *module)
{
	int i = 0;
	struct munit_case *mcase = module->mcase;

	for (i = 0; i < module->num; i++)
	{
		if (!strcmp(string, mcase[i].name))
			return &mcase[i];
	}

	return NULL;
}

static struct munit_module *munit_search_test_module(const char *module_name)
{
	struct list_head *itr;
	struct munit_module *module;

	mutex_lock(&case_module_mutex);
	if (!list_empty(&case_module)) {
		list_for_each(itr, &case_module) {
			module = list_entry(itr, struct munit_module, list);
			if (strstr(module_name, module->name)) {
				mutex_unlock(&case_module_mutex);
				return module;
			}
		}
	}
	mutex_unlock(&case_module_mutex);

	return NULL;
}

static struct munit_module *munit_get_case_module_struct(const char *string)
{
	return munit_search_test_module(string);
}

void munit_get_module_case_log(const char *module_name, const char *fmt, ...)
{
	int case_id, len;
	va_list args;
	char *log;
	char case_log[PAGE_SIZE];
	struct munit_module *module;
	struct munit_case *mcase;

	module = munit_get_case_module_struct(module_name);
	if (!module) {
		MLOG_ERR("munit:%s get this module struct fail\n", module_name);
		return;
	}

	mcase = munit_get_module_case(module_name, module);
	if (!mcase) {
		MLOG_ERR("munit:%s get this case id fail\n", module_name);
		return;
	}

	log = mcase->log;
	len = PAGE_SIZE - strlen(log) - 1;
	if (len <= 0) {
		MLOG_ERR("munit:%s this log have no space\n", module_name);
		return;
	}

	va_start(args, fmt);
	vsnprintf(case_log, sizeof(case_log), fmt, args);
	va_end(args);

	strncat(log, case_log, len);

	return;
}
EXPORT_SYMBOL(munit_get_module_case_log);

static void munit_create_debugfs_file(char *name, struct dentry *root, void *data,
					struct file_operations *ops)
{
	debugfs_create_file(name, 0644, root, data, ops);
}

static void munit_module_case_create_fs(struct munit_module *module,
					struct munit_case *mcase)
{
	int i;
	struct munit_case *mtest;

	module->module_dir = debugfs_create_dir(module->name, munit_debugfs_rootdir);

	for (i = 0; i < module->num; i++) {
		mtest = &mcase[i];

		mtest->case_dir = debugfs_create_dir(mtest->name,
					module->module_dir);
		munit_create_debugfs_file("run", mtest->case_dir, mtest,
					&munit_fops);
		munit_create_debugfs_file("res", mtest->case_dir, mtest,
					&munit_res_fops);
		munit_create_debugfs_file("log", mtest->case_dir, mtest,
					&munit_fops);
		mtest->log = kmalloc(PAGE_SIZE, GFP_KERNEL);
		if (!mtest->log)
			return;
	}
}

static void munit_module_case_decreate_fs(struct munit_module *module)
{
	int i;
	struct munit_case *mtest;

	for (i = 0; i < module->num; i++) {
		mtest = &module->mcase[i];

		debugfs_remove_recursive(mtest->case_dir);
	}

	debugfs_remove_recursive(module->module_dir);
}

void register_cases(char *name, struct munit_case *mcase)
{
	int num = 0;
	struct munit_module *module;
	struct munit_case *tcase;

	module = kmalloc(sizeof(*module), GFP_KERNEL);
	module->name = kmalloc(strlen(name) + 1, GFP_KERNEL);
	memcpy(module->name, name, strlen(name));
	module->name[strlen(name)] = '\0';

	for (tcase = mcase; tcase->func_run; tcase++)
		num++;

	module->num = num;
	module->mcase = mcase;

	mutex_lock(&case_module_mutex);
	list_add_tail(&module->list, &case_module);
	mutex_unlock(&case_module_mutex);

	munit_module_case_create_fs(module, mcase);
}
EXPORT_SYMBOL(register_cases);

void unregister_cases(char *name)
{
	struct list_head *itr;
	struct munit_module *module;

	mutex_lock(&case_module_mutex);
	if (!list_empty(&case_module)) {
		list_for_each(itr, &case_module) {
			module = list_entry(itr, struct munit_module, list);
			if (!strcmp(module->name, name)) {
				munit_module_case_decreate_fs(module);
				list_del(&module->list);
				kfree(module->name);
				kfree(module);
				break;
			}
		}
	}
	mutex_unlock(&case_module_mutex);
}
EXPORT_SYMBOL(unregister_cases);

static int munit_open(struct inode *inode, struct file *filp)
{
	return 0;
}

static ssize_t munit_write(struct file *filp, const char __user *buffer,
			   size_t count, loff_t *pos)
{
	int ret;
	unsigned long long enable;
	struct munit_case *mcase;

	ret = kstrtoull_from_user(buffer, count, 10, &enable);
	if (ret)
		return -EINVAL;

	if (enable != 1)
		return -EINVAL;

	mcase = (struct munit_case *)filp->f_inode->i_private;
	if (mcase->func_run){
		memset(mcase->log, 0, PAGE_SIZE);
		ret = mcase->func_run();
		if (ret == 0)
			mcase->succ = true;
		else
			mcase->succ = false;

		if (mcase->succ) {
			strcpy(mcase->res, "PASS");
			MLOG_INFO("case:%s PASS", mcase->name);
		} else {
			strcpy(mcase->res, "NOT PASS");
			MLOG_ERR("case:%s NOT PASS", mcase->name);
		}
	}

		return count;
}

static ssize_t munit_read(struct file *filp, char __user *buffer,
			  size_t count, loff_t *pos)
{
	int ret;
	char *temp;
	struct munit_case *mcase;

	mcase = (struct munit_case *)filp->f_inode->i_private;

	if (*pos >= strlen(mcase->log))
		return 0;

	count = min_t(size_t, count, strlen(mcase->log) - *pos);
	count = min_t(size_t, count, PAGE_SIZE);

	temp = kmalloc(count, GFP_KERNEL);

	memcpy(temp, mcase->log + *pos, count);

	ret = copy_to_user(buffer, temp, count);
	if (ret)
		return -EFAULT;

	*pos += count;

	kfree(temp);

	return count;
}

static ssize_t munit_res_read(struct file *filp, char __user *buffer,
			      size_t count, loff_t *pos)
{
	int ret;
	char res;
	char *temp;
	struct munit_case *mcase;

	mcase = (struct munit_case *)filp->f_inode->i_private;

	if (!count)
		return 0;

	if (*pos >= strlen(mcase->res))
		return 0;

	count = min_t(size_t, count, strlen(mcase->res) - *pos);

	temp = kmalloc(count, GFP_KERNEL);

	memcpy(temp, mcase->res + *pos, count);

	ret = copy_to_user(buffer, temp, count);
	if (ret)
		return -EFAULT;

	*pos += count;

	kfree(temp);

	return count;
}

struct file_operations munit_fops = {
	.open = munit_open,
	.write = munit_write,
	.read = munit_read,
};

struct file_operations munit_res_fops = {
	.open = munit_open,
	.read = munit_res_read,
};

static void munit_init_debugfs(void)
{
	if (!munit_debugfs_rootdir)
		munit_debugfs_rootdir = debugfs_create_dir(DEVICE_NAME, NULL);
}

static void munit_deinit_debugfs(void)
{
	debugfs_remove_recursive(munit_debugfs_rootdir);
}

static int munit_cdev_register(struct munit_dev *mdev)
{
	struct device *device;

	if (alloc_chrdev_region(&mdev->devno, 0, 1, DEVICE_NAME)) {
		MLOG_ERR("failed to create the munit cdev region\n");
		return -EFAULT;
	}

	mdev->class = class_create(THIS_MODULE, DEVICE_NAME);
	if (IS_ERR(mdev->class)) {
		MLOG_ERR("failed to create the munit class\n");
		goto munit_unregister_region;
	}

	cdev_init(&mdev->cdev, &munit_fops);
	if (cdev_add(&mdev->cdev, mdev->devno, 1)) {
		MLOG_ERR("failed to add the munit cdev\n");
		goto munit_free_class;
	}

	device = device_create(mdev->class, NULL, mdev->devno,
			       NULL, DEVICE_NAME);
	if (IS_ERR(device)) {
		MLOG_ERR("failed to create device for munit\n");
		goto munit_del_cdev;
	}

	return 0;

munit_del_cdev:
	cdev_del(&mdev->cdev);
munit_free_class:
	class_destroy(mdev->class);
munit_unregister_region:
	unregister_chrdev_region(mdev->devno, 1);;

	return -EFAULT;
}

static int __init munit_dev_init(void)
{
	int ret = 0;
	struct munit_dev *munit = &mdev;

	ret = munit_cdev_register(munit);
	if (ret) {
		MLOG_ERR("munit failed to register a char device\n");
		return ret;
	}

	munit_init_debugfs();

	return ret;
}

static void __exit munit_dev_exit(void)
{
	struct munit_dev *munit = &mdev;

	device_destroy(munit->class, munit->devno);
	cdev_del(&munit->cdev);
	class_destroy(munit->class);
	unregister_chrdev_region(munit->devno, 1);
	munit_deinit_debugfs();
}

module_init(munit_dev_init);
module_exit(munit_dev_exit);

MODULE_DESCRIPTION("Alibaba munit driver");
MODULE_AUTHOR("Jiankang Chen <jkchen@linux.alibaba.com>");
MODULE_LICENSE("GPL");
