# File trees

        ├── beecases
        │   ├── beekernels
        │   │   ├── example
        │   │   │   ├── example.c
        │   │   │   └── Makefile
        │   │   └── Makefile
        │   ├── beeusers
        │   │   ├── example
        │   │   │   ├── bee_example.sh
        │   │   │   └── Makefile
        │   │   └── Makefile
        │   └── Makefile
        ├── beeconfigs
        │   ├── bee_env.mk
        │   ├── bee_module.mk
        │   ├── bee_target.mk
        │   └── bee_user.mk
        ├── beeinc
        │   └── munit.h
        ├── beetools
        │   ├── Makefile
        │   └── munit
        │       ├── example
        │       │   ├── example_drv.c
        │       │   └── Makefile
        │       ├── Makefile
        │       └── munit_drv.c
        ├── bt
        │   ├── bee_run_kernels
        │   ├── bee_run_users
        │   └── bt
        ├── build.sh
        └── Makefile


# How to add and run a new kernel case in Bee

We add a example case, you can see [kernel case](../package/bee/src/beecases/beekernels/example/) and [user case](../package/bee/src/beecases/beeusers/example/)

## Makefile for new kernel case

All kernel cases in Bee are running based on Munit. Bee has a simple build system that requires just two steps to add a new use case:

1. Add your source file to beecases/beekernels
2. Add Makefile as follow:

        bee_srcdir ?= ../../..
        include $(bee_srcdir)/beeconfigs/bee_env.mk
        TARGET := example
        INSTALL_SH := *.sh
        INSTALL_KO := *.ko
        obj-m := $(TARGET).o
        include $(bee_srcdir)/beeconfigs/bee_module.mk

## Run your case

You can then build and run your cases using Bee framework like [how to build and run Bee](bee.md)

# How to add and run a new user case in Bee

## Makefile for a new use case

1. Add your source file to beecases/beeusers
2. Add Makefile as follow:

        bee_srcdir ?= ../../..
        include $(bee_srcdir)/beeconfigs/bee_env.mk
        INSTALL_TARGETS := *.sh
        MAKE_TARGETS :=
        include $(bee_srcdir)/beeconfigs/bee_user.mk

## Run your case

There are some differences between user and kernel cases. A shell script is needed for each user case.

1. Script naming rules: bee_casesname.sh, where 'casename' is the name of the case you want to add e.g. bee_mpam.sh etc.
2. Script template: you need add some functions to run single or all cases and declare your case list in script, like [this](../package/bee/src/beecases/beeusers/example/bee_example.sh).
