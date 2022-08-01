# Background

Currently, kernel use kunit as the unit testing framework, but kunit has some flaws and does not suitable for silicon validation, the main problems are.；

- Kunit runs all the use cases synchronously when insmod, and does not allow for individual execution of use cases, and uniform execution of use cases asynchronously；
- Not all kernels have kunit；
- Kunit does not have a user-friendly log system, and it is not possible to easily get the failed use cases and view the logs of the failed use cases.；

Therefore, we introduced Munit, a unit testing framework that can run single and multiple cases **asynchronously**, **easy to run on all kernel versions**, and have **a well-developed logging system**.

The comparisons of Munit and Kunit are as follows:

|   | KUNIT | MUNIT |
|  ----     | ----     |  ----    |
| Portability |  NO | YES |
| Asynchronously |  NO | YES |
| Case log record |  NO | YES |
| Case result record |  NO | YES |

# Munit basic framework

<img src=munit.png width="50%">

The Munit module registers each unit test set as a separate kernel module and manages these modules uniformly. In addition to the kernel modules, some user space softwares handle munit logs and control the execution of cases, these software is used to provide upper layer services in general.

&emsp;&emsp;&emsp;\+ ----- +\
&emsp;&emsp;&emsp;&nbsp;|&ensp;shell&ensp;|\
&emsp;&emsp;&emsp;\+ ----- +\
&emsp;&emsp;&emsp;&emsp;&ensp;∧\
&emsp;&emsp;&emsp;&emsp;&ensp; |\
&emsp;&emsp;&emsp;&emsp;&ensp; |\
&emsp;&emsp;&ensp;\+ --------- +\
&emsp;&emsp;&ensp; |&nbsp;&nbsp;debug fs&nbsp;&nbsp;|\
&emsp;&emsp;&ensp;\+ --------- +\
&emsp;&emsp;&emsp;&emsp;&ensp;∧\
&emsp;&emsp;&emsp;&emsp;&ensp; |\
&emsp;&emsp;&emsp;&emsp;&ensp; |&emsp;&emsp;&emsp;&emsp;&emsp;\
&emsp;&emsp;&ensp;\+ --------- +&emsp;<--------&emsp;\+ ---------------- +\
&emsp;&emsp;&ensp; |&emsp;&emsp;&emsp;&emsp;&emsp;|&ensp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;|&nbsp;&nbsp;smmu cases ko&nbsp;&nbsp;|\
&emsp;&emsp;&ensp; |&emsp;&emsp;&emsp;&emsp;&emsp;|&emsp;-------->&emsp;\+ ---------------- +\
&emsp;&emsp;&ensp; |&nbsp;&nbsp; munit ko&nbsp;&nbsp;|\
&emsp;&emsp;&ensp; |&emsp;&emsp;&emsp;&emsp;&emsp;|&emsp;-------->&emsp;\+ -------------- +\
&emsp;&emsp;&ensp; |&emsp;&emsp;&emsp;&emsp;&emsp;|&ensp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;|&nbsp;&nbsp;pcie cases ko&nbsp;&nbsp;|\
&emsp;&emsp;&ensp;\+ --------- +&emsp;<--------&emsp;\+ -------------- +

As shown in the figure above, the munit framework module is first registered to the kernel and then serves the user-developed kernel use case unit test module, which finally registers all use cases in the munit framework module.
munit provides debugfs interface services to the user space and uses the shell scripting framework to execute cases and collect results.

# Munit debugfs and case execution

Munit's debugfs are divided into 3 levels, the first level is the munit directory, the second level is the module name, the third level is the use case directory, and each use case directory is divided into 3 files:

    /sys/kernel/debug/munit/
    └── example
        └── example_0
            ├── log
            ├── res
            └── run

- res: case execution results;
- log: case log;
- run: run case;

You can run case with follow command:

    echo 1 > /sys/kernel/debug/munit/example/example_0/run

# Munit log collection and result statistics

You can see case log with follow command:

    cat /sys/kernel/debug/munit/example/example_0/log

You can see case result with follow command:

    cat /sys/kernel/debug/munit/example/example_0/res

# Munit macro

Munit has the same macros as kunit：

    MUNIT_EXPECT_EQ
    MUNIT_EXPECT_TRUE
    MUNIT_EXPECT_FALSE
    MUNIT_EXPECT_NE

# Munit case example

See [example](../package/bee/src/beecases/beekernels/example/example.c)

build, install, and insmod munit and example kernel module, then you can see：

    /sys/kernel/debug/munit/
    └── example
        └── example_0
            ├── log
            ├── res
            └── run

run example_0：

    #echo 1 > /sys/kernel/debug/munit/example/example_0/run

the result is：

    #cat /sys/kernel/debug/munit/example/example_0/res
    FAIL

the log is；

    #cat /sys/kernel/debug/munit/example/example_0/log
    run example_0 test!
